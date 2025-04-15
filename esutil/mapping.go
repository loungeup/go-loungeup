package esutil

import (
	"strings"
	"sync/atomic"

	"github.com/loungeup/go-loungeup/errors"
)

var globalMappingKeys atomic.Pointer[MappingKeys]

//nolint:gochecknoinits
func init() {
	globalMappingKeys.Store(newMappingKeys())
}

func GlobalMappingKeys() *MappingKeys { return globalMappingKeys.Load() }

type MappingKeys struct {
	Booking            *BookingMappingKeys
	Guest              *GuestMappingKeys
	ComputedAttributes *ComputedAttributesMappingKeys
}

type ComputedAttributesMappingKeys struct {
	Account *ScopedComputedAttributesMappingKeys
	Group   *ScopedComputedAttributesMappingKeys
	Chain   *ScopedComputedAttributesMappingKeys
}

type ScopedComputedAttributesMappingKeys struct {
	Boolean string
	Date    string
	Number  string
	Text    string
}

type ScopedMappingKeys struct {
	Booking            *BookingMappingKeys
	Guest              *ScopedGuestMappingKeys
	ComputedAttributes *ScopedComputedAttributesMappingKeys
}

type BookingMappingKeys struct {
	Arrival                  string
	ArrivalDay               string
	ArrivalTime              string
	Balance                  string
	BookingDate              string
	Channel                  string
	CicoHasCompletedPayment  string
	CicoHasFilledPoliceForm  string
	CicoHasFilledPrestayForm string
	CustomFields             string
	CustomFieldsBoolean      string
	CustomFieldsDate         string
	CustomFieldsList         string
	CustomFieldsNumber       string
	CustomFieldsText         string
	Departure                string
	DepartureDay             string
	EntityID                 string
	Fare                     string
	FareCode                 string
	GuestID                  string
	ID                       string
	Index                    string
	InstayDates              string
	Pass                     string
	PaxAdults                string
	PaxBabies                string
	PaxChildren              string
	PMSBookingID             string
	PMSBookingParentID       string
	Room                     string
	RoomType                 string
	Status                   string
	StayLength               string
	TouristTax               string
	UpdatedAt                string
	Wildcard                 string
}

type GuestMappingKeys struct {
	Account *ScopedGuestMappingKeys
	Chain   *ScopedGuestMappingKeys
	Group   *ScopedGuestMappingKeys
}

type ScopedGuestMappingKeys struct {
	Birthdate           string
	BirthplaceCountry   string
	City                string
	Company             string
	ComposedWith        string
	Country             string
	CustomFields        string
	CustomFieldsBoolean string
	CustomFieldsDate    string
	CustomFieldsList    string
	CustomFieldsNumber  string
	CustomFieldsText    string
	Emails              string
	EntityID            string
	FirstName           string
	Gender              string
	ID                  string
	Languages           string
	LastName            string
	Nationalities       string
	OptedOutMarketing   string
	Phones              string
	PMSID               string
	State               string
	Title               string
	UpdatedAt           string
	Zipcode             string
	Wildcard            string
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

func (scope MappingKeysScope) computedAttributePrefix() string {
	switch scope {
	case MappingKeysScopeAccount:
		return "typedComputedAttributes.account"
	case MappingKeysScopeChain:
		return "typedComputedAttributes.chain"
	case MappingKeysScopeGroup:
		return "typedComputedAttributes.group"
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

func newMappingKeys() *MappingKeys {
	return &MappingKeys{
		Booking: newBookingMappingKeys(),
		Guest: &GuestMappingKeys{
			Account: newScopedGuestMappingKeys(MappingKeysScopeAccount),
			Chain:   newScopedGuestMappingKeys(MappingKeysScopeChain),
			Group:   newScopedGuestMappingKeys(MappingKeysScopeGroup),
		},
		ComputedAttributes: &ComputedAttributesMappingKeys{
			Account: newScopedComputedAttributeMappingKeys(MappingKeysScopeAccount),
			Group:   newScopedComputedAttributeMappingKeys(MappingKeysScopeGroup),
			Chain:   newScopedComputedAttributeMappingKeys(MappingKeysScopeChain),
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
		ComputedAttributes: &ScopedComputedAttributesMappingKeys{
			Boolean: newScopedComputedAttributeMappingKeys(scope).Boolean,
			Date:    newScopedComputedAttributeMappingKeys(scope).Date,
			Number:  newScopedComputedAttributeMappingKeys(scope).Number,
			Text:    newScopedComputedAttributeMappingKeys(scope).Text,
		},
	}, nil
}

func newBookingMappingKeys() *BookingMappingKeys {
	prefix := "booking"

	return &BookingMappingKeys{
		Arrival:                  joinMappingKeyParts(prefix, "arrival"),
		ArrivalDay:               joinMappingKeyParts(prefix, "arrivalDay"),
		ArrivalTime:              joinMappingKeyParts(prefix, "data.arrivalTime"),
		Balance:                  joinMappingKeyParts(prefix, "balance"),
		BookingDate:              joinMappingKeyParts(prefix, "bookingDate"),
		Channel:                  joinMappingKeyParts(prefix, "channel"),
		CicoHasCompletedPayment:  joinMappingKeyParts(prefix, "cico.hasCompletedPayment"),
		CicoHasFilledPoliceForm:  joinMappingKeyParts(prefix, "cico.hasFilledPoliceForm"),
		CicoHasFilledPrestayForm: joinMappingKeyParts(prefix, "cico.hasFilledPrestayForm"),
		CustomFields:             joinMappingKeyParts(prefix, "customFields"),
		CustomFieldsBoolean:      joinMappingKeyParts(prefix, "customFields.boolean"),
		CustomFieldsDate:         joinMappingKeyParts(prefix, "customFields.date"),
		CustomFieldsList:         joinMappingKeyParts(prefix, "customFields.list"),
		CustomFieldsNumber:       joinMappingKeyParts(prefix, "customFields.number"),
		CustomFieldsText:         joinMappingKeyParts(prefix, "customFields.text"),
		Departure:                joinMappingKeyParts(prefix, "departure"),
		DepartureDay:             joinMappingKeyParts(prefix, "departureDay"),
		EntityID:                 joinMappingKeyParts(prefix, "entityId"),
		Fare:                     joinMappingKeyParts(prefix, "fare"),
		FareCode:                 joinMappingKeyParts(prefix, "fareCode"),
		GuestID:                  joinMappingKeyParts(prefix, "guestId"),
		ID:                       joinMappingKeyParts(prefix, "id"),
		Index:                    joinMappingKeyParts(prefix, "data.index"),
		InstayDates:              joinMappingKeyParts(prefix, "instayDates"),
		Pass:                     joinMappingKeyParts(prefix, "pass"),
		PaxAdults:                joinMappingKeyParts(prefix, "paxAdults"),
		PaxBabies:                joinMappingKeyParts(prefix, "paxBabies"),
		PaxChildren:              joinMappingKeyParts(prefix, "paxChildren"),
		PMSBookingID:             joinMappingKeyParts(prefix, "pmsBookingId"),
		PMSBookingParentID:       joinMappingKeyParts(prefix, "pmsBookingParentId"),
		Room:                     joinMappingKeyParts(prefix, "room"),
		RoomType:                 joinMappingKeyParts(prefix, "roomType"),
		Status:                   joinMappingKeyParts(prefix, "status"),
		StayLength:               joinMappingKeyParts(prefix, "stayLength"),
		TouristTax:               joinMappingKeyParts(prefix, "touristTax"),
		UpdatedAt:                joinMappingKeyParts(prefix, "updatedAt"),
		Wildcard:                 joinMappingKeyParts(prefix, "*"),
	}
}

func newScopedGuestMappingKeys(scope MappingKeysScope) *ScopedGuestMappingKeys {
	prefix := scope.guestPrefix()

	return &ScopedGuestMappingKeys{
		Birthdate:           joinMappingKeyParts(prefix, "birthdate"),
		BirthplaceCountry:   joinMappingKeyParts(prefix, "birthplace.country"),
		City:                joinMappingKeyParts(prefix, "city"),
		Company:             joinMappingKeyParts(prefix, "company"),
		ComposedWith:        joinMappingKeyParts(prefix, "composedWith"),
		Country:             joinMappingKeyParts(prefix, "country"),
		CustomFields:        joinMappingKeyParts(prefix, "customFields"),
		CustomFieldsBoolean: joinMappingKeyParts(prefix, "customFields.boolean"),
		CustomFieldsDate:    joinMappingKeyParts(prefix, "customFields.date"),
		CustomFieldsList:    joinMappingKeyParts(prefix, "customFields.list"),
		CustomFieldsNumber:  joinMappingKeyParts(prefix, "customFields.number"),
		CustomFieldsText:    joinMappingKeyParts(prefix, "customFields.text"),
		Emails:              joinMappingKeyParts(prefix, "emails"),
		EntityID:            joinMappingKeyParts(prefix, "entityId"),
		FirstName:           joinMappingKeyParts(prefix, "firstname"),
		Gender:              joinMappingKeyParts(prefix, "gender"),
		ID:                  joinMappingKeyParts(prefix, "id"),
		Languages:           joinMappingKeyParts(prefix, "languages"),
		LastName:            joinMappingKeyParts(prefix, "lastname"),
		Nationalities:       joinMappingKeyParts(prefix, "nationalities"),
		OptedOutMarketing:   joinMappingKeyParts(prefix, "optedOut.marketing"),
		Phones:              joinMappingKeyParts(prefix, "phones"),
		PMSID:               joinMappingKeyParts(prefix, "pmsId"),
		State:               joinMappingKeyParts(prefix, "state"),
		Title:               joinMappingKeyParts(prefix, "title"),
		UpdatedAt:           joinMappingKeyParts(prefix, "updatedAt"),
		Zipcode:             joinMappingKeyParts(prefix, "zipcode"),
		Wildcard:            joinMappingKeyParts(prefix, "*"),
	}
}

func newScopedComputedAttributeMappingKeys(scope MappingKeysScope) *ScopedComputedAttributesMappingKeys {
	prefix := scope.computedAttributePrefix()

	return &ScopedComputedAttributesMappingKeys{
		Boolean: joinMappingKeyParts(prefix, "boolean"),
		Date:    joinMappingKeyParts(prefix, "date"),
		Number:  joinMappingKeyParts(prefix, "number"),
		Text:    joinMappingKeyParts(prefix, "text"),
	}
}

func joinMappingKeyParts(parts ...string) string {
	return strings.Join(parts, ".")
}
