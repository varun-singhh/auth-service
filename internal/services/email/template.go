package email

import (
	"bytes"
	"fmt"
	"gofr.dev/pkg/errors"
	"html/template"
)

// EmailType represents the type of email to send.
type EmailType int

const (
	PasswordReset EmailType = iota
	VerificationCode
	Custom
)

// EmailData contains the data needed for composing the email.
type EmailData struct {
	EmailType
	DataMap map[string]string
}

// generateTemplateByType returns an email template of the specified type with the provided data.
func generateTemplateByType(emailData EmailData) ([]byte, error) {
	switch emailData.EmailType {
	case PasswordReset:
		return generateCustomEmailTemplate(resetPasswordEmailTemplate, emailData.DataMap)
	case VerificationCode:
		return generateCustomEmailTemplate(verifyCodeEmailTemplate, emailData.DataMap)
	case Custom:
		return generateCustomEmailTemplate(customEmailTemplate, emailData.DataMap)
	default:
		return nil, &errors.InvalidParam{Param: []string{"emailType"}}
	}
}

func generateCustomEmailTemplate(templateContent string, data map[string]string) ([]byte, error) {
	// Parse the template
	tmpl, err := template.New("emailTemplate").Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("error parsing email template: %w", err)
	}

	// Execute the template with the data
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return nil, fmt.Errorf("error executing email template: %w", err)
	}

	return tpl.Bytes(), nil
}
