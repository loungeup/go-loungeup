package esutil

import (
	"encoding/json"
	"strings"

	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	essortorder "github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/jsonutil"
	"github.com/loungeup/go-loungeup/pointer"
)

//nolint:forcetypeassert
var computedAttrAggConfigs = map[string]ComputedAttrAggConfig{
	"accountIds": {
		Agg: estypes.Aggregations{
			Terms: &estypes.TermsAggregation{
				Field: pointer.From("booking.entityId"),
				Size:  pointer.From(20), //nolint:mnd
			},
		},
		AggType: ComputedAttrAggTypeNumber,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			buckets := aggregate.(*estypes.StringTermsAggregate).Buckets.([]estypes.StringTermsBucket)
			if len(buckets) == 0 {
				return nil
			}

			entityIds := uuid.UUIDs{}
			for _, bucket := range buckets {
				key, ok := bucket.Key.(string)
				if !ok {
					continue
				}

				entityId, _ := uuid.Parse(key)

				entityIds = append(entityIds, entityId)
			}

			if len(entityIds) == 0 {
				return nil
			}

			return entityIds
		},
	},
	"averageRevenue": {
		Agg: estypes.Aggregations{
			Avg: &estypes.AverageAggregation{
				Field: pointer.From("booking.fare"),
			},
		},
		AggType: ComputedAttrAggTypeNumber,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return mapESFloat64(aggregate.(*estypes.AvgAggregate).Value)
		},
	},
	"mostRelevantBookingId": {
		Agg: estypes.Aggregations{
			ScriptedMetric: &estypes.ScriptedMetricAggregation{
				InitScript: &estypes.Script{
					Id: pointer.From("compute-guest-current-booking-init"),
				},
				MapScript: &estypes.Script{
					Id: pointer.From("compute-guest-current-booking-map"),
				},
				CombineScript: &estypes.Script{
					Id: pointer.From("compute-guest-current-booking-combine"),
				},
				ReduceScript: &estypes.Script{
					Id: pointer.From("compute-guest-current-booking-reduce"),
				},
				Params: map[string]json.RawMessage{
					"entityType": nil, // This is a placeholder for the actual value.
				},
			},
		},
		AggType: ComputedAttrAggTypeText,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			return jsonutil.Document(aggregate.(*estypes.ScriptedMetricAggregate).Value).Child("id")
		},
	},
	"nextAccountId": {
		Agg: estypes.Aggregations{
			TopHits: &estypes.TopHitsAggregation{
				Size: pointer.From(1),
				Sort: []estypes.SortCombinations{
					estypes.SortOptions{
						SortOptions: map[string]estypes.FieldSort{
							"booking.arrival": {
								Order: &essortorder.Asc,
							},
						},
					},
				},
				Source_: estypes.SourceFilter{
					Includes: []string{"booking.entityId"},
				},
			},
		},
		AggType: ComputedAttrAggTypeText,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			hits := aggregate.(*estypes.TopHitsAggregate).Hits.Hits
			if len(hits) == 0 {
				return nil
			}

			return jsonutil.Document(hits[0].Source_).UUID("booking.entityId")
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
			return mapESFloat64(aggregate.(*estypes.MinAggregate).Value)
		},
	},
	"previousAccountId": {
		Agg: estypes.Aggregations{
			TopHits: &estypes.TopHitsAggregation{
				Size: pointer.From(1),
				Sort: []estypes.SortCombinations{
					estypes.SortOptions{
						SortOptions: map[string]estypes.FieldSort{
							"booking.departure": {
								Order: &essortorder.Desc,
							},
						},
					},
				},
				Source_: estypes.SourceFilter{
					Includes: []string{"booking.entityId"},
				},
			},
		},
		AggType: ComputedAttrAggTypeText,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			hits := aggregate.(*estypes.TopHitsAggregate).Hits.Hits
			if len(hits) == 0 {
				return nil
			}

			return jsonutil.Document(hits[0].Source_).UUID("booking.entityId")
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
			return mapESFloat64(aggregate.(*estypes.MaxAggregate).Value)
		},
	},
	"totalAccounts": {
		Agg: estypes.Aggregations{
			Terms: &estypes.TermsAggregation{
				Field: pointer.From("booking.entityId"),
			},
		},
		AggType: ComputedAttrAggTypeNumber,
		MapValueFunc: func(aggregate estypes.Aggregate) any {
			buckets := aggregate.(*estypes.StringTermsAggregate).Buckets.([]estypes.StringTermsBucket)

			return len(buckets)
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
			return mapESFloat64(aggregate.(*estypes.ValueCountAggregate).Value)
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
			return mapESFloat64(aggregate.(*estypes.SumAggregate).Value)
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
			return mapESFloat64(aggregate.(*estypes.SumAggregate).Value)
		},
	},
}

func GetComputedAttrAggConfig[T ~string](v T, entityType models.EntityType) ComputedAttrAggConfig {
	if result, ok := computedAttrAggConfigs[string(v)]; ok {
		if string(v) == "mostRelevantBookingId" {
			result.Agg.ScriptedMetric.Params["entityType"] = json.RawMessage(`"` + entityType.String() + `"`)
		}

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

func mapESFloat64(value *estypes.Float64) float64 {
	if value == nil {
		return 0
	}

	return float64(*value)
}
