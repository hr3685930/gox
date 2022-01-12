package gin

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
)

var trans ut.Translator

// LoadValidatorLocal 初始化语言包
func LoadValidatorLocal(local string) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() //chinese
		enT := en.New() //english
		uni := ut.New(enT, zhT, enT)
		var o bool
		trans, o = uni.GetTranslator(local)
		if !o {
			return errors.New("uni.GetTranslator failed")
		}
		// register translate
		var err error
		switch local {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		}

		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func Translate(err error) (errMsgs []string) {
	errs := err.(validator.ValidationErrors)
	for _, err := range errs {
		errMsg := err.Translate(trans)
		errMsgs = append(errMsgs, errMsg)
	}
	return errMsgs
}
