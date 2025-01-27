package esutil

import (
	"strings"

	"github.com/loungeup/go-loungeup/pkg/errors"
)

type MappingKeys struct {
	Booking *BookingMappingKeys
	Guest   *GuestMappingKeys
}

type ScopedMappingKeys struct {
	Booking *BookingMappingKeys
	Guest   *ScopedGuestMappingKeys
}

type BookingMappingKeys struct {
	ID                  string
	EntityID            string
	Arrival             string
	Departure           string
	PMSBookingID        string
	PMSBookingIDKeyword string
}

type GuestMappingKeys struct {
	Account *ScopedGuestMappingKeys
	Chain   *ScopedGuestMappingKeys
	Group   *ScopedGuestMappingKeys
}

type ScopedGuestMappingKeys struct {
	ID           string
	EntityID     string
	ComposedWith string
	FirstName    string
	LastName     string
	Emails       string
}

type MappingKeysScope int

const (
	MappingKeysScopeUnknown MappingKeysScope = iota
	MappingKeysScopeAccount
	MappingKeysScopeChain
	MappingKeysScopeGroup
)

func NewMappingKeysScope(value string) MappingKeysScope {
	switch {
	case strings.EqualFold(value, "account"):
		return MappingKeysScopeAccount
	case strings.EqualFold(value, "chain"):
		return MappingKeysScopeChain
	case strings.EqualFold(value, "group"):
		return MappingKeysScopeGroup
	default:
		return MappingKeysScopeUnknown
	}
}

func (scope MappingKeysScope) guestPrefix() string {
	switch scope {
	case MappingKeysScopeAccount:
		return "guest.account"
	case MappingKeysScopeChain:
		return "guest.chain"
	case MappingKeysScopeGroup:
		return "guest.group"
	default:
		return ""
	}
}

func (scope MappingKeysScope) validate() error {
	switch scope {
	case MappingKeysScopeAccount, MappingKeysScopeChain, MappingKeysScopeGroup:
		return nil
	default:
		return &errors.Error{Code: errors.CodeInvalid, Message: "Invalid ES mapping keys scope"}
	}
}

func NewMappingKeys() *MappingKeys {
	return &MappingKeys{
		Booking: newBookingMappingKeys(),
		Guest: &GuestMappingKeys{
			Account: newScopedGuestMappingKeys(MappingKeysScopeAccount),
			Chain:   newScopedGuestMappingKeys(MappingKeysScopeChain),
			Group:   newScopedGuestMappingKeys(MappingKeysScopeGroup),
		},
	}
}

func NewScopedMappingKeys(scope MappingKeysScope) (*ScopedMappingKeys, error) {
	if err := scope.validate(); err != nil {
		return nil, err
	}

	return &ScopedMappingKeys{
		Booking: newBookingMappingKeys(),
		Guest:   newScopedGuestMappingKeys(scope),
	}, nil
}

func newBookingMappingKeys() *BookingMappingKeys {
	prefix := "booking"

	return &BookingMappingKeys{
		ID:                  prefix + ".id",
		EntityID:            prefix + ".entityId",
		Arrival:             prefix + ".arrival",
		Departure:           prefix + ".departure",
		PMSBookingID:        prefix + ".pmsBookingId",
		PMSBookingIDKeyword: prefix + ".pmsBookingId.keyword",
	}
}

func newScopedGuestMappingKeys(scope MappingKeysScope) *ScopedGuestMappingKeys {
	prefix := scope.guestPrefix()

	return &ScopedGuestMappingKeys{
		ID:           prefix + ".id",
		EntityID:     prefix + ".entityId",
		ComposedWith: prefix + ".composedWith",
		FirstName:    prefix + ".firstname",
		LastName:     prefix + ".lastname",
		Emails:       prefix + ".emails",
	}
}
