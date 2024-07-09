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

	app.POST("/api/login", httpHandler.Login)
	app.POST("/api/signup", httpHandler.Signup)
	app.GET("/api/code", httpHandler.GenerateSignupVerificationCode)
	app.POST("/api/reset-password", httpHandler.ResetPassword)
	app.POST("/api/forgot-password", httpHandler.ForgotPassword)
	app.GET("/api/validate/{token}", httpHandler.ValidateToken)
	app.POST("/api/verify", httpHandler.VerifyAccount)
	app.GET("/api/refresh-token/{token}", httpHandler.RefreshToken)

	app.Start()
}
