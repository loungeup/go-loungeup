package esutil

import (
	"bytes"
	"encoding/json"
)

type GuestBookingDocument struct {
	Booking Booking `json:"booking"`
	GuestCardDocument
}

type ScopedGuestBookingDocument struct {
	Booking                 Booking                        `json:"booking"`
	CampaignStats           json.RawMessage                `json:"campaignStats,omitempty"`
	Device                  json.RawMessage                `json:"device,omitempty"`
	Guest                   *ScopedGuest                   `json:"guest,omitempty"`
	SurveyAnswers           json.RawMessage                `json:"surveyAnswers,omitempty"`
	TypedComputedAttributes *ScopedTypedComputedAttributes `json:"typedComputedAttributes,omitempty"`
	Aggregations            *ScopedAggregations            `json:"aggregations,omitempty"`
}

type GuestCardDocument struct {
	CampaignStats           json.RawMessage          `json:"campaignStats,omitempty"`
	Device                  json.RawMessage          `json:"device,omitempty"`
	Guest                   Guest                    `json:"guest"`
	SurveyAnswers           json.RawMessage          `json:"surveyAnswers,omitempty"`
	TypedComputedAttributes *TypedComputedAttributes `json:"typedComputedAttributes,omitempty"`
	Aggregations            Aggregations             `json:"aggregations,omitempty"`
}

type ScopedGuestCardDocument struct {
	CampaignStats           json.RawMessage                `json:"campaignStats,omitempty"`
	Device                  json.RawMessage                `json:"device,omitempty"`
	Guest                   *ScopedGuest                   `json:"guest,omitempty"`
	SurveyAnswers           json.RawMessage                `json:"surveyAnswers,omitempty"`
	TypedComputedAttributes *ScopedTypedComputedAttributes `json:"typedComputedAttributes,omitempty"`
	Aggregations            *ScopedAggregations            `json:"aggregations,omitempty"`
}

type TypedComputedAttributes struct {
	Account *ScopedTypedComputedAttributes `json:"account,omitempty"`
	Chain   *ScopedTypedComputedAttributes `json:"chain,omitempty"`
	Group   *ScopedTypedComputedAttributes `json:"group,omitempty"`
}

type ScopedTypedComputedAttributes struct {
	Number  []ScopedTypedComputedAttribute `json:"number,omitempty"`
	Text    []ScopedTypedComputedAttribute `json:"text,omitempty"`
	Boolean []ScopedTypedComputedAttribute `json:"boolean,omitempty"`
	Date    []ScopedTypedComputedAttribute `json:"date,omitempty"`
}

type ScopedTypedComputedAttribute struct {
	ID        string `json:"id,omitempty"`
	Value     any    `json:"value,omitempty"`
	AccountID string `json:"accountId,omitempty"`
}

type ScopedAggregations struct {
	CounterFutureBookings float64                     `json:"counterFutureBookings"`
	LastDeparture         string                      `json:"lastDeparture,omitempty"`
	CounterBookings       float64                     `json:"counterBookings"`
	EntityID              string                      `json:"entityId,omitempty"`
	ConvertedAvgFare      *BookingConvertedCurrencies `json:"convertedAvgFare,omitempty"`
	CounterPastBookings   float64                     `json:"counterPastBookings"`
	AvgFares              float64                     `json:"avgFares"`
	SumFares              float64                     `json:"sumFares"`
	ConvertedSumFare      *BookingConvertedCurrencies `json:"convertedSumFare,omitempty"`
	NextArrival           string                      `json:"nextArrival,omitempty"`
	AccountIDs            Array[string]               `json:"accountIds,omitempty"`
	AccountIDsCounter     float64                     `json:"accountIdsCounter"`
	LastAccountID         string                      `json:"lastAccountId,omitempty"`
	ID                    string                      `json:"id,omitempty"`
	NextAccountID         string                      `json:"nextAccountId,omitempty"`
	UpdatedAt             string                      `json:"updatedAt,omitempty"`
}

type Aggregations struct {
	Account *ScopedAggregations `json:"account,omitempty"`
	Chain   *ScopedAggregations `json:"chain,omitempty"`
	Group   *ScopedAggregations `json:"group,omitempty"`
}

type Booking struct {
	AggregateAt           string                      `json:"aggregateAt,omitempty"`
	Arrival               string                      `json:"arrival,omitempty"`
	ArrivalDay            string                      `json:"arrivalDay,omitempty"`
	ArrivalDow            string                      `json:"arrivalDow,omitempty"`
	Balance               any                         `json:"balance,omitempty"`
	BookingDate           string                      `json:"bookingDate,omitempty"`
	BookingDateArrival    string                      `json:"bookingDateArrival,omitempty"`
	Channel               string                      `json:"channel,omitempty"`
	Cico                  *BookingCico                `json:"cico,omitempty"`
	Closed                bool                        `json:"closed,omitempty"`
	ConvertedBalance      *BookingConvertedCurrencies `json:"convertedBalance,omitempty"`
	ConvertedFare         *BookingConvertedCurrencies `json:"convertedFare,omitempty"`
	ConvertedFarePerNight *BookingConvertedCurrencies `json:"convertedFarePerNight,omitempty"`
	ConvertedTouristTax   *BookingConvertedCurrencies `json:"convertedTouristTax,omitempty"`
	CustomFields          *CustomFields               `json:"customFields,omitempty"`
	Data                  *BookingData                `json:"data,omitempty"`
	Departure             string                      `json:"departure,omitempty"`
	DepartureDay          string                      `json:"departureDay,omitempty"`
	DepartureDow          string                      `json:"departureDow,omitempty"`
	EntityID              string                      `json:"entityId,omitempty"`
	ExternalIDs           *BookingExternalIDs         `json:"externalIds,omitempty"`
	Fare                  any                         `json:"fare,omitempty"`
	FareCode              string                      `json:"fareCode,omitempty"`
	FarePerNight          float64                     `json:"farePerNight,omitempty"`
	FilledCustomFields    Array[string]               `json:"filledCustomFields,omitempty"`
	GuestID               string                      `json:"guestId,omitempty"`
	ID                    int                         `json:"id,omitempty"`
	InstayDates           Array[string]               `json:"instayDates,omitempty"`
	InstayDows            Array[string]               `json:"instayDows,omitempty"`
	Last                  string                      `json:"last,omitempty"`
	Partner               string                      `json:"partner,omitempty"`
	Pass                  string                      `json:"pass,omitempty"`
	PaxAdults             any                         `json:"paxAdults,omitempty"`
	PaxBabies             any                         `json:"paxBabies,omitempty"`
	PaxChildren           any                         `json:"paxChildren,omitempty"`
	PMSBookingID          string                      `json:"pmsBookingId,omitempty"`
	PMSBookingParentID    string                      `json:"pmsBookingParentId,omitempty"`
	Purposes              Array[string]               `json:"purposes,omitempty"`
	Room                  string                      `json:"room,omitempty"`
	RoomType              string                      `json:"roomType,omitempty"`
	Start                 string                      `json:"start,omitempty"`
	Status                string                      `json:"status,omitempty"`
	StayLength            int                         `json:"stayLength,omitempty"`
	Tags                  Array[string]               `json:"tags,omitempty"`
	TouristTax            any                         `json:"touristTax,omitempty"`
	UpdatedAt             string                      `json:"updatedAt,omitempty"`
	Weekend               bool                        `json:"weekend,omitempty"`
}

var _ json.Unmarshaler = (*BookingConvertedCurrencies)(nil)

// UnmarshalJSON implements the json.Unmarshaler interface.
// This method is required to be compatible with the Guest Profile PHP server.
func (m *BookingConvertedCurrencies) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, json.RawMessage(`[]`)) {
		*m = BookingConvertedCurrencies{}

		return nil
	}

	type Alias BookingConvertedCurrencies

	alias := Alias{}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	*m = BookingConvertedCurrencies(alias)

	return nil
}

type BookingExternalIDs struct {
	ExternalID       string `json:"externalid,omitempty"`
	MonewebAccountID string `json:"moneweb_account_id,omitempty"`
	Qualitelis       string `json:"qualitelis,omitempty"`
}

type BookingCico struct {
	HasCompletedPayment  bool `json:"hasCompletedPayment,omitempty"`
	HasFilledPoliceForm  bool `json:"hasFilledPoliceForm,omitempty"`
	HasFilledPrestayForm bool `json:"hasFilledPrestayForm,omitempty"`
}

type BookingData struct {
	App             string `json:"app,omitempty"`
	ArrivalTime     string `json:"arrivalTime,omitempty"`
	Converted       string `json:"converted,omitempty"`
	Index           string `json:"index,omitempty"`
	NextStay        string `json:"nextStay,omitempty"`
	PmsCreatedAt    string `json:"pmsCreatedAt,omitempty"`
	PmsImportAt     string `json:"pmsImportAt,omitempty"`
	PrevStay        string `json:"prevStay,omitempty"`
	PreviousStatus  string `json:"previousStatus,omitempty"`
	ReindexGuest    string `json:"reindexGuest,omitempty"`
	StatusUpdatedAt string `json:"statusUpdatedAt,omitempty"`
}

type BookingConvertedCurrencies struct {
	AUD float64 `json:"aud,omitempty"`
	CAD float64 `json:"cad,omitempty"`
	CHF float64 `json:"chf,omitempty"`
	CNY float64 `json:"cny,omitempty"`
	EUR float64 `json:"eur,omitempty"`
	GBP float64 `json:"gbp,omitempty"`
	JPY float64 `json:"jpy,omitempty"`
	KRW float64 `json:"krw,omitempty"`
	SGD float64 `json:"sgd,omitempty"`
	USD float64 `json:"usd,omitempty"`
}

type BookingCustomFields struct{}

type Guest struct {
	ID        string `json:"id,omitempty"`
	EntityID  string `json:"entityId,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	IndexedAt string `json:"indexedAt,omitempty"`

	Account *ScopedGuest `json:"account,omitempty"`
	Chain   *ScopedGuest `json:"chain,omitempty"`
	Group   *ScopedGuest `json:"group,omitempty"`
}

// ScopedGuest represents a guest in an Guest format. It is used to add the
// representations based on the entities of the guest.
type ScopedGuest struct {
	AnonymizedAt        string                  `json:"anonymizedAt,omitempty"`
	Anonymous           bool                    `json:"anonymous"`
	Birthdate           string                  `json:"birthdate,omitempty"`
	Birthplace          *GuestAddress           `json:"birthplace,omitempty"`
	City                Array[string]           `json:"city,omitempty"`
	Company             string                  `json:"company,omitempty"`
	ComposedWith        Array[string]           `json:"composedWith,omitempty"`
	ComputeAggregations bool                    `json:"computeAggregations,omitempty"`
	Country             Array[string]           `json:"country,omitempty"`
	CreatedAt           string                  `json:"createdAt,omitempty"`
	CustomFields        *CustomFields           `json:"customFields,omitempty"`
	Documents           *GuestDocuments         `json:"documents,omitempty"`
	EmailDomains        Array[string]           `json:"emailDomains,omitempty"`
	Emails              Array[string]           `json:"emails,omitempty"`
	EmailsMergeableAt   string                  `json:"emailsMergeableAt,omitempty"`
	EntityID            string                  `json:"entityId,omitempty"`
	FilledCustomFields  Array[string]           `json:"filledCustomFields,omitempty"`
	Firstname           string                  `json:"firstname,omitempty"`
	Gender              string                  `json:"gender,omitempty"`
	HasPersonnalEmail   bool                    `json:"hasPersonnalEmail,omitempty"`
	ID                  string                  `json:"id,omitempty"`
	Languages           Array[string]           `json:"languages,omitempty"`
	LastChannel         string                  `json:"lastChannel,omitempty"`
	Lastname            string                  `json:"lastname,omitempty"`
	Loyalty             *GuestLoyalty           `json:"loyalty,omitempty"`
	MessengerID         *GuestMessengerID       `json:"messengerId,omitempty"`
	Metadata            *GuestMetadata          `json:"metadata,omitempty"`
	Middlename          string                  `json:"middlename,omitempty"`
	Nationalities       Array[string]           `json:"nationalities,omitempty"`
	Notes               Array[string]           `json:"notes,omitempty"`
	OptedOut            *GuestOptedOut          `json:"optedOut,omitempty"`
	Phones              Array[string]           `json:"phones,omitempty"`
	PhonesMergeableAt   string                  `json:"phonesMergeableAt,omitempty"`
	PMSID               Array[string]           `json:"pmsId,omitempty"`
	PreferredContacts   *GuestPreferredContacts `json:"preferredContacts,omitempty"`
	Socials             *GuestSocials           `json:"socials,omitempty"`
	Source              Array[string]           `json:"source,omitempty"`
	State               Array[string]           `json:"state,omitempty"`
	Street              Array[string]           `json:"street,omitempty"`
	SubdivisionCode     Array[string]           `json:"subdivisionCode,omitempty"`
	Tags                Array[string]           `json:"tags,omitempty"`
	Timezone            Array[string]           `json:"timezone,omitempty"`
	Title               string                  `json:"title,omitempty"`
	TrustableContacts   *GuestTrustableContacts `json:"trustableContacts,omitempty"`
	UpdatedAt           string                  `json:"updatedAt,omitempty"`
	ZipCode             Array[string]           `json:"zipcode,omitempty"`
}

type GuestAddress struct {
	City            string `json:"city,omitempty"`
	Country         string `json:"country,omitempty"`
	State           string `json:"state,omitempty"`
	Street          string `json:"street,omitempty"`
	SubdivisionCode string `json:"subdivisionCode,omitempty"`
	Timezone        string `json:"timezone,omitempty"`
	Zipcode         string `json:"zipcode,omitempty"`
}

type CustomFields struct {
	Boolean []*GuestCustomField `json:"boolean,omitempty"`
	Date    []*GuestCustomField `json:"date,omitempty"`
	List    []*GuestCustomField `json:"list,omitempty"`
	Number  []*GuestCustomField `json:"number,omitempty"`
	Text    []*GuestCustomField `json:"text,omitempty"`
}

type GuestCustomField struct {
	Key   string `json:"key,omitempty"`
	Value any    `json:"value,omitempty"`
}

type GuestDocuments struct {
	DrivingLicenceID string `json:"drivingLicenceId,omitempty"`
	IdentityCardID   string `json:"identityCardId,omitempty"`
	PassportID       string `json:"passportId,omitempty"`
}

type GuestMessengerID struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type GuestLoyalty struct {
	ID      string `json:"id,omitempty"`
	Points  int    `json:"points,omitempty"`
	Program string `json:"program,omitempty"`
	Status  string `json:"status,omitempty"`
}

type GuestMetadata struct {
	Filename []string `json:"filename,omitempty"`
}

type GuestOptedOut struct {
	App             bool `json:"app,omitempty"`
	Loyalty         bool `json:"loyalty,omitempty"`
	Marketing       bool `json:"marketing,omitempty"`
	CustomerAccount bool `json:"customeraccount,omitempty"`
	Prestay         bool `json:"prestay,omitempty"`
	Sendinblue      bool `json:"sendinblue,omitempty"`
}

type GuestPreferredContacts struct {
	Address     *GuestAddress `json:"address,omitempty"`
	Email       string        `json:"email,omitempty"`
	Language    string        `json:"language,omitempty"`
	Nationality string        `json:"nationality,omitempty"`
	Phone       string        `json:"phone,omitempty"`
}

type GuestPrivacyPolicy struct {
	PolicyID string `json:"policyId,omitempty"`
}

type GuestSocials struct {
	Avatar            string `json:"avatar,omitempty"`
	Facebook          bool   `json:"facebook,omitempty"`
	Linkedin          bool   `json:"linkedin,omitempty"`
	LinkedinFollowers int    `json:"linkedinFollowers,omitempty"`
	Twitter           bool   `json:"twitter,omitempty"`
	TwitterFollowers  int    `json:"twitterFollowers,omitempty"`
}

type GuestTrustableContacts struct {
	Emails []string `json:"emails,omitempty"`
	Phones []string `json:"phones,omitempty"`
}
