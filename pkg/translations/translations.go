// Package translations provides a type for translated values.
package translations

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/loungeup/go-loungeup/pkg/errors"
)

type (
	Translations     map[TranslationKey]TranslationValue
	TranslationKey   string
	TranslationValue string
)

func (t Translations) Equal(other Translations) bool {
	if len(t) != len(other) {
		return false
	}

	for key, value := range t {
		if otherValue, ok := other[key]; !ok || value != otherValue {
			return false
		}
	}

	return true
}

func (t Translations) Get(k TranslationKey, options ...GetOption) TranslationValue {
	config := &getConfig{
		defaultKey:        "en",
		getFirstByDefault: false,
	}
	for _, option := range options {
		option(config)
	}

	if result := t[k]; result != "" {
		return result
	}

	if result := t[config.defaultKey]; result != "" {
		return result
	}

	if config.getFirstByDefault {
		return t.first()
	}

	return ""
}

type GetOption func(*getConfig)

type getConfig struct {
	defaultKey        TranslationKey
	getFirstByDefault bool
}

func DefaultKey(key TranslationKey) GetOption { return func(c *getConfig) { c.defaultKey = key } }
func GetFirstByDefault() GetOption            { return func(c *getConfig) { c.getFirstByDefault = true } }

var _ (sql.Scanner) = (*Translations)(nil)

func (t Translations) IsZero() bool { return len(t) == 0 }

func (t *Translations) Scan(v any) error {
	vAsBytes, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("could not convert value to []byte")
	}

	return json.Unmarshal(vAsBytes, t)
}

func (t Translations) Validate() error {
	for k, v := range t {
		if k == "" {
			return &errors.Error{Code: errors.CodeInvalid, Message: "Translation key must not be empty"}
		}

		if v == "" {
			return &errors.Error{Code: errors.CodeInvalid, Message: "Translation value must not be empty"}
		}
	}

	return nil
}

var _ (driver.Valuer) = (Translations)(nil)

func (t Translations) Value() (driver.Value, error) { return json.Marshal(t) }

func (t Translations) first() TranslationValue {
	for _, v := range t {
		return v
	}

	return ""
}
