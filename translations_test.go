package loungeup

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranslationsScan(t *testing.T) {
	tests := map[string]struct {
		in      any
		want    Translations
		wantErr bool
	}{
		"valid JSON": {
			in:      []byte(`{"en": "Hello", "fr": "Bonjour"}`),
			want:    Translations{"en": "Hello", "fr": "Bonjour"},
			wantErr: false,
		},
		"invalid JSON": {
			in:      []byte(`{"en": "Hello", "fr": }`),
			wantErr: true,
		},
		"invalid type": {
			in:      123,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Translations{}
			if err := got.Scan(tt.in); tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestTranslationsValue(t *testing.T) {
	tests := map[string]struct {
		in   Translations
		want driver.Value
	}{
		"valid translations": {
			in:   Translations{"en": "Hello", "fr": "Bonjour"},
			want: []byte(`{"en":"Hello","fr":"Bonjour"}`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tt.in.Value()
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
