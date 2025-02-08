package locales

import (
	"github.com/go-playground/locales/en"
	indonesia "github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
)

var (
	trans ut.Translator
)

func InitIdTrans() ut.Translator {
	if trans != nil {
		return trans
	}
	idn := indonesia.New()
	uni := ut.New(idn, idn)
	var found bool
	trans, found = uni.GetTranslator("id")
	if !found {
		panic("Translator id not found")
	}
	return trans
}

func InitDefaultTrans() ut.Translator {
	if trans != nil {
		return trans
	}
	enn := en.New()
	uni := ut.New(enn, enn)
	var found bool
	trans, found = uni.GetTranslator("en")
	if !found {
		panic("Translator en not found")
	}
	return trans
}

func GetTrans() ut.Translator {
	if trans != nil {
		return trans
	}

	return InitDefaultTrans() // default
}
