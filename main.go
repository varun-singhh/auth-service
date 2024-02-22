package main

import (
	"gofr.dev/pkg/gofr"

	httpLayer "github.com/varun-singhh/auth-service/internal/http"
	"github.com/varun-singhh/auth-service/internal/services"
	"github.com/varun-singhh/auth-service/internal/stores"
)

func main() {

	app := gofr.New()

	// store layer setup
	authStore := stores.New()

	// service layer setup
	authSvc := services.New(authStore)

	// handler setup
	httpHandler := httpLayer.New(authSvc)

	app.POST("/login", httpHandler.Login)
	app.POST("/signup", httpHandler.Signup)
	app.GET("/code", httpHandler.GenerateSignupVerificationCode)
	app.POST("/reset-password", httpHandler.ResetPassword)
	app.POST("/forgot-password", httpHandler.ForgotPassword)
	app.POST("/validate/{token}", httpHandler.ValidateToken)
	app.POST("/verify/{token}", httpHandler.VerifyAccount)
	app.GET("/refresh-token/{token}", httpHandler.RefreshToken)

	app.Start()
}
