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

func (t Translations) Get(k TranslationKey) TranslationValue { return t[k] }

var _ (sql.Scanner) = (Translations)(nil)

func (t Translations) IsZero() bool { return len(t) == 0 }

func (t Translations) Scan(v any) error {
	vAsBytes, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("could not convert value to []byte")
	}

	return json.Unmarshal(vAsBytes, &t)
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
