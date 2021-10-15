//go:build unit
// +build unit

package usecase

import (
	"fmt"
	catalog "github.com/c0llinn/ebook-store/internal/catalog/model"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/internal/shop/mock"
	"github.com/c0llinn/ebook-store/internal/shop/model"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"strings"
	"testing"
)

const (
	findOrdersByQueryMethod   = "FindByQuery"
	findOrderByIDMethod       = "FindByID"
	createOrderMethod         = "Create"
	updateOrderMethod         = "Update"
	findBookByIDMethod        = "FindBookByID"
	createPaymentIntentMethod = "CreatePaymentIntentForOrder"
	getBookContent            = "GetBookContent"
)

type ShopUseCaseTestSuite struct {
	suite.Suite
	repo           *mock.OrderRepository
	paymentClient  *mock.PaymentClient
	catalogService *mock.CatalogService
	useCase        ShopUseCase
}

func (s *ShopUseCaseTestSuite) SetupTest() {
	s.repo = new(mock.OrderRepository)
	s.paymentClient = new(mock.PaymentClient)
	s.catalogService = new(mock.CatalogService)

	s.useCase = NewShopUseCase(s.repo, s.paymentClient, s.catalogService)
}

func TestShopUseCaseRun(t *testing.T) {
	suite.Run(t, new(ShopUseCaseTestSuite))
}

func (s *ShopUseCaseTestSuite) TestFindOrders_Successfully() {
	paginated := factory.NewPaginatedOrders()
	s.repo.On(findOrdersByQueryMethod, model.OrderQuery{}).Return(paginated, nil)

	actual, err := s.useCase.FindOrders(model.OrderQuery{})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), paginated, actual)
	s.repo.AssertCalled(s.T(), findOrdersByQueryMethod, model.OrderQuery{})
}

func (s *ShopUseCaseTestSuite) TestFindOrders_WithError() {
	s.repo.On(findOrdersByQueryMethod, model.OrderQuery{}).Return(model.PaginatedOrders{}, fmt.Errorf("some error"))

	_, err := s.useCase.FindOrders(model.OrderQuery{})

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrdersByQueryMethod, model.OrderQuery{})
}

func (s *ShopUseCaseTestSuite) TestFindOrderByID_Successfully() {
	order := factory.NewOrder()
	s.repo.On(findOrderByIDMethod, order.ID).Return(order, nil)

	actual, err := s.useCase.FindOrderByID(order.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), order, actual)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, order.ID)
}

func (s *ShopUseCaseTestSuite) TestFindOrderByID_WithError() {
	id := uuid.NewString()
	s.repo.On(findOrderByIDMethod, id).Return(model.Order{}, fmt.Errorf("some error"))

	_, err := s.useCase.FindOrderByID(id)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, id)
}

func (s *ShopUseCaseTestSuite) TestCreateOrder_WhenCatalogServiceFails() {
	order := factory.NewOrder()
	s.catalogService.On(findBookByIDMethod, order.BookID).Return(catalog.Book{}, fmt.Errorf("some error"))

	err := s.useCase.CreateOrder(&order)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, order.BookID)
	s.paymentClient.AssertNotCalled(s.T(), createPaymentIntentMethod, &order)
	s.repo.AssertNotCalled(s.T(), createOrderMethod, &order)
}

func (s *ShopUseCaseTestSuite) TestCreateOrder_WhenPaymentClientFails() {
	order := factory.NewOrder()
	book := factory.NewBook()
	s.catalogService.On(findBookByIDMethod, order.BookID).Return(book, nil)
	s.paymentClient.On(createPaymentIntentMethod, &order).Return(fmt.Errorf("some error"))

	err := s.useCase.CreateOrder(&order)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, &order)
	s.catalogService.AssertNotCalled(s.T(), findBookByIDMethod, order.ID)
	s.repo.AssertNotCalled(s.T(), createOrderMethod, &order)
}

func (s *ShopUseCaseTestSuite) TestCreateOrder_WhenRepositoryFails() {
	order := factory.NewOrder()
	book := factory.NewBook()
	s.paymentClient.On(createPaymentIntentMethod, &order).Return(nil)
	s.catalogService.On(findBookByIDMethod, order.BookID).Return(book, nil)

	updatedOrder := order
	updatedOrder.Total = int64(book.Price)
	s.repo.On(createOrderMethod, &updatedOrder).Return(fmt.Errorf("some error"))

	err := s.useCase.CreateOrder(&order)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, &order)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, order.BookID)
	s.repo.AssertCalled(s.T(), createOrderMethod, &updatedOrder)
}

func (s *ShopUseCaseTestSuite) TestCreateOrder_Successfully() {
	order := factory.NewOrder()
	book := factory.NewBook()
	s.paymentClient.On(createPaymentIntentMethod, &order).Return(nil)
	s.catalogService.On(findBookByIDMethod, order.BookID).Return(book, nil)

	updatedOrder := order
	updatedOrder.Total = int64(book.Price)
	s.repo.On(createOrderMethod, &updatedOrder).Return(nil)

	err := s.useCase.CreateOrder(&order)

	assert.Nil(s.T(), err)

	s.paymentClient.AssertCalled(s.T(), createPaymentIntentMethod, &order)
	s.catalogService.AssertCalled(s.T(), findBookByIDMethod, order.BookID)
	s.repo.AssertCalled(s.T(), createOrderMethod, &updatedOrder)
}

func (s *ShopUseCaseTestSuite) TestUpdateOrder_Successfully() {
	order := factory.NewOrder()
	s.repo.On(updateOrderMethod, &order).Return(nil)

	err := s.useCase.UpdateOrder(&order)

	assert.Nil(s.T(), err)
	s.repo.AssertCalled(s.T(), updateOrderMethod, &order)
}

func (s *ShopUseCaseTestSuite) TestUpdateOrder_WithError() {
	order := factory.NewOrder()
	s.repo.On(updateOrderMethod, &order).Return(fmt.Errorf("some error"))

	err := s.useCase.UpdateOrder(&order)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), updateOrderMethod, &order)
}

func (s *ShopUseCaseTestSuite) TestCompleteOrder_WhenOrderCouldNotBeFound() {
	id := uuid.NewString()
	s.repo.On(findOrderByIDMethod, id).Return(model.Order{}, fmt.Errorf("some error"))

	err := s.useCase.CompleteOrder(id)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, id)
	s.repo.AssertNumberOfCalls(s.T(), updateOrderMethod, 0)
}

func (s *ShopUseCaseTestSuite) TestCompleteOrder_WhenUpdateFails() {
	order := factory.NewOrder()
	s.repo.On(findOrderByIDMethod, order.ID).Return(order, nil)

	order.Complete()
	s.repo.On(updateOrderMethod, &order).Return(fmt.Errorf("some error"))

	err := s.useCase.CompleteOrder(order.ID)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, order.ID)
	s.repo.AssertCalled(s.T(), updateOrderMethod, &order)
}

func (s *ShopUseCaseTestSuite) TestCompleteOrder_Successfully() {
	order := factory.NewOrder()
	s.repo.On(findOrderByIDMethod, order.ID).Return(order, nil)

	order.Complete()
	s.repo.On(updateOrderMethod, &order).Return(nil)

	err := s.useCase.CompleteOrder(order.ID)

	assert.Nil(s.T(), err)
	s.repo.AssertCalled(s.T(), findOrderByIDMethod, order.ID)
	s.repo.AssertCalled(s.T(), updateOrderMethod, &order)
}

func (s *ShopUseCaseTestSuite) TestDownloadOrder_WhenOrderCouldNotBeFound() {
	order := factory.NewOrder()
	s.repo.On(findOrderByIDMethod, order.ID).Return(model.Order{}, fmt.Errorf("some error"))

	_, err := s.useCase.DownloadOrder(order.ID)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, order.ID)
	s.catalogService.AssertNotCalled(s.T(), getBookContent, order.BookID)
}

func (s *ShopUseCaseTestSuite) TestDownloadOrder_WhenOrderIsNotPaid() {
	order := factory.NewOrder()
	order.Status = model.Pending
	s.repo.On(findOrderByIDMethod, order.ID).Return(order, nil)

	_, err := s.useCase.DownloadOrder(order.ID)

	assert.Equal(s.T(), &common.ErrOrderNotPaid{Err: fmt.Errorf("only books from paid order can be downloaded")}, err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, order.ID)
	s.catalogService.AssertNotCalled(s.T(), getBookContent, order.BookID)
}

func (s *ShopUseCaseTestSuite) TestDownloadOrder_WithError() {
	order := factory.NewOrder()
	order.Status = model.Paid
	s.repo.On(findOrderByIDMethod, order.ID).Return(order, nil)
	s.catalogService.On(getBookContent, order.BookID).Return(nil, fmt.Errorf("some error"))

	_, err := s.useCase.DownloadOrder(order.ID)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, order.ID)
	s.catalogService.AssertCalled(s.T(), getBookContent, order.BookID)
}

func (s *ShopUseCaseTestSuite) TestDownloadOrder_Successfully() {
	order := factory.NewOrder()
	order.Status = model.Paid
	s.repo.On(findOrderByIDMethod, order.ID).Return(order, nil)
	s.catalogService.On(getBookContent, order.BookID).Return(io.NopCloser(strings.NewReader("test")), nil)

	actual, err := s.useCase.DownloadOrder(order.ID)

	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)

	s.repo.AssertCalled(s.T(), findOrderByIDMethod, order.ID)
	s.catalogService.AssertCalled(s.T(), getBookContent, order.BookID)
}