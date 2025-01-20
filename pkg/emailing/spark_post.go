package emailing

import (
	"encoding/base64"
	"fmt"

	"github.com/SparkPost/gosparkpost"
	"github.com/loungeup/go-loungeup/pkg/errors"
)

type sparkPostClient struct{ baseClient *gosparkpost.Client }

func NewSparkPostClient(apiKey string) (*sparkPostClient, error) {
	const (
		defaultBaseURL    = "https://api.sparkpost.com"
		defaultAPIVersion = 1
	)

	baseClient := &gosparkpost.Client{}
	if err := baseClient.Init(&gosparkpost.Config{
		BaseUrl:    defaultBaseURL,
		ApiKey:     apiKey,
		ApiVersion: defaultAPIVersion,
	}); err != nil {
		return nil, fmt.Errorf("could not initialize SparkPost client: %w", err)
	}

	return &sparkPostClient{baseClient}, nil
}

var _ (Sender) = (*sparkPostClient)(nil)

func (c *sparkPostClient) Send(email *Email) error {
	if err := email.Validate(); err != nil {
		return &errors.Error{Code: errors.CodeInvalid, Message: "Invalid email", UnderlyingError: err}
	}

	if _, _, err := c.baseClient.Send(makeSparkPostTransmission(email)); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	return nil
}

func makeSparkPostTransmission(email *Email) *gosparkpost.Transmission {
	if email.hasTemplate() {
		return &gosparkpost.Transmission{
			Recipients: mapSparkPostRecipients(email.Recipients),
			Content: map[string]any{
				"template_id": email.TemplateID,
			},
			SubstitutionData: email.TemplateVariables,
		}
	}

	return &gosparkpost.Transmission{
		Recipients: mapSparkPostRecipients(email.Recipients),
		Content: gosparkpost.Content{
			From:        email.Sender,
			Subject:     email.Subject,
			HTML:        string(email.HTMLBody),
			Text:        string(email.TextBody),
			Attachments: mapSparkPostAttachments(email.Attachments),
		},
	}
}

func mapSparkPostAttachments(attachments []*Attachment) []gosparkpost.Attachment {
	result := []gosparkpost.Attachment{}
	for _, attachment := range attachments {
		result = append(result, gosparkpost.Attachment{
			Filename: attachment.Filename,
			MIMEType: attachment.ContentType,
			B64Data:  base64.StdEncoding.EncodeToString(attachment.Content),
		})
	}

	return result
}

func mapSparkPostRecipients(recipients []string) []gosparkpost.Recipient {
	result := []gosparkpost.Recipient{}
	for _, recipient := range recipients {
		result = append(result, gosparkpost.Recipient{
			Address: gosparkpost.Address{
				Email: recipient,
			},
		})
	}

	return result
}
