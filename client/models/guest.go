package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	GenderFemale string = "F"
	GenderMale   string = "M"
	GenderX      string = "X"

	LandlinePhone  string = "landline"
	MobilePhone    string = "mobile"
	PhoneTypeOther string = "other"

	TitleMR    string = "MR"
	TitleMRMRS string = "MR&MRS"
	TitleMRS   string = "MRS"
	TitleMS    string = "MS"
)

type Guest struct {
	ID       uuid.UUID `json:"id,omitempty"`
	EntityID uuid.UUID `json:"entityId,omitempty"`

	Type string `json:"type,omitempty"`

	Title         DataValue[StructuredValue[string]]      `json:"title,omitempty"`
	Gender        DataValue[StructuredValue[string]]      `json:"gender,omitempty"`
	Firstname     DataValue[StructuredValue[string]]      `json:"firstname,omitempty"`
	Lastname      DataValue[StructuredValue[string]]      `json:"lastname,omitempty"`
	Languages     DataValue[StructuredValueSlice[string]] `json:"languages,omitempty"`
	Nationalities DataValue[StructuredValueSlice[string]] `json:"nationalities,omitempty"`

	Company DataValue[StructuredValue[Company]]    `json:"company,omitempty"`
	Emails  DataValue[StructuredValueSlice[Email]] `json:"emails,omitempty"`
	Phones  DataValue[StructuredPhoneValueSlice]   `json:"phones,omitempty"`
	Socials DataValue[StructuredValueMap[string]]  `json:"socials,omitempty"`

	Addresses  DataValue[StructuredValueSlice[Address]] `json:"addresses,omitempty"`
	Birthdate  DataValue[StructuredValue[BirthDate]]    `json:"birthdate,omitempty"`
	Birthplace DataValue[StructuredValue[BirthPlace]]   `json:"birthplace,omitempty"`

	EmailsMergeableAt time.Time              `json:"emailsMergeableAt,omitempty"`
	PhonesMergeableAt time.Time              `json:"phonesMergeableAt,omitempty"`
	HardMergeAt       time.Time              `json:"hardMergeAt,omitempty"`
	ComposedWith      DataValue[[]uuid.UUID] `json:"composedWith,omitempty"`

	PrivacyPolicies DataValue[StructuredValueSlice[PrivacyPolicy]] `json:"privacyPolicies,omitempty"`
	Credentials     DataValue[StructuredValue[Credentials]]        `json:"credentials,omitempty"`
	CustomFields    DataValue[StructuredValueMap[any]]             `json:"customFields,omitempty"`
	Documents       DataValue[Documents]                           `json:"documents,omitempty"`
	LastChannel     DataValue[StructuredValue[string]]             `json:"lastChannel,omitempty"`
	Loyalty         DataValue[StructuredValue[Loyalty]]            `json:"loyalty,omitempty"`
	Notes           DataValue[StructuredValueSlice[string]]        `json:"notes,omitempty"`
	OptedOut        DataValue[StructuredValueMap[bool]]            `json:"optedOut,omitempty"`
	Revenue         float64                                        `json:"revenue,omitempty"`
	Sources         DataValue[[]Source]                            `json:"sources,omitempty"`
	Tags            DataValue[[]string]                            `json:"tags,omitempty"`

	Extra       uuid.UUID                               `json:"extra,omitempty"`
	MessengerID DataValue[StructuredValue[MessengerID]] `json:"messengerId,omitempty"`
	PMSID       DataValue[StructuredValue[PMSID]]       `json:"pmsId,omitempty"`
	PMSIDs      DataValue[StructuredValueSlice[PMSID]]  `json:"pmsIds,omitempty"`

	AnonymizedAt time.Time `json:"anonymizedAt,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
}

type SearchGuestsRequest struct {
	Query *SearchGuestsQuery `json:"query"`
}

type SearchGuestsQuery struct {
	Logic    SearchGuestsLogic       `json:"logic"`
	Criteria []*SearchGuestsCriteria `json:"criteria"`
}

type SearchGuestsCriteria struct {
	Logic    SearchGuestsLogic          `json:"logic"`
	Criteria []*SearchGuestsSubCriteria `json:"criteria"`
}

type SearchGuestsSubCriteria struct {
	Logic    SearchGuestsLogic          `json:"logic,omitempty"`
	Criteria []*SearchGuestsSubCriteria `json:"criteria,omitempty"`
	Field    string                     `json:"field,omitempty"`
	Operator SearchGuestOperator        `json:"operator,omitempty"`
	Value    any                        `json:"value,omitempty"`
}

type SearchGuestsLogic string

const (
	SearchGuestsLogicAnd SearchGuestsLogic = "AND"
	SearchGuestsLogicOr  SearchGuestsLogic = "OR"
)

type SearchGuestOperator string

const SearchGuestOperatorEqual SearchGuestOperator = "="

type CountGuestsResponse struct {
	Count    int64 `json:"count"`
	Accurate bool  `json:"accurate"`
}

type StructuredValue[T any] struct {
	Preferred bool      `json:"preferred"`
	Value     T         `json:"value"`
	UpdatedAt time.Time `json:"updated"`
}

type StructuredValueMap[T any] map[string]StructuredValue[T]

type StructuredValueSlice[T any] []StructuredValue[T]

func (s StructuredValueSlice[T]) IsZero() bool { return len(s) == 0 }

// Select the preferred value or the last one.
func (s StructuredValueSlice[T]) Select() StructuredValue[T] {
	result := StructuredValue[T]{}

	for _, v := range s {
		if v.Preferred {
			return v
		}

		result = v
	}

	return result
}

// MostRecent value of the slice.
func (s StructuredValueSlice[T]) MostRecent() StructuredValue[T] {
	if len(s) < 1 {
		return StructuredValue[T]{}
	}

	result := s[0]
	for _, v := range s {
		if v.UpdatedAt.After(result.UpdatedAt) {
			result = v
		}
	}

	return result
}

type StructuredPhoneValueSlice StructuredValueSlice[Phone]

func (s StructuredPhoneValueSlice) byType(t string) StructuredValue[Phone] {
	result := StructuredValue[Phone]{}

	for _, p := range s {
		if p.Value.Type != t {
			continue
		}

		if p.Preferred {
			return p // Early-return preferred value.
		}

		result = p
	}

	return result
}

func (s StructuredPhoneValueSlice) Mobile() StructuredValue[Phone]   { return s.byType(MobilePhone) }
func (s StructuredPhoneValueSlice) Landline() StructuredValue[Phone] { return s.byType(LandlinePhone) }

func (s StructuredPhoneValueSlice) IsZero() bool { return StructuredValueSlice[Phone](s).IsZero() }

func (s StructuredPhoneValueSlice) Select() StructuredValue[Phone] {
	return StructuredValueSlice[Phone](s).Select()
}

type Email struct {
	Email string `json:"email"`
}

func (e *Email) IsZero() bool { return e.Email == "" }

type MessengerID struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type PrivacyPolicy struct {
	PolicyID string `json:"policyId"`
}

type Company struct {
	Company string `json:"company"`
	Title   string `json:"title"`
}

func (c *Company) IsZero() bool {
	return c.Company == "" && c.Title == ""
}

type Phone struct {
	Phone       string `json:"phone"`
	Type        string `json:"type"`
	CountryCode string `json:"countryCode"`
}

func (p *Phone) IsZero() bool { return p.Phone == "" }

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Address struct {
	Street          string `json:"street"`
	ZipCode         string `json:"zipcode"`
	City            string `json:"city"`
	State           string `json:"state"`
	SubdivisionCode string `json:"subdivisionCode,omitempty"`
	Country         string `json:"country"`
}

func (a *Address) IsZero() bool {
	return a.Street == "" && a.ZipCode == "" && a.City == "" && a.State == "" && a.Country == ""
}

type BirthDate struct {
	BirthDate time.Time
}

func (b *BirthDate) Format(layout string) string {
	if b.BirthDate.IsZero() {
		return ""
	}

	return b.BirthDate.Format(layout)
}

func (d *BirthDate) UnmarshalJSON(data []byte) error {
	rawDate := &struct {
		Birthdate string `json:"birthdate"`
	}{}
	if err := json.Unmarshal(data, &rawDate); err != nil {
		return err
	} else if rawDate.Birthdate == "" {
		return nil
	}

	parsedDate, err := time.Parse("2006-01-02", strings.Split(rawDate.Birthdate, "T")[0])
	if err != nil {
		return err
	}

	d.BirthDate = parsedDate

	return nil
}

type BirthPlace struct {
	City    string `json:"city"`
	ZipCode string `json:"zipcode"`
	Country string `json:"country"`
}

func (p *BirthPlace) String() string {
	return joinStrings(", ", joinStrings(" ", p.ZipCode, p.City), p.Country)
}

type Documents struct {
	DrivingLicence *StructuredValue[Document]           `json:"drivingLicence,omitempty"`
	Passport       *StructuredValue[DocumentWithExpire] `json:"passport,omitempty"`
	IdentityCard   *StructuredValue[DocumentWithExpire] `json:"identityCard,omitempty"`
}

type Document struct {
	Id      string    `json:"id"`
	Country string    `json:"country,omitempty"`
	Issued  time.Time `json:"issued,omitempty"`
}

type DocumentWithExpire struct {
	Document
	Expires time.Time `json:"expires"`
}

type Loyalty struct {
	Program string  `json:"program"`
	Number  string  `json:"number"`
	Status  string  `json:"status"`
	Points  float32 `json:"points"`
}

type PMSID struct {
	ID  string `json:"id"`
	PMS string `json:"pms"`
}

type Source struct {
	Source   string `json:"source"`
	Metadata struct {
		Filename string `json:"filename"`
	} `json:"metadata"`
}

func joinStrings(separator string, values ...string) string {
	notZeroValues := []string{}

	for _, value := range values {
		if value != "" {
			notZeroValues = append(notZeroValues, value)
		}
	}

	return strings.Join(notZeroValues, separator)
}
