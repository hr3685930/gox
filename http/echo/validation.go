package echo

import (
    "github.com/go-playground/locales/en"
    "github.com/go-playground/locales/zh"
    "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    zh_trans "github.com/go-playground/validator/v10/translations/zh"
    "sync"
)

var uni *ut.UniversalTranslator

type CustomValidator struct {
    lock sync.Mutex
    validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
    return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
    cv.lock.Lock()
    zhTrans := zh.New()
    enTrans := en.New()
    uni = ut.New(zhTrans, zhTrans, enTrans)
    trans, _ := uni.GetTranslator("zh")
    _ = zh_trans.RegisterDefaultTranslations(cv.validator, trans)
    err := cv.validator.Struct(i)
    cv.lock.Unlock()
    return err
}
