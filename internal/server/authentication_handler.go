package server

import (
	"context"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Authenticator interface {
	Register(context.Context, auth.RegisterRequest) (auth.CredentialsResponse, error)
	Login(context.Context, auth.LoginRequest) (auth.CredentialsResponse, error)
	ResetPassword(context.Context, auth.PasswordResetRequest) error
}

type AuthenticationHandler struct {
	authenticator Authenticator
}

func NewAuthenticatorHandler(authenticator Authenticator) *AuthenticationHandler {
	return &AuthenticationHandler{
		authenticator: authenticator,
	}
}

func (h *AuthenticationHandler) Routes() []Route {
	return []Route{
		{Method: http.MethodPost, Path: "/register", Handler: h.register, Public: true},
		{Method: http.MethodPost, Path: "/login", Handler: h.login, Public: true},
		{Method: http.MethodPatch, Path: "/password-reset", Handler: h.resetPassword, Public: true},
	}
}

// register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce  json
// @Param payload body auth.RegisterRequest true "Register Payload"
// @Success 201 {object} auth.CredentialsResponse
// @Failure 400 {object} api.Error
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /register [post]
func (h *AuthenticationHandler) register(c *gin.Context) {
	var request auth.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "RegisterRequest", Err: err})
		return
	}

	response, err := h.authenticator.Register(c.Request.Context(), request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// login godoc
// @Summary Login using email and password
// @Tags Auth
// @Accept json
// @Produce  json
// @Param payload body auth.LoginRequest true "Login Payload"
// @Success 201 {object} auth.CredentialsResponse
// @Failure 400 {object} api.Error
// @Failure 401 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /login [post]
func (h *AuthenticationHandler) login(c *gin.Context) {
	var request auth.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "LoginRequest", Err: err})
		return
	}

	response, err := h.authenticator.Login(c.Request.Context(), request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// resetPassword godoc
// @Summary Reset the password for the given email
// @Tags Auth
// @Accept json
// @Produce  json
// @Param payload body auth.PasswordResetRequest true "Register Payload"
// @Success 204 "success"
// @Failure 400 {object} api.Error
// @Failure 500 {object} api.Error
// @Router /password-reset [post]
func (h *AuthenticationHandler) resetPassword(c *gin.Context) {
	var request auth.PasswordResetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(&common.ErrNotValid{Input: "PasswordResetRequest", Err: err})
		return
	}

	if err := h.authenticator.ResetPassword(c.Request.Context(), request); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
