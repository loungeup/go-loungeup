package esutil

import (
	"strings"

	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/loungeup/go-loungeup/pkg/pointer"
)

//nolint:forcetypeassert
var computedAttrAggConfigs = map[string]ComputedAttrAggConfig{
	"accountIds": {}, // TODO(remyduthu): Implement.
	"averageRevenue": {
		Agg: estypes.Aggregations{
			Avg: &estypes.AverageAggregation{
				Field: pointer.From("booking.fare"),
			},
		},
		AggType: ComputedAttrAggTypeNumber,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return aggregate.(*estypes.AvgAggregate).Value
		},
	},
	"nextAccountId": {
		Agg:     estypes.Aggregations{}, // TODO(remyduthu): Implement.
		AggType: ComputedAttrAggTypeText,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return nil // TODO: Implement.
		},
	},
	"nextBookingArrival": {
		Agg: estypes.Aggregations{
			Min: &estypes.MinAggregation{
				Field: pointer.From("booking.arrival"),
			},
		},
		AggType: ComputedAttrAggTypeDate,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return aggregate.(*estypes.MinAggregate).Value
		},
	},
	"previousAccountId": {
		Agg:     estypes.Aggregations{}, // TODO(remyduthu): Implement.
		AggType: ComputedAttrAggTypeText,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return nil // TODO: Implement.
		},
	},
	"previousBookingDeparture": {
		Agg: estypes.Aggregations{
			Max: &estypes.MaxAggregation{
				Field: pointer.From("booking.departure"),
			},
		},
		AggType: ComputedAttrAggTypeDate,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return aggregate.(*estypes.MaxAggregate).Value
		},
	},
	"totalAccounts": {
		Agg: estypes.Aggregations{
			ValueCount: &estypes.ValueCountAggregation{
				Field: pointer.From("booking.entityId"),
			},
		},
		AggType: ComputedAttrAggTypeNumber,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return aggregate.(*estypes.ValueCountAggregate).Value
		},
	},
	"totalBookings": {
		Agg: estypes.Aggregations{
			ValueCount: &estypes.ValueCountAggregation{
				Field: pointer.From("booking.id"),
			},
		},
		AggType: ComputedAttrAggTypeNumber,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return aggregate.(*estypes.ValueCountAggregate).Value
		},
	},
	"totalDistinctBookings": {
		Agg: estypes.Aggregations{
			Cardinality: &estypes.CardinalityAggregation{
				Script: &estypes.Script{
					Source: pointer.From("doc['booking.arrival'].value + '_' + doc['booking.departure'].value"),
				},
			},
		},
		AggType: ComputedAttrAggTypeNumber,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return aggregate.(*estypes.CardinalityAggregate).Value
		},
	},
	"totalNights": {
		Agg: estypes.Aggregations{
			Sum: &estypes.SumAggregation{
				Field: pointer.From("booking.stayLength"),
			},
		},
		AggType: ComputedAttrAggTypeNumber,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return aggregate.(*estypes.SumAggregate).Value
		},
	},
	"totalRevenue": {
		Agg: estypes.Aggregations{
			Sum: &estypes.SumAggregation{
				Field: pointer.From("booking.fare"),
			},
		},
		AggType: ComputedAttrAggTypeNumber,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return aggregate.(*estypes.SumAggregate).Value
		},
	},
}

func GetComputedAttrAggConfig[T ~string](v T) ComputedAttrAggConfig {
	if result, ok := computedAttrAggConfigs[string(v)]; ok {
		return result
	}

	return ComputedAttrAggConfig{}
}

type ComputedAttrAggConfig struct {
	Agg          estypes.Aggregations
	AggType      ComputedAttrAggType
	MapValueFunc func(aggregate estypes.Aggregate) any
}

type ComputedAttrAggType string

const (
	ComputedAttrAggTypeUnknown ComputedAttrAggType = ""
	ComputedAttrAggTypeBoolean ComputedAttrAggType = "boolean"
	ComputedAttrAggTypeDate    ComputedAttrAggType = "date"
	ComputedAttrAggTypeNumber  ComputedAttrAggType = "number"
	ComputedAttrAggTypeText    ComputedAttrAggType = "text"
)

func (t ComputedAttrAggType) String() string { return string(t) }

func MapEntityType[T ~string](v T) string {
	switch {
	case strings.EqualFold(string(v), "account"):
		return "account"
	case strings.EqualFold(string(v), "chain"):
		return "chain"
	case strings.EqualFold(string(v), "group"):
		return "group"
	default:
		return "" // We are never using other types in ES (e.g. resellers).
	}
}
