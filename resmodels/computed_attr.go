package resmodels

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/timeutil"
	"github.com/loungeup/go-loungeup/translations"
)

type ComputedAttr struct {
	ID         uuid.UUID                                `json:"id"`
	EntityID   uuid.UUID                                `json:"entityId"`
	Type       ComputedAttrType                         `json:"type"`
	ValueType  ComputedAttrValueType                    `json:"valueType"`
	Scope      ComputedAttrScope                        `json:"scope"`
	Name       res.DataValue[translations.Translations] `json:"name"`
	Conditions *res.DataValue[json.RawMessage]          `json:"conditions"`
	TaskRID    string                                   `json:"taskRid"`
	ComputedAt timeutil.RFC3339Time                     `json:"computedAt"`
	CreatedAt  timeutil.RFC3339Time                     `json:"createdAt"`
	DisabledAt timeutil.RFC3339Time                     `json:"disabledAt"`
	EnabledAt  timeutil.RFC3339Time                     `json:"enabledAt"`
	UpdatedAt  timeutil.RFC3339Time                     `json:"updatedAt"`
}

type ComputedAttrType string

const (
	ComputedAttrTypeUnknown  ComputedAttrType = ""
	ComputedAttrTypeCustom   ComputedAttrType = "custom"
	ComputedAttrTypeStandard ComputedAttrType = "standard"
)

func (t ComputedAttrType) String() string { return string(t) }

type ComputedAttrValueType string

const (
	ComputedAttrValueTypeUnknown                  ComputedAttrValueType = ""
	ComputedAttrValueTypeAccountIDs               ComputedAttrValueType = "accountIds"
	ComputedAttrValueTypeAverageRevenue           ComputedAttrValueType = "averageRevenue"
	ComputedAttrValueTypeNextAccountID            ComputedAttrValueType = "nextAccountId"
	ComputedAttrValueTypeNextBookingArrival       ComputedAttrValueType = "nextBookingArrival"
	ComputedAttrValueTypePreviousAccountID        ComputedAttrValueType = "previousAccountId"
	ComputedAttrValueTypePreviousBookingDeparture ComputedAttrValueType = "previousBookingDeparture"
	ComputedAttrValueTypeTotalAccounts            ComputedAttrValueType = "totalAccounts"
	ComputedAttrValueTypeTotalBookings            ComputedAttrValueType = "totalBookings"
	ComputedAttrValueTypeTotalDistinctBookings    ComputedAttrValueType = "totalDistinctBookings"
	ComputedAttrValueTypeTotalNights              ComputedAttrValueType = "totalNights"
	ComputedAttrValueTypeTotalRevenue             ComputedAttrValueType = "totalRevenue"
)

func (t ComputedAttrValueType) String() string { return string(t) }

type ComputedAttrScope string

const (
	ComputedAttrScopeUnknown  ComputedAttrScope = ""
	ComputedAttrScopeAccounts ComputedAttrScope = "accounts"
	ComputedAttrScopeEntity   ComputedAttrScope = "entity"
)

func (s ComputedAttrScope) String() string { return string(s) }
