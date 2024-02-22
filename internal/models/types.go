package models

import "net/smtp"

type User struct {
	ID         int    `json:"id"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Password   string `json:"password"`
	Permission string `json:"permission"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

type ForgotPassword struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type ResetPassword struct {
	CurrentPassword   string `json:"current_password"`
	NewPassword       string `json:"new_password"`
	VerificationToken string `json:"verification_token"`
}

type EmailSender struct {
	Auth   smtp.Auth
	config *EmailConfig
}

type EmailConfig struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FromEmail string `json:"from_email"`
	FromName  string `json:"from_name"`
	SmtpAuth  string `json:"smtp_auth"`
	IsSmtp    string `json:"is_smtp"`
}

type VerifyAccount struct {
	Code       string `json:"code"`
	Permission string `json:"permission"`
}
