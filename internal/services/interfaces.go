package services

import (
	"github.com/varun-singhh/auth-service/internal/models"
	"gofr.dev/pkg/gofr"
)

type Authorizer interface {
	Login(ctx *gofr.Context, requestUser *models.User) (*models.AuthorizedResponse, error)
	Signup(ctx *gofr.Context, requestUser *models.User) (*models.AuthorizedResponse, error)
	RefreshToken(ctx *gofr.Context, signedToken string) (*models.AuthorizedResponse, error)
	ValidateToken(ctx *gofr.Context, signedToken string) (bool, error)
	ForgotPassword(ctx *gofr.Context, email, permission string) (interface{}, error)
	ResetPassword(ctx *gofr.Context, permission, newPassword, token string) (string, error)
	VerifyAccount(ctx *gofr.Context, code, permission, token string) (interface{}, error)
	GenerateSignupVerificationCode(ctx *gofr.Context, email string) (interface{}, error)
}
