package translations

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslationsGet(t *testing.T) {
	testTranslations := Translations{"en": "Hello", "fr": "Bonjour"}

	tests := map[string]struct {
		in   TranslationKey
		want TranslationValue
	}{
		"existing key": {in: "en", want: "Hello"},
		"unknown key":  {in: "de", want: ""},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, testTranslations.Get(tt.in))
		})
	}
}

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

func TestTranslationsValidate(t *testing.T) {
	tests := map[string]struct {
		in      Translations
		wantErr bool
	}{
		"valid translations": {in: Translations{"en": "Hello", "fr": "Bonjour"}, wantErr: false},
		"empty key":          {in: Translations{"": "Hello", "fr": "Bonjour"}, wantErr: true},
		"empty value":        {in: Translations{"en": "", "fr": "Bonjour"}, wantErr: true},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, tt.in.Validate() != nil)
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
