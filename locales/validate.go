package locales

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/id"
)

func InitValidate(trans ut.Translator, locale string) (*validator.Validate, error) {
	var err error
	validate := validator.New(validator.WithRequiredStructEnabled())
	switch locale {
	case "id":
		InitIdTrans()
		err = id.RegisterDefaultTranslations(validate, trans)
	default:
		InitDefaultTrans()
		err = en.RegisterDefaultTranslations(validate, trans)
	}
	return validate, err
}
