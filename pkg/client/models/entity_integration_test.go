package models

import (
	"testing"

	"github.com/loungeup/go-loungeup/pkg/matcher"
	"github.com/stretchr/testify/assert"
)

func TestGetEntityIntegrationValue(t *testing.T) {
	tests := map[string]struct {
		getValue  func() (any, error)
		want      any
		wantError bool
	}{
		"string value": {
			getValue: func() (any, error) {
				return GetEntityIntegrationValue[string](EntityIntegrationValues{"foo": "bar"}, "foo")
			},
			want: "bar",
		},
		"matcher value": {
			getValue: func() (any, error) {
				return GetEntityIntegrationValue[matcher.Matcher](EntityIntegrationValues{
					"foo": map[string]string{
						matcher.WellKnownMatcherKeys.UpgradeProductID.String(): "bar",
					},
				}, "foo")
			},
			want: matcher.Matcher{
				matcher.WellKnownMatcherKeys.UpgradeProductID: "bar",
			},
		},
		"invalid key": {
			getValue: func() (any, error) {
				return GetEntityIntegrationValue[string](EntityIntegrationValues{"foo": "bar"}, "baz")
			},
			wantError: true,
		},
		"invalid value type": {
			getValue: func() (any, error) {
				return GetEntityIntegrationValue[string](EntityIntegrationValues{"foo": 42}, "foo")
			},
			wantError: true,
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got, err := tt.getValue()
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
