package stores

import (
	"database/sql"
	"fmt"
	"github.com/varun-singhh/auth-service/internal/models"
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
)

type auth struct {
}

func New() *auth {
	return &auth{}
}

func (a *auth) Get(ctx *gofr.Context, user *models.User) (*models.User, error) {

	var (
		respUser models.User
		query    string
		params   []interface{}
	)

	if user.Phone != "" && user.Email != "" {
		query = `SELECT id, email, phone, password, permissions, status, created_at FROM users WHERE (id=$1 OR (email=$2 AND phone=$3)) AND permissions=$4`
		params = append(params, user.ID, user.Email, user.Phone, user.Permission)

	} else if user.Phone == "" {
		query = `SELECT id, email, phone, password, permissions, status, created_at FROM users WHERE (id=$1 OR email=$2) AND permissions=$3`
		params = append(params, user.ID, user.Email, user.Permission)
	} else if user.Email == "" {
		query = `SELECT id, email, phone, password, permissions, status, created_at FROM users WHERE (id=$1 OR phone=$2) AND permissions=$3`
		params = append(params, user.ID, user.Phone, user.Permission)
	}

	err := ctx.DB().QueryRow(query, params...).Scan(&respUser.ID, &respUser.Email, &respUser.Phone, &respUser.Password, &respUser.Permission, &respUser.Status, &respUser.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, &errors.EntityNotFound{
			Entity: user.Permission,
			ID:     user.Email,
		}
	}

	if err != nil {
		ctx.Logger.Errorf("error while fetching user detail from db, Err: %v", err)
		return nil, &errors.DB{Err: err}
	}

	return &respUser, nil
}

func (a *auth) Create(ctx *gofr.Context, user *models.User) (*models.User, error) {
	query := `INSERT INTO users (email,phone,password,permissions) VALUES($1,$2,$3,$4)`

	_, err := ctx.DB().Exec(query, user.Email, user.Phone, user.Password, user.Permission)
	if err != nil {
		ctx.Logger.Errorf(fmt.Sprintf("error while creating user, Error is %v", err.Error()))
		return nil, &errors.DB{Err: err}
	}

	return a.Get(ctx, user)
}

func (a *auth) ResetPassword(ctx *gofr.Context, user *models.User, newPassword string) error {
	query := `UPDATE users SET password = $1 WHERE (email = $2 OR phone =$3) AND permissions = $4`
	_, err := ctx.DB().Exec(query, newPassword, user.Email, user.Phone, user.Permission)
	if err != nil {
		ctx.Logger.Errorf("error while resetting user password")
		return err
	}

	return nil
}

func (a *auth) VerifyAccount(ctx *gofr.Context, user *models.User, status string) error {

	query := `UPDATE  users  SET status = $1 WHERE id = $2`

	_, err := ctx.DB().Exec(query, status, user.ID)
	if err != nil {
		ctx.Logger.Errorf("failed to update account verification status, err:%v", err)
		return &errors.DB{Err: err}
	}

	return nil
}
