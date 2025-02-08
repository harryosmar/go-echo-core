package presentation

import (
	"github.com/go-playground/validator/v10"
	coreError "github.com/harryosmar/go-echo-core/error"
	"github.com/harryosmar/go-echo-core/locales"
	"strings"
)

func ResponseErr(err error) error {
	var errCode coreError.CodeErrEntity
	errCode = coreError.ErrGeneral.WithError(err)

	if codeErr, ok := err.(coreError.CodeErrEntity); ok {
		errCode = codeErr
	} else if typeOfErrEntity, isTypeOfErrEntity := err.(coreError.TypeOfErrEntity); isTypeOfErrEntity {
		errCode = typeOfErrEntity.GetCodeErrEntity()
	}

	return NewResponseEntity().
		WithStatusCode(errCode.Status).
		WithContentStatus(false).
		WithMessage(errCode.Message, errCode.Args...).
		WithErrorCode(errCode.Code).
		WithError(errCode)
}

func ResponseErrValidation(err error) error {
	if errors, ok := err.(validator.ValidationErrors); ok {
		data := map[string]interface{}{}
		for _, fieldError := range errors {
			fieldName := strings.ToLower(fieldError.Field())
			errStr := fieldError.Translate(locales.GetTrans())
			data[fieldName] = errStr
		}

		return NewResponseEntity().
			WithStatusCode(coreError.ErrValidation.Status()).
			WithContentStatus(false).
			WithMessage(coreError.ErrValidation.String()).
			WithErrorCode(coreError.ErrValidation.Code()).
			WithData(data)
	}

	return ResponseErr(err)
}
