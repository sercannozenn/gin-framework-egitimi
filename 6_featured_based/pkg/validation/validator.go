package validation

import (
	"log"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	esTranslations "github.com/go-playground/validator/v10/translations/es"
	trTranslations "github.com/go-playground/validator/v10/translations/tr"
)

var translator ut.Translator

func Init(lang string) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("Validator engine not found")
	}
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get("json")
		if tag == "" {
			return field.Name
		}
		name := strings.Split(tag, ",")[0]
		if name == "-" || name == "" {
			return field.Name
		}

		return name
	})

	trLocale := tr.New()
	enLocale := en.New()
	esLocale := es.New()

	uni := ut.New(enLocale, enLocale, trLocale, esLocale)

	var found bool
	translator, found = uni.GetTranslator(lang)
	if !found {
		translator, _ = uni.GetTranslator("en")
	}

	switch lang {
	case "en":
		enTranslations.RegisterDefaultTranslations(v, translator)
	case "tr":
		trTranslations.RegisterDefaultTranslations(v, translator)
	case "es":
		esTranslations.RegisterDefaultTranslations(v, translator)
	default:
		enTranslations.RegisterDefaultTranslations(v, translator)
	}
	if err := enTranslations.RegisterDefaultTranslations(v, translator); err != nil {
		log.Printf("translation register error: %v", err)
	}
}

func MapValidationErrors(ve validator.ValidationErrors) map[string][]string {
	out := make(map[string][]string)

	for _, fe := range ve {
		// fe.Field() => struct alan adı Name Email Age
		// fe.Tag() => failed tag required min email
		// fe.Param() => min=2
		field := fe.Field() // name, email
		msg := fe.Translate(translator)

		out[field] = append(out[field], msg)
	}
	return out
}