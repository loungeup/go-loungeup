package email

import (
	"testing"

	"github.com/loungeup/go-loungeup/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := map[string]struct {
		modifyInput func(*Email)
		wantError   bool
	}{
		"valid email": {
			modifyInput: func(e *Email) {},
			wantError:   false,
		},
		"missing sender": {
			modifyInput: func(e *Email) { e.Sender = "" },
			wantError:   true,
		},
		"missing recipients": {
			modifyInput: func(e *Email) { e.Recipients = nil },
			wantError:   true,
		},
		"empty recipients": {
			modifyInput: func(e *Email) { e.Recipients = []string{} },
			wantError:   true,
		},
		"missing subject": {
			modifyInput: func(e *Email) { e.Subject = "" },
			wantError:   true,
		},
		"missing body and template": {
			modifyInput: func(e *Email) { e.HTMLBody, e.TextBody, e.TemplateID = nil, nil, nil },
			wantError:   true,
		},
		"body and template": {
			modifyInput: func(e *Email) { e.HTMLBody, e.TextBody, e.TemplateID = []byte("html"), []byte("text"), "template" },
			wantError:   true,
		},
		"invalid attachment": {
			modifyInput: func(e *Email) { e.Attachments = []*Attachment{{Filename: ""}} },
			wantError:   true,
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			in := &Email{
				Sender:     "ops@loungeup.com",
				Recipients: []string{"john@loungeup.com", "jane@loungeup.com"},
				Subject:    "Hello",
				HTMLBody:   []byte("<h1>Hello</h1>"),
				TextBody:   []byte("Hello"),
				Attachments: []*Attachment{
					{
						Filename:    "hello.txt",
						Content:     []byte("Hello"),
						ContentType: "text/plain",
					},
				},
			}
			tt.modifyInput(in)

			if err := in.Validate(); tt.wantError {
				assert.Equal(t, errors.CodeInvalid, errors.ErrorCode(err))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
