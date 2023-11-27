package loungeup

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type (
	// Translations represents a map of translations.
	Translations map[TranslationLanguageCode]TranslatedContent

	// TranslationLanguageCode represents the language code of the translation. It is an ISO 639-1 code.
	TranslationLanguageCode string

	// TranslatedContent represents the content of a translation.
	TranslatedContent string
)

var _ (sql.Scanner) = (Translations)(nil)

func (t Translations) Scan(value any) error {
	valueAsBytes, ok := value.([]byte)
	if !ok {
		return errors.New("could not convert value to []byte")
	}

	return json.Unmarshal(valueAsBytes, &t)
}

var _ (driver.Valuer) = (Translations)(nil)

func (t Translations) Value() (driver.Value, error) { return json.Marshal(t) }
