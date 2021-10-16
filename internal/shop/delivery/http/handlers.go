package http

import (
	"bytes"
	"fmt"
	auth "github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/internal/shop/delivery/dto"
	"github.com/c0llinn/ebook-store/internal/shop/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type ShopService interface {
	FindOrders(query model.OrderQuery) (model.PaginatedOrders, error)
	FindOrderByID(id string) (model.Order, error)
	CreateOrder(order *model.Order) error
	UpdateOrder(order *model.Order) error
	CompleteOrder(orderID string) error
	DownloadOrder(orderID string) (io.ReadCloser, error)
}

type IDGenerator interface {
	NewID() string
}

type ShopHandler struct {
	service     ShopService
	idGenerator IDGenerator
}

func NewShopHandler(service ShopService, idGenerator IDGenerator) ShopHandler {
	return ShopHandler{service: service, idGenerator: idGenerator}
}

// getOrders godoc
// @Summary Fetch Orders
// @Tags Shop
// @Produce  json
// @Param payload body dto.SearchOrders true "Filters"
// @Success 200 {object} dto.PaginatedOrdersResponse
// @Failure 500 {object} api.Error
// @Router /orders [get]
func (h ShopHandler) getOrders(context *gin.Context) {
	var s dto.SearchOrders
	if err := context.ShouldBindQuery(&s); err != nil {
		context.Error(&common.ErrNotValid{Input: "SearchOrders", Err: err})
		return
	}

	user, err := h.getUserFromContext(context)
	if err != nil {
		context.Error(err)
		return
	}

	query := s.ToDomain()
	if !user.IsAdmin() {
		query.UserID = user.ID
	}

	paginatedOrders, err := h.service.FindOrders(query)
	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, dto.FromPaginatedOrders(paginatedOrders))
}

// getOrder godoc
// @Summary Fetch Order by ID
// @Tags Shop
// @Produce  json
// @Param id path string true "orderId ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /orders/{id} [get]
func (h ShopHandler) getOrder(context *gin.Context) {
	order, err := h.service.FindOrderByID(context.Param("id"))
	if err != nil {
		context.Error(err)
		return
	}

	user, err := h.getUserFromContext(context)
	if err != nil {
		context.Error(err)
		return
	}

	if !user.IsAdmin() && order.UserID != user.ID {
		context.Error(fmt.Errorf("you don't have permission to see this order"))
		return
	}

	context.JSON(http.StatusOK, dto.FromOrder(order))
}

// createOrder godoc
// @Summary Create a new Order
// @Tags Shop
// @Accept json
// @Produce  json
// @Param payload body dto.CreateOrder true "Order Payload"
// @Success 201 {object} dto.OrderResponse
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /orders [post]
func (h ShopHandler) createOrder(context *gin.Context) {
	var c dto.CreateOrder
	if err := context.ShouldBindJSON(&c); err != nil {
		context.Error(&common.ErrNotValid{Input: "CreateOrder", Err: err})
		return
	}

	user, err := h.getUserFromContext(context)
	if err != nil {
		context.Error(err)
		return
	}

	order := c.ToDomain(h.idGenerator.NewID(), user.ID)

	if err = h.service.CreateOrder(&order); err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusCreated, dto.FromOrder(order))
}

func (h ShopHandler) getUserFromContext(context *gin.Context) (user auth.User, err error) {
	value, exists := context.Get("user")
	if !exists {
		err = fmt.Errorf("not authenticated user")
		return
	}

	user = value.(auth.User)
	return
}

// downloadOrder godoc
// @Summary Download the book for the given Order
// @Tags Shop
// @Produce  application/pdf
// @Param payload body dto.CreateOrder true "Order Payload"
// @Success 200 "Success"
// @Failure 402 {object} api.Error
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /orders/{id}/download [get]
func (h ShopHandler) downloadOrder(context *gin.Context) {
	content, err := h.service.DownloadOrder(context.Param("id"))
	if err != nil {
		context.Error(err)
		return
	}

	buffer := new(bytes.Buffer)
	if _, err = buffer.ReadFrom(content); err != nil {
		context.Error(err)
		return
	}

	context.DataFromReader(http.StatusOK, int64(buffer.Len()), "application/pdf", buffer, nil)
}

// handleStripeWebhook godoc
// @Summary Handle stripe webhooks
// @Tags Shop
// @Accept json
// @Produce  json
// @Success 200 "Success"
// @Success 500 {object} api.Error
// @Router /stripe/webhook [post]
func (h ShopHandler) handleStripeWebhook(context *gin.Context) {
	var w dto.HandleStripeWebhook
	if err := context.ShouldBindJSON(&w); err != nil {
		context.Error(err)
		return
	}

	if w.Type == "payment_intent.succeeded" {
		orderID := w.Data["object"].(map[string]interface{})["metadata"].(map[string]interface{})["orderID"].(string)
		if err := h.service.CompleteOrder(orderID); err != nil {
			context.Error(err)
			return
		}
	}

	context.Status(http.StatusOK)
}
