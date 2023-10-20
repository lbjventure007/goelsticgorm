package model

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type User struct {
	Username string `json:"username" form:"username" path:"username" validate:"required,min=6,max=20"`
	Password string `json:"password" form:"password" path:"password" validate:"required,min=6,max=20"`
}

func (User) TableName() string {
	return "user"
}
func (u *User) Login() (*User, string) {
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")

	//实例化验证器
	validate := validator.New()
	// 注册翻译器到校验器
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return u, err.Error()
	}
	errs := validate.Struct(u)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			return u, err.Translate(trans)
		}
	}
	return u, ""
}
