package loungeup

import (
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/errors"
)

// DefaultPrivacyPolicy used as the privacy policy applied by default to all entities.
var DefaultPrivacyPolicy = &PrivacyPolicy{
	PortalConsentIsImplicit: true,
	EnabledAt:               time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC),
}

type PrivacyPolicy struct {
	// ID of the privacy policy.
	ID uuid.UUID `json:"id"`

	// EntityID is the ID of the entity the privacy policy is applied to.
	EntityID uuid.UUID `json:"entityId"`

	// PortalConsentContent represents the translations of the content of the portal consent.
	PortalConsentContent Translations `json:"portalConsentContent,omitempty"`

	// PortalConsentIsImplicit represents whether the portal consent is implicit or not.
	PortalConsentIsImplicit bool `json:"portalConsentIsImplicit"`

	// FooterContent represents the translations of the footer content.
	FooterContent Translations `json:"footerContent,omitempty"`

	// ExternalURL represents the translations of the external URL.
	ExternalURL Translations `json:"externalUrl,omitempty"`

	// GuestsRetentionDays represents the number of days the guests are retained.
	GuestsRetentionDays int `json:"guestsRetentionDays,omitempty"`

	// MessagesRetentionDays represents the number of days the messages are retained.
	MessagesRetentionDays int `json:"messagesRetentionDays,omitempty"`

	// EnabledAt represents the date at which the privacy policy is enabled.
	EnabledAt time.Time `json:"enabledAt"`

	CreatedAt time.Time `json:"createdAt"`
	CreatedBy uuid.UUID `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy uuid.UUID `json:"updatedBy"`
}

func (p *PrivacyPolicy) IsDefault() bool { return p.ID == DefaultPrivacyPolicy.ID }

func (p *PrivacyPolicy) Validate() error {
	if p.ID == uuid.Nil {
		return &errors.Error{Code: errors.CodeInternal, Message: "ID must not be empty"}
	}

	if p.EntityID == uuid.Nil {
		return &errors.Error{Code: errors.CodeInternal, Message: "Entity ID must not be empty"}
	}

	if p.EnabledAt.IsZero() {
		return &errors.Error{Code: errors.CodeInternal, Message: "Enabled at must not be empty"}
	}

	if p.CreatedAt.IsZero() {
		return &errors.Error{Code: errors.CodeInternal, Message: "Created at must not be empty"}
	}

	if p.CreatedBy == uuid.Nil {
		return &errors.Error{Code: errors.CodeInternal, Message: "Created by must not be empty"}
	}

	if p.UpdatedAt.IsZero() {
		return &errors.Error{Code: errors.CodeInternal, Message: "Updated at must not be empty"}
	}

	if p.UpdatedBy == uuid.Nil {
		return &errors.Error{Code: errors.CodeInternal, Message: "Updated by must not be empty"}
	}

	return nil
}

type PrivacyPoliciesSelector struct {
	EntityID                 uuid.UUID `json:"entityId"`
	IsGuestsRetentionDaysSet bool      `json:"isGuestsRetentionDaysSet"`
}

type PrivacyPolicySelector struct {
	EntityID uuid.UUID `json:"entityId"`
}

type PrivacyPoliciesReadWriter interface {
	PrivacyPoliciesReader
	PrivacyPoliciesWriter
}

type PrivacyPoliciesReader interface {
	// ReadPrivacyPolicies reads the privacy policies matching the given selector.
	ReadPrivacyPolicies(givenSelector *PrivacyPoliciesSelector) ([]*PrivacyPolicy, error)

	// ReadPrivacyPolicy reads the privacy policy matching the given selector.
	ReadPrivacyPolicy(givenSelector *PrivacyPolicySelector) (*PrivacyPolicy, error)
}

type PrivacyPoliciesWriter interface {
	// UpdatePrivacyPolicy updates the given privacy policy.
	UpdatePrivacyPolicy(privacyPolicyToUpdate *PrivacyPolicy) error

	// DeletePrivacyPolicy deletes the given privacy policy.
	DeletePrivacyPolicy(privacyPolicyToDelete *PrivacyPolicy) error
}
