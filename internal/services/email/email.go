package email

import (
	"crypto/tls"
	"fmt"
	"github.com/varun-singhh/auth-service/internal/models"
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

type Sender struct {
	Auth   smtp.Auth
	config *models.EmailConfig
}

func NewSender(c *models.EmailConfig) (*Sender, error) {
	if c != nil && (c.Host == "" || c.Port == "") {
		return nil, &errors.Response{StatusCode: http.StatusBadRequest, Code: http.StatusText(http.StatusBadRequest), Reason: "missing smtp host or port"}
	}

	return &Sender{config: c}, nil
}

func (s *Sender) Send(ctx *gofr.Context, to, cc, bcc []string, emailData EmailData) error {
	tlsConfig := &tls.Config{
		ServerName: s.config.Host,
	}

	var (
		c   *smtp.Client
		err error
	)

	// Establish SMTP connection
	if s.config.IsSmtp == "true" {
		var conn *tls.Conn
		conn, err = tls.Dial("tcp", s.config.Host+":"+s.config.Port, tlsConfig)
		if err != nil {
			return fmt.Errorf("error creating TLS connection: %w", err)
		}
		defer conn.Close()
		c, err = smtp.NewClient(conn, s.config.Host)
	} else {
		c, err = smtp.Dial(s.config.Host + ":" + s.config.Port)
	}
	if err != nil {
		return fmt.Errorf("error creating SMTP client: %w", err)
	}
	defer c.Close()

	// Start TLS if not already using it
	if s.config.IsSmtp != "true" {
		err = c.StartTLS(tlsConfig)
		if err != nil {
			return fmt.Errorf("error starting TLS: %w", err)
		}
	}

	// Authenticate
	switch s.config.SmtpAuth {
	case `CRAM-MD5`:
		s.Auth = smtp.CRAMMD5Auth(os.Getenv("DEFAULT_SMTP_USERNAME"), os.Getenv("DEFAULT_SMTP_PASSWORD"))
	default:
		s.Auth = smtp.PlainAuth("", ctx.Config.Get("DEFAULT_SMTP_USERNAME"), ctx.Config.Get("DEFAULT_SMTP_PASSWORD"), ctx.Config.Get("DEFAULT_SMTP_HOST"))
	}
	if err := c.Auth(s.Auth); err != nil {
		return fmt.Errorf("authentication error: %w", err)
	}

	// Generate HTML content for the email
	parsedHTMLContent, err := generateTemplateByType(emailData)
	if err != nil {
		fmt.Println("Error generating email HTML:", err)
	}

	// Compose email message with headers
	message := []byte("From: " + s.config.FromEmail + "\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Cc: " + strings.Join(cc, ",") + "\r\n" +
		"Bcc: " + strings.Join(bcc, ",") + "\r\n" +
		"Subject: Your Subject\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"utf-8\"\r\n" +
		"\r\n" +
		string(parsedHTMLContent))

	// Send email
	if err := s.emailWriteAndClose(c, s.config.FromEmail, to, cc, bcc, message); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	ctx.Logger.Infof("email sent to %+v", to)
	return nil
}

func (s *Sender) emailWriteAndClose(c *smtp.Client, from string, to, cc, bcc []string, msg []byte) error {
	err := c.Mail(from)
	if err != nil {
		return err
	}

	// SMTP RCPT to addresses
	for _, addr := range to {
		err = c.Rcpt(addr)
		if err != nil {
			return err
		}
	}

	if len(cc) != 0 {
		for _, addr := range cc {
			err = c.Rcpt(addr)
			if err != nil {
				return err
			}
		}
	}

	if len(bcc) != 0 {
		for _, addr := range bcc {
			err = c.Rcpt(addr)
			if err != nil {
				return err
			}
		}
	}
	// SMTP DATA
	w, err := c.Data()
	if err != nil {
		return err
	}

	// SMTP Write
	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	// SMTP Close Connection
	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}

	return nil
}

func getHTMLForSignupCode(code string) []byte {
	var html = []byte(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Email Verification Code</title>
<style>
    body {
        font-family: Arial, sans-serif;
        background-color: #f4f4f4;
        margin: 0;
        padding: 0;
    }
    .container {
        max-width: 600px;
        margin: 20px auto;
        background-color: #fff;
        padding: 20px;
        border-radius: 10px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    }
    h2 {
        color: #333;
    }
    p {
        color: #666;
        line-height: 1.6;
    }
    .verification-code {
        font-size: 24px;
        font-weight: bold;
        padding: 10px 20px;
        color: #808080;
        border-radius: 5px;
        margin-top: 20px;
    }
</style>
</head>
<body>
<div class="container">
    <h2>Email Verification Code</h2>
    <p>Dear User,</p>
    <p>Your verification code is:</p>
    <div class="verification-code">` + code + `</div>
        <p class="expires">This code expires in 2 minutes. Please use it promptly.</p>
    <p>Please use this code to verify your email address.</p>
    <p>If you didn't request this verification code, you can safely ignore this email.</p>
    <p>Thank you,<br>AIMSS Chamiana Shimla Himachal Pradesh<br>Phone: 01773501627,01773501628<br>Website: <a href="http://www.aimsschamiana.edu.in/">aimsschamiana.edu.in</a></p>
    
</div>
</body>
</html>`)

	return html
}
