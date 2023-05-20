package validate

import (
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator ut.Translator

func init() {

	// Instantiate a validator.
	validate = validator.New()

	// Create a translator for english so the error messages are
	// more human-readable than technical.
	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, translator)
}

// Check validates the provided model against it's declared tags.
func Check(val any) error {
	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		if len(verrors) < 1 {
			return nil
		}

		// Return only the first error.
		return errors.New(verrors[0].Translate(translator))
	}

	return nil
}

// GenerateID generate a unique id for entities.
func GenerateID() string {
	return uuid.NewString()
}

// CheckID validates that the format of an id is valid.
func CheckID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("ID is not in its proper form")
	}
	return nil
}
