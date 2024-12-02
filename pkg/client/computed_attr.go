package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/translations"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

type computedAttrsClient struct{ baseClient *Client }

func (c *computedAttrsClient) ReadOne(selector *ComputedAttrSelector) (*ComputedAttr, error) {
	return transport.GetRESModel[*ComputedAttr](c.baseClient.resClient, selector.rid(), resprot.Request{})
}

type ComputedAttr struct {
	ID         uuid.UUID                                `json:"id"`
	EntityID   uuid.UUID                                `json:"entityId"`
	Type       ComputedAttrType                         `json:"type"`
	ValueType  ComputedAttrValueType                    `json:"valueType"`
	Scope      ComputedAttrScope                        `json:"scope"`
	Name       res.DataValue[translations.Translations] `json:"name,omitempty"`
	Conditions *res.DataValue[json.RawMessage]          `json:"conditions,omitempty"`
	TaskRID    string                                   `json:"taskRid,omitempty"`
	CreatedAt  time.Time                                `json:"createdAt"`
	DisabledAt time.Time                                `json:"disabledAt,omitempty"`
	EnabledAt  time.Time                                `json:"enabledAt,omitempty"`
	UpdatedAt  time.Time                                `json:"updatedAt"`
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

type ComputedAttrSelector struct {
	AttrID   uuid.UUID
	EntityID uuid.UUID
}

func (s *ComputedAttrSelector) rid() string {
	if s.EntityID == uuid.Nil {
		return "guestprofile.computed-attributes." + s.AttrID.String()
	}

	return "guestprofile.entities." + s.EntityID.String() + ".computed-attributes." + s.AttrID.String()
}
