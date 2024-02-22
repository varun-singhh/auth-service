package stores

import (
	"github.com/varun-singhh/auth-service/internal/models"
	"gofr.dev/pkg/gofr"
)

type Authorizer interface {
	Get(ctx *gofr.Context, user *models.User) (*models.User, error)
	Create(ctx *gofr.Context, user *models.User) (*models.User, error)
	ResetPassword(ctx *gofr.Context, user *models.User, newPassword string) error
	VerifyAccount(ctx *gofr.Context, user *models.User, status string) error
}
