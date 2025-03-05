package translations

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslationsEqual(t *testing.T) {
	tests := map[string]struct {
		a, b Translations
		want bool
	}{
		"empty": {
			a:    Translations{},
			b:    Translations{},
			want: true,
		},
		"equal": {
			a:    Translations{"en": "Hello", "fr": "Bonjour"},
			b:    Translations{"fr": "Bonjour", "en": "Hello"},
			want: true,
		},
		"extra key in b": {
			a:    Translations{"en": "Hello"},
			b:    Translations{"en": "Hello", "fr": "Bonjour"},
			want: false,
		},
		"extra key in a": {
			a:    Translations{"en": "Hello", "fr": "Bonjour"},
			b:    Translations{"en": "Hello"},
			want: false,
		},
		"not equal": {
			a:    Translations{"en": "Hello", "fr": "Bonjour"},
			b:    Translations{"en": "Hello", "fr": "Salut"},
			want: false,
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.a.Equal(tt.b))
		})
	}
}

func TestTranslationsGet(t *testing.T) {
	t.Run("existing key", func(t *testing.T) {
		translations := Translations{"en": "Hello", "fr": "Bonjour"}
		assert.Equal(t, TranslationValue("Hello"), translations.Get("en"))
	})

	t.Run("unknown key", func(t *testing.T) {
		translations := Translations{"fr": "Bonjour"}
		assert.Equal(t, TranslationValue(""), translations.Get("de"))
	})

	t.Run("default to en", func(t *testing.T) {
		translations := Translations{"en": "Hello"}
		assert.Equal(t, TranslationValue("Hello"), translations.Get("foo"))
	})

	t.Run("default to fr", func(t *testing.T) {
		translations := Translations{"fr": "Bonjour"}
		assert.Equal(t, TranslationValue("Bonjour"), translations.Get("foo", DefaultKey("fr")))
	})

	t.Run("get first by default", func(t *testing.T) {
		translations := Translations{"de": "Hallo", "fr": "Bonjour"}
		assert.Equal(t, TranslationValue("Hallo"), translations.Get("foo", GetFirstByDefault()))
	})
}

func TestTranslationsIsZero(t *testing.T) {
	tests := map[string]struct {
		in   Translations
		want bool
	}{
		"empty":     {in: Translations{}, want: true},
		"not empty": {in: Translations{"en": "Hello"}, want: false},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.in.IsZero())
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
