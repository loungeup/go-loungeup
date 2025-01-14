package models

import (
	"testing"
	"time"
)

func TestStructuredValueSlice_MostRecent(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		slice    StructuredValueSlice[string]
		expected StructuredValue[string]
	}{
		{
			name:     "empty slice",
			slice:    StructuredValueSlice[string]{},
			expected: StructuredValue[string]{},
		},
		{
			name: "single element",
			slice: StructuredValueSlice[string]{
				{Value: "one", UpdatedAt: now},
			},
			expected: StructuredValue[string]{Value: "one", UpdatedAt: now},
		},
		{
			name: "multiple elements",
			slice: StructuredValueSlice[string]{
				{Value: "one", UpdatedAt: now.Add(-2 * time.Hour)},
				{Value: "two", UpdatedAt: now.Add(-1 * time.Hour)},
				{Value: "three", UpdatedAt: now},
			},
			expected: StructuredValue[string]{Value: "three", UpdatedAt: now},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.slice.MostRecent()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
