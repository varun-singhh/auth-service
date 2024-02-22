package services

import (
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/varun-singhh/auth-service/internal/services/email"
	"gofr.dev/pkg/gofr"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gofr.dev/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/varun-singhh/auth-service/internal/models"
)

type JWTClaim struct {
	Email    string
	Password string
	UserID   string
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

var jwtResetKey = []byte(os.Getenv("JWT_RESET_SECRET"))

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPassword(userPassword, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func generateJWT(email, password, id string) (tokenString string, err error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &JWTClaim{
		Email:    email,
		Password: password,
		UserID:   id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func validateToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusUnauthorized,
			Code:       http.StatusText(http.StatusUnauthorized),
			Reason:     "invalid token",
		}
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, &errors.Response{
			StatusCode: http.StatusUnauthorized,
			Code:       http.StatusText(http.StatusUnauthorized),
			Reason:     "couldn't parse claims",
		}
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, &errors.Response{
			StatusCode: http.StatusUnauthorized,
			Code:       http.StatusText(http.StatusUnauthorized),
			Reason:     "token expired",
		}
	}
	return claims, err
}

func generateResetPasswordJWT(email string, password string) (tokenString string, err error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &JWTClaim{
		Email:    email,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtResetKey)
	return
}

func validateResetPasswordToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtResetKey), nil
		},
	)
	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusBadRequest,
			Code:       http.StatusText(http.StatusBadRequest),
			Reason:     "token invalid",
		}
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, &errors.Response{
			StatusCode: http.StatusBadRequest,
			Code:       http.StatusText(http.StatusBadRequest),
			Reason:     "couldn't parse claims",
		}
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, &errors.Response{
			StatusCode: http.StatusBadRequest,
			Code:       http.StatusText(http.StatusBadRequest),
			Reason:     "token expired",
		}
	}
	return claims, err
}

func validateUser(user *models.User) error {
	if user.Password == "" {
		return errors.MissingParam{Param: []string{"password"}}
	}

	if user.Email == "" && user.Phone == "" {
		return errors.MissingParam{Param: []string{"email/phone"}}
	}

	if user.Permission == "" {
		return errors.MissingParam{Param: []string{"permission"}}
	}

	if !(user.Permission == "PATIENT" || user.Permission == "DOCTOR" || user.Permission == "ADMIN" || user.Permission == "MANAGER") {
		return &errors.Response{
			StatusCode: http.StatusBadRequest,
			Code:       http.StatusText(http.StatusBadRequest),
			Reason:     "permission type " + user.Permission + " not allowed",
		}
	}

	return nil
}

func verifyCode(ctx *gofr.Context, code, userID string) (bool, error) {
	strCmd := ctx.Redis.Get(ctx, userID)
	if strCmd.Err() != nil {
		return false, &errors.Response{
			StatusCode: 400,
			Code:       "400",
			Reason:     "verification code expired",
		}
	}

	if strCmd.Val() != code {
		return false, &errors.Response{
			StatusCode: 400,
			Code:       "400",
			Reason:     "invalid code, try again",
		}
	}

	return true, nil
}

func generateCode(ctx *gofr.Context, key string) (string, error) {
	code, err := gonanoid.Generate("0123456789", 8)
	if err != nil {
		ctx.Logger.Errorf(fmt.Sprintf("error generating verification code, Error: %v", err))
		return "", err
	}

	ctx.Redis.Set(ctx, key, code, 2*time.Minute)
	return code, nil
}

func sendEmail(ctx *gofr.Context, to, cc, bcc []string, msg email.EmailData) error {
	sender, err := email.NewSender(&models.EmailConfig{
		Host:      ctx.Config.Get("DEFAULT_SMTP_HOST"),
		Port:      ctx.Config.Get("DEFAULT_SMTP_PORT"),
		Username:  ctx.Config.Get("DEFAULT_SMTP_USERNAME"),
		Password:  ctx.Config.Get("DEFAULT_SMTP_PASSWORD"),
		FromEmail: ctx.Config.Get("DEFAULT_SMTP_FROM_EMAIL"),
		FromName:  ctx.Config.Get("DEFAULT_SMTP_FROM_NAME"),
		SmtpAuth:  ctx.Config.Get("DEFAULT_SMTP_AUTH_TYPE"),
		IsSmtp:    ctx.Config.Get("DEFAULT_IS_SMTPS"),
	})
	if err != nil {
		return err
	}

	err = sender.Send(ctx, to, cc, bcc, msg)
	if err != nil {
		return err
	}

	return nil
}
