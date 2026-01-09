package user

import (
	"errors"
	"featured-base-starter-kit/pkg/api"
	"featured-base-starter-kit/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateUserHandler(c *gin.Context) {
	var req CreateUserRequest

	// 1- Json parse eder
	// 2- Struct alanlarına map eder
	// 3- binding taglerine göre doğrulama yapar

	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			// ve: validation errors
			// ok: doğrulama başarılı mı?

			c.JSON(http.StatusUnprocessableEntity, api.APIErrorResponse{
				Message: "Validation failed",
				Errors:  validation.MapValidationErrors(ve),
			})
			return
		}

		c.JSON(http.StatusBadRequest, api.APIErrorResponse{
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

		c.JSON(http.StatusInternalServerError, api.APIErrorResponse{
			Message: "Internal server error",
		})
		return
	}

	format := c.Query("format")

	response := api.APISuccessResponse{
		Message: "User created",
		Data: gin.H{
			"name":  req.Name,
			"email": req.Email,
			"age":   req.Age,
		},
	}

	switch format {
	case "xml":
		c.XML(http.StatusCreated, response)
	case "yaml":
		c.YAML(http.StatusCreated, response)
	default:
		c.JSON(http.StatusCreated, response)
	}	
}

func pretendDBInsert(req CreateUserRequest) error {
	if len(req.Email) >= 5 && req.Email[:5] == "fail@" {
		return errors.New("db insert failed")
	}

	return nil
}