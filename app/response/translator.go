package response

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var translator ut.Translator

func init() {
	english := en.New()
	chinese := zh.New()
	uni := ut.New(english, chinese)

	var found bool
	translator, found = uni.GetTranslator("zh")
	if !found {
		panic(fmt.Errorf("translator [%s] not found", "en"))
	}

	//获取gin的validator实例
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := zh_translations.RegisterDefaultTranslations(validate, translator); err != nil {
			panic(fmt.Errorf("register default translation fail:%w", err))
		}

		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func translateErrors(validationErrs validator.ValidationErrors) []FieldError {
	errs := make([]FieldError, 0, len(validationErrs))
	for _, e := range validationErrs {
		var field string
		segments := strings.Split(e.Namespace(), ".")
		if len(segments) == 2 {
			field = segments[1]
		}
		errs = append(errs, FieldError{
			Message: e.Translate(translator),
			Field:   field,
			Value:   e.Value(),
		})
	}
	return errs
}
