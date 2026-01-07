package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
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

func main() {

	cfg := loadConfig()

	initValidator(cfg)

	r := gin.Default()

	r.POST("/users", createUserHandler)

	r.Run(":8080")
}

type APIErrorResponse struct {
	Message string              `json:"message"` // struct tag
	Errors  map[string][]string `json:"errors,omitempty"`
	//omitempty: eğer bu alan boş işse JSON'a dahil etme
	/*
		{
			"message": "Validation Failed",
			"errors": {
				"email": [
					"Email is required",
					"Email must be a valid email address"
				],
				"password": [
					"Password is required"
				]
			}

		}
	*/
}

type APISuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Request DTO: Bind + Validation
// DTO: Data Transfer Obkect
// Veri Transfer Nesnesi
// required, min,max, email gibi kuralları business validation tarafıdır
// Tip uyuşmazlığı (int beklerken string gelmesi) ise "bind/parsing" hatası
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required"`
}

type Config struct {
	Lang string // "en", es, tr
}

func loadConfig() Config {
	lang := strings.TrimSpace(strings.ToLower(os.Getenv("APP_LANG"))) // En en

	if lang == "" {
		lang = "en"
	}

	switch lang {
	case "tr", "en", "es":
		// bunlardan birisi geliyorsa lang olduğu gibi kalsın
	default:
		lang = "en"
	}

	return Config{Lang: lang}
}

func initValidator(cfg Config) {
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
	translator, found = uni.GetTranslator(cfg.Lang)
	if !found {
		translator, _ = uni.GetTranslator("en")
	}

	switch cfg.Lang {
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

func createUserHandler(c *gin.Context) {
	var req CreateUserRequest

	// 1- Json parse eder
	// 2- Struct alanlarına map eder
	// 3- binding taglerine göre doğrulama yapar

	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			// ve: validation errors
			// ok: doğrulama başarılı mı?

			c.JSON(http.StatusUnprocessableEntity, APIErrorResponse{
				Message: "Validation failed",
				Errors:  mapValidationErrors(ve),
			})
			return
		}

		c.JSON(http.StatusBadRequest, APIErrorResponse{
			Message: "Invalid request payload.",
		})
		return
	}
	// Buraya kadar kodumuz geldiuse
	// json doğrulanmıştır
	// tipler doğrudur
	// Validation kuralları sağlanmıştır.

	// business logic kontrolü
	if err := pretendDBInsert(req); err != nil {
		// örnek olarak DB insert hatası

		c.JSON(http.StatusInternalServerError, APIErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, APISuccessResponse{
		Message: "User created",
		Data: gin.H{
			"name":  req.Name,
			"email": req.Email,
			"age":   req.Age,
		},
	})
}

func mapValidationErrors(ve validator.ValidationErrors) map[string][]string {
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

func pretendDBInsert(req CreateUserRequest) error {
	if len(req.Email) >= 5 && req.Email[:5] == "fail@" {
		return errors.New("db insert failed")
	}

	return nil
}
