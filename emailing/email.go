package emailing

import (
	"github.com/loungeup/go-loungeup/errors"
)

type Email struct {
	Sender     string
	Recipients []string
	Subject    string

	HTMLBody    []byte
	TextBody    []byte
	Attachments []*Attachment

	TemplateID        any
	TemplateVariables map[string]any
}

func (e *Email) Validate() error {
	if e.Sender == "" {
		return &errors.Error{Code: errors.CodeInvalid, Message: "Sender must be set"}
	}

	if len(e.Recipients) == 0 {
		return &errors.Error{Code: errors.CodeInvalid, Message: "Recipients must be set"}
	}

	if e.Subject == "" {
		return &errors.Error{Code: errors.CodeInvalid, Message: "Subject must be set"}
	}

	if !e.hasBody() && !e.hasTemplate() {
		return &errors.Error{Code: errors.CodeInvalid, Message: "Body or template must be set"}
	}

	if e.hasBody() && e.hasTemplate() {
		return &errors.Error{Code: errors.CodeInvalid, Message: "Body and template cannot be set at the same time"}
	}

	for _, attachment := range e.Attachments {
		if err := attachment.Validate(); err != nil {
			return &errors.Error{Code: errors.CodeInvalid, Message: "Invalid attachment", UnderlyingError: err}
		}
	}

	return nil
}

func (e *Email) hasBody() bool     { return len(e.HTMLBody) > 0 || len(e.TextBody) > 0 }
func (e *Email) hasTemplate() bool { return e.TemplateID != nil }

type Attachment struct {
	Filename    string
	Content     []byte
	ContentType string
}

func (a *Attachment) Validate() error {
	if a.Filename == "" {
		return &errors.Error{Code: errors.CodeInvalid, Message: "Filename must be set"}
	}

	if len(a.Content) == 0 {
		return &errors.Error{Code: errors.CodeInvalid, Message: "Content must be set"}
	}

	if a.ContentType == "" {
		return &errors.Error{Code: errors.CodeInvalid, Message: "Content type must be set"}
	}

	return nil
}
