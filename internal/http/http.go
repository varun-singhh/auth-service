package http

import (
	"github.com/varun-singhh/auth-service/internal/models"
	"github.com/varun-singhh/auth-service/internal/services"
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
)

type handler struct {
	authService services.Authorizer
}

func New(authService services.Authorizer) *handler {
	return &handler{authService: authService}
}

func (h *handler) Login(ctx *gofr.Context) (interface{}, error) {

	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"json request body"}}
	}

	resp, err := h.authService.Login(ctx, &user)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *handler) Signup(ctx *gofr.Context) (interface{}, error) {

	var user models.User
	err := ctx.Bind(&user)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"json request body"}}
	}

	resp, err := h.authService.Signup(ctx, &user)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *handler) ForgotPassword(ctx *gofr.Context) (interface{}, error) {
	var user models.User

	err := ctx.Bind(&user)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"json request body"}}
	}

	resp, err := h.authService.ForgotPassword(ctx, user.Email, user.Permission)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *handler) ResetPassword(ctx *gofr.Context) (interface{}, error) {
	params := ctx.Request().URL.Query()

	var user models.User

	err := ctx.Bind(&user)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"json request body"}}
	}

	resp, err := h.authService.ResetPassword(ctx, user.Permission, user.Password, params.Get("token"))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *handler) RefreshToken(ctx *gofr.Context) (interface{}, error) {
	token := ctx.PathParam("token")

	resp, err := h.authService.RefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *handler) ValidateToken(ctx *gofr.Context) (interface{}, error) {
	token := ctx.PathParam("token")

	_, err := h.authService.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *handler) VerifyAccount(ctx *gofr.Context) (interface{}, error) {
	var input models.VerifyAccount

	err := ctx.Bind(&input)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"json request body"}}
	}

	token := ctx.PathParam("token")

	res, err := h.authService.VerifyAccount(ctx, input.Code, input.Permission, token)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *handler) GenerateSignupVerificationCode(ctx *gofr.Context) (interface{}, error) {
	q := ctx.Request().URL.Query()
	email := q.Get("email")

	resp, err := h.authService.GenerateSignupVerificationCode(ctx, email)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
