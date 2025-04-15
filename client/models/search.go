package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pointer"
)

type SearchCriteriaOperator string

const (
	SearchOperatorDateAfterThan            SearchCriteriaOperator = "&>"
	SearchOperatorDateAfterThanEqual       SearchCriteriaOperator = "&>="
	SearchOperatorDateAfterThanKeywordDay  SearchCriteriaOperator = "&>%"
	SearchOperatorAnonymous                SearchCriteriaOperator = "@"
	SearchOperatorDateBeforeThan           SearchCriteriaOperator = "&<"
	SearchOperatorDateBeforeThanEqual      SearchCriteriaOperator = "&<="
	SearchOperatorDateBeforeThanKeywordDay SearchCriteriaOperator = "&<%"
	SearchOperatorCampaignAnswered         SearchCriteriaOperator = "A"
	SearchOperatorCampaignNotAnswered      SearchCriteriaOperator = "!A"
	SearchOperatorCampaignNotOpened        SearchCriteriaOperator = "!O"
	SearchOperatorCampaignNotReceived      SearchCriteriaOperator = "!R"
	SearchOperatorCampaignNotScheduled     SearchCriteriaOperator = "!S"
	SearchOperatorCampaignOpened           SearchCriteriaOperator = "O"
	SearchOperatorCampaignReceived         SearchCriteriaOperator = "R"
	SearchOperatorCampaignScheduled        SearchCriteriaOperator = "S"
	SearchOperatorContains                 SearchCriteriaOperator = "%"
	SearchOperatorDateAwayLess             SearchCriteriaOperator = "-"
	SearchOperatorDateAwayMore             SearchCriteriaOperator = "+"
	SearchOperatorDow                      SearchCriteriaOperator = "DOW"
	SearchOperatorDateEqual                SearchCriteriaOperator = "&="
	SearchOperatorEquals                   SearchCriteriaOperator = "="
	SearchOperatorFilled                   SearchCriteriaOperator = "F"
	SearchOperatorHadSession               SearchCriteriaOperator = "1"
	SearchOperatorHasCurrentSession        SearchCriteriaOperator = "C"
	SearchOperatorInferior                 SearchCriteriaOperator = "<"
	SearchOperatorNo                       SearchCriteriaOperator = "N"
	SearchOperatorNotAnonymous             SearchCriteriaOperator = "!@"
	SearchOperatorNotContains              SearchCriteriaOperator = "!%"
	SearchOperatorNotEquals                SearchCriteriaOperator = "!="
	SearchOperatorNotFilled                SearchCriteriaOperator = "!F"
	SearchOperatorNotHadSession            SearchCriteriaOperator = "!1"
	SearchOperatorNotStartsWith            SearchCriteriaOperator = "!S"
	SearchOperatorRange                    SearchCriteriaOperator = "RANGE"
	SearchOperatorStartsWith               SearchCriteriaOperator = "S"
	SearchOperatorSuperior                 SearchCriteriaOperator = ">"
	SearchOperatorYes                      SearchCriteriaOperator = "Y"
)

func (s SearchCriteriaOperator) String() string { return string(s) }

type SearchKeywordDate string

const (
	SearchKeywordDateUnknown   SearchKeywordDate = ""
	SearchKeywordDateToday     SearchKeywordDate = "today"
	SearchKeywordDateTomorrow  SearchKeywordDate = "tomorrow"
	SearchKeywordDateYesterday SearchKeywordDate = "yesterday"
)

const hoursInDay = 24

func ParseSearchKeywordDate(keyword any) SearchKeywordDate {
	switch keyword {
	case "today":
		return SearchKeywordDateToday
	case "tomorrow":
		return SearchKeywordDateTomorrow
	case "yesterday":
		return SearchKeywordDateYesterday
	}

	return SearchKeywordDateUnknown
}

func (d SearchKeywordDate) Duration() *time.Duration {
	switch d {
	case SearchKeywordDateToday:
		return pointer.From(0 * time.Hour)
	case SearchKeywordDateTomorrow:
		return pointer.From(hoursInDay * time.Hour)
	case SearchKeywordDateYesterday:
		return pointer.From(-hoursInDay * time.Hour)
	}

	return nil
}

type SearchCriteriaField string

const (
	SearchCriteriaFieldApp                  SearchCriteriaField = "app"
	SearchCriteriaFieldArrival              SearchCriteriaField = "arrival"
	SearchCriteriaFieldArrivalDOW           SearchCriteriaField = "arrivaldow"
	SearchCriteriaFieldBalance              SearchCriteriaField = "balance"
	SearchCriteriaFieldBirthdate            SearchCriteriaField = "birthdate"
	SearchCriteriaFieldBirthday             SearchCriteriaField = "birthday"
	SearchCriteriaFieldBookingDate          SearchCriteriaField = "bookingdate"
	SearchCriteriaFieldBookingID            SearchCriteriaField = "bookingId"
	SearchCriteriaFieldBookingPurpose       SearchCriteriaField = "bookingPurpose"
	SearchCriteriaFieldBookingStatus        SearchCriteriaField = "bookingstatus"
	SearchCriteriaFieldBookingStatusDiff    SearchCriteriaField = "bookingstatusdiff"
	SearchCriteriaFieldBookingTags          SearchCriteriaField = "bookingTags"
	SearchCriteriaFieldBookingWindow        SearchCriteriaField = "bookingarrival"
	SearchCriteriaFieldCampaignApp          SearchCriteriaField = "msgcampaign"
	SearchCriteriaFieldCampaignEmail        SearchCriteriaField = "campaign"
	SearchCriteriaFieldCampaignNewsletter   SearchCriteriaField = "newsletter"
	SearchCriteriaFieldCampaignScheduled    SearchCriteriaField = "scheduledCampaign"
	SearchCriteriaFieldCampaignSearch       SearchCriteriaField = "campaignSearch"
	SearchCriteriaFieldCampaignSMS          SearchCriteriaField = "smscampaign"
	SearchCriteriaFieldCampaignWhatsApp     SearchCriteriaField = "whatsappcampaign"
	SearchCriteriaFieldChannel              SearchCriteriaField = "channel"
	SearchCriteriaFieldCity                 SearchCriteriaField = "city"
	SearchCriteriaFieldCompany              SearchCriteriaField = "company"
	SearchCriteriaFieldCountry              SearchCriteriaField = "country"
	SearchCriteriaFieldCustomFieldsBooking  SearchCriteriaField = "visitCustomField"
	SearchCriteriaFieldCustomFieldsGuest    SearchCriteriaField = "userCustomField"
	SearchCriteriaFieldDeparture            SearchCriteriaField = "departure"
	SearchCriteriaFieldDepartureDOW         SearchCriteriaField = "departuredow"
	SearchCriteriaFieldEmail                SearchCriteriaField = "email"
	SearchCriteriaFieldEmailCollected       SearchCriteriaField = "emailcollected"
	SearchCriteriaFieldEmailDomain          SearchCriteriaField = "emaildomain"
	SearchCriteriaFieldEntityGuestUUIDs     SearchCriteriaField = "entityGuestUuids"
	SearchCriteriaFieldEntityID             SearchCriteriaField = "entityId"
	SearchCriteriaFieldEntityObjectUUID     SearchCriteriaField = "entityObjectUuid"
	SearchCriteriaFieldFare                 SearchCriteriaField = "fare"
	SearchCriteriaFieldFareAvg              SearchCriteriaField = "avgfare"
	SearchCriteriaFieldFareCode             SearchCriteriaField = "farecode"
	SearchCriteriaFieldFareRatio            SearchCriteriaField = "fareratio"
	SearchCriteriaFieldFareSum              SearchCriteriaField = "sumfare"
	SearchCriteriaFieldFidelity             SearchCriteriaField = "fidelity"
	SearchCriteriaFieldFidelityStatus       SearchCriteriaField = "bwrewardsstatus"
	SearchCriteriaFieldGuestTags            SearchCriteriaField = "guestTags"
	SearchCriteriaFieldGuestUUID            SearchCriteriaField = "guestId"
	SearchCriteriaFieldHasFacebook          SearchCriteriaField = "hasfacebook"
	SearchCriteriaFieldHasLinkedIn          SearchCriteriaField = "haslinkedin"
	SearchCriteriaFieldHasPersonnalEmail    SearchCriteriaField = "hasPersonnalEmail"
	SearchCriteriaFieldHasTwitter           SearchCriteriaField = "hastwitter"
	SearchCriteriaFieldIDMasterResa         SearchCriteriaField = "idmasterresa"
	SearchCriteriaFieldIDResa               SearchCriteriaField = "idresa"
	SearchCriteriaFieldInStay               SearchCriteriaField = "instay"
	SearchCriteriaFieldInStayDate           SearchCriteriaField = "instaydate"
	SearchCriteriaFieldInStayDates          SearchCriteriaField = "instayDates"
	SearchCriteriaFieldLangs                SearchCriteriaField = "langs"
	SearchCriteriaFieldLastConnexion        SearchCriteriaField = "lastconnexion"
	SearchCriteriaFieldLastName             SearchCriteriaField = "lastname"
	SearchCriteriaFieldLinkedInFollowers    SearchCriteriaField = "linkedinfollowers"
	SearchCriteriaFieldMetadata             SearchCriteriaField = "metadata"
	SearchCriteriaFieldNationality          SearchCriteriaField = "nationality"
	SearchCriteriaFieldNbNextStays          SearchCriteriaField = "nbnextstays"
	SearchCriteriaFieldNbPreviousStays      SearchCriteriaField = "nbpreviousstays"
	SearchCriteriaFieldNbStays              SearchCriteriaField = "nbstays"
	SearchCriteriaFieldNextStay             SearchCriteriaField = "nextstay"
	SearchCriteriaFieldOptinAuto            SearchCriteriaField = "optoutprestay"
	SearchCriteriaFieldOptinCustomerAccount SearchCriteriaField = "optoutcustomeraccount"
	SearchCriteriaFieldOptinLoyalty         SearchCriteriaField = "optoutloyalty"
	SearchCriteriaFieldOptinMarketing       SearchCriteriaField = "optoutmarketing"
	SearchCriteriaFieldOptinSendInBlue      SearchCriteriaField = "optoutsendinblue"
	SearchCriteriaFieldPartner              SearchCriteriaField = "partner"
	SearchCriteriaFieldPaxAdults            SearchCriteriaField = "paxadults"
	SearchCriteriaFieldPaxBabies            SearchCriteriaField = "paxbabies"
	SearchCriteriaFieldPaxChildren          SearchCriteriaField = "paxchildren"
	SearchCriteriaFieldPhone                SearchCriteriaField = "phone"
	SearchCriteriaFieldPreferredEmail       SearchCriteriaField = "preferredEmail"
	SearchCriteriaFieldPreviousStay         SearchCriteriaField = "prevstay"
	SearchCriteriaFieldPushNotification     SearchCriteriaField = "pushNotification"
	SearchCriteriaFieldRoomNumber           SearchCriteriaField = "roomnumber"
	SearchCriteriaFieldRoomType             SearchCriteriaField = "roomtype"
	SearchCriteriaFieldSearch               SearchCriteriaField = "search"
	SearchCriteriaFieldSegment              SearchCriteriaField = "segment"
	SearchCriteriaFieldSourceImport         SearchCriteriaField = "sourceImport"
	SearchCriteriaFieldSourceType           SearchCriteriaField = "sourceType"
	SearchCriteriaFieldStayLength           SearchCriteriaField = "staylength"
	SearchCriteriaFieldTags                 SearchCriteriaField = "tags"
	SearchCriteriaFieldTouristTax           SearchCriteriaField = "touristtax"
	SearchCriteriaFieldTwitterFollowers     SearchCriteriaField = "twitterfollowers"
	SearchCriteriaFieldUpdatedAt            SearchCriteriaField = "updatedAt"
	SearchCriteriaFieldWeekend              SearchCriteriaField = "weekend"
	SearchCriteriaFieldZipCode              SearchCriteriaField = "zipcode"
)

type SearchConditions struct {
	Logic    string             `json:"logic"`
	Criteria []*SearchCriterion `json:"criteria"`
}

type SearchCriterion struct {
	Logic    string            `json:"logic"`
	Criteria []*SearchCriteria `json:"criteria"`
}

type SearchCriteria struct {
	Field    SearchCriteriaField    `json:"field"`
	Operator SearchCriteriaOperator `json:"operator"`
	Value    any                    `json:"value"`
}

type StructuredValueCriterion[T any] struct {
	Preferred bool       `json:"preferred,omitempty"`
	Value     T          `json:"value,omitempty"`
	UpdatedAt *time.Time `json:"updated,omitempty"`
	From      string     `json:"from,omitempty"`
}

type PMSIDCriterion struct {
	ID  string `json:"id,omitempty"`
	PMS string `json:"pms,omitempty"`
}
type MessengerIDCriterion struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type PhoneCriterion struct {
	Phone       string `json:"phone,omitempty"`
	Type        string `json:"type,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

type SearchByContactMatch string

const (
	SearchByContactMatchOne SearchByContactMatch = "one"
	SearchByContactMatchAll SearchByContactMatch = "all"
)

type SearchByContactSelector struct {
	EntityID uuid.UUID `json:"-"`

	Match SearchByContactMatch `json:"match,omitempty"`

	Emails      []*StructuredValueCriterion[Email]              `json:"emails,omitempty"`
	Phones      []*StructuredValueCriterion[PhoneCriterion]     `json:"phones,omitempty"`
	MessengerID *StructuredValueCriterion[MessengerIDCriterion] `json:"messengerId,omitempty"`
	PMSID       *StructuredValueCriterion[PMSIDCriterion]       `json:"pmsId,omitempty"`
	Credentials *StructuredValueCriterion[Credentials]          `json:"credentials,omitempty"`

	LastKey string `json:"lastKey,omitempty"`
	Size    int    `json:"size,omitempty"`
}
