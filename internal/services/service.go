package services

import (
	errors2 "errors"
	"fmt"
	"github.com/varun-singhh/auth-service/internal/services/email"
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/types"
	"strconv"

	"github.com/varun-singhh/auth-service/internal/models"
	"github.com/varun-singhh/auth-service/internal/stores"
)

type auth struct {
	authStore   stores.Authorizer
	emailConfig *models.EmailConfig
}

func New(authStore stores.Authorizer) *auth {
	return &auth{authStore: authStore}
}

func (a *auth) Login(ctx *gofr.Context, requestUser *models.User) (*models.AuthorizedResponse, error) {
	err := validateUser(requestUser)
	if err != nil {
		return nil, err
	}

	user, err := a.authStore.Get(ctx, requestUser)
	if err != nil {
		return nil, err
	}

	if err = checkPassword(user.Password, requestUser.Password); err != nil {
		return nil, &errors.Response{
			StatusCode: 400,
			Code:       "400",
			Reason:     "invalid password, please try again",
		}
	}

	tokenString, err := generateJWT(user.Email, user.Password, strconv.Itoa(user.ID))

	return &models.AuthorizedResponse{User: user, Token: tokenString, Message: "user logged in successfully"}, nil
}

func (a *auth) Signup(ctx *gofr.Context, requestUser *models.User) (*models.AuthorizedResponse, error) {
	var entityNotFound *errors.EntityNotFound

	err := validateUser(requestUser)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hashPassword(requestUser.Password)
	if err != nil {
		return nil, err
	}

	requestUser.Password = hashedPassword

	isExistingUser, err := a.authStore.Get(ctx, requestUser)
	if !errors2.As(err, &entityNotFound) && err != nil {
		return nil, err
	}

	if isExistingUser != nil && isExistingUser.ID > 0 {
		return nil, &errors.EntityAlreadyExists{}
	}

	user, err := a.authStore.Create(ctx, requestUser)
	if err != nil {
		return nil, err
	}

	return &models.AuthorizedResponse{User: user, Token: "", Message: "user created successfully, please verify your registered email"}, nil
}

func (a *auth) RefreshToken(ctx *gofr.Context, signedToken string) (*models.AuthorizedResponse, error) {

	claims, err := validateToken(signedToken)
	if err != nil {
		return nil, err
	}

	tokenString, err := generateJWT(claims.Email, claims.Password, claims.UserID)

	err = sendEmail(ctx, []string{"varunsingh00712@gmail.com"}, nil, nil, email.EmailData{EmailType: email.PasswordReset, DataMap: map[string]string{"code": tokenString, "user": "varun singh"}})
	if err != nil {
		ctx.Logger.Errorf("failed to send mail, err: %v", err)
	}

	return &models.AuthorizedResponse{Token: tokenString, Message: "token refreshed successfully"}, nil
}

func (a *auth) ValidateToken(ctx *gofr.Context, signedToken string) (bool, error) {

	claims, err := validateToken(signedToken)
	if err != nil {
		return false, err
	}

	if claims.UserID != "" {
		return true, nil
	}

	return false, nil
}

func (a *auth) ForgotPassword(ctx *gofr.Context, mail, permission string) (interface{}, error) {
	user, err := a.authStore.Get(ctx, &models.User{Email: mail, Permission: permission})
	if err != nil {
		return nil, errors.Error(err.Error())
	}

	if user.ID < 1 {
		return nil, errors.EntityNotFound{Entity: mail}
	}

	token, err := generateResetPasswordJWT(mail, user.Password)
	if err != nil {
		return nil, err
	}

	if user.Email != "" {
		err = sendEmail(ctx, []string{"varunsingh00712@gmail.com"}, nil, nil, email.EmailData{EmailType: email.PasswordReset, DataMap: map[string]string{"resetLink": "localhost:2222/reset-password?token=" + token, "user": "varun singh"}})
		if err != nil {
			ctx.Logger.Errorf("failed to send mail for verification code, err: %v", err)
		}
	}

	return types.Response{Data: "reset link sent on registered mail"}, nil
}

func (a *auth) ResetPassword(ctx *gofr.Context, permission, newPassword, token string) (string, error) {
	claims, err := validateToken(token)
	if err != nil {
		return "", err
	}

	user, err := a.authStore.Get(ctx, &models.User{Email: claims.Email, Permission: permission})
	if err != nil {
		return "", err
	}

	if user.Password != claims.Password {
		return "", &errors.Response{
			StatusCode: 400,
			Code:       "400",
			Reason:     "user not authorized to reset",
		}
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return "", err
	}

	if err := a.authStore.ResetPassword(ctx, user, hashedPassword); err != nil {
		return "", err
	}

	return "password reset successfully", nil
}

func (a *auth) VerifyAccount(ctx *gofr.Context, code, permission, token string) (interface{}, error) {
	claims, err := validateToken(token)
	if err != nil {
		return nil, err
	}

	id, _ := strconv.Atoi(claims.Id)

	user, err := a.authStore.Get(ctx, &models.User{Email: claims.Email, ID: id, Permission: permission})
	if err != nil {
		return nil, err
	}

	if user.Status == "VERIFIED" {
		return types.Response{
			Data: "account already verified",
		}, nil
	}

	if user.Password != claims.Password {
		return nil, &errors.Response{
			StatusCode: 400,
			Code:       "400",
			Reason:     "invalid token, user not authorized to perform action",
		}
	}

	status := "PENDING"

	isVerified, err := verifyCode(ctx, code, user.Email)
	if err != nil {
		return nil, err
	}

	if isVerified {
		status = "VERIFIED"
	}

	if err = a.authStore.VerifyAccount(ctx, user, status); err != nil {
		return nil, err
	}

	return types.Response{Data: "account verified successfully"}, nil
}

func (a *auth) GenerateSignupVerificationCode(ctx *gofr.Context, mail string) (interface{}, error) {
	val := ctx.Redis.Get(ctx, mail)
	if val.Val() != "" {
		return types.Response{Data: fmt.Sprintf("code already sent, retry after 2 minutes")}, nil
	}

	code, err := generateCode(ctx, mail)

	err = sendEmail(ctx, []string{"varunsingh00712@gmail.com"}, nil, nil, email.EmailData{EmailType: email.VerificationCode, DataMap: map[string]string{"code": code, "user": "varun singh"}})
	if err != nil {
		ctx.Logger.Errorf("failed to send mail for verification code, err: %v", err)
	}

	return types.Response{
		Data: "new code sent to registered email",
	}, nil
}
