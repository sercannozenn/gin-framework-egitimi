package api

import "github.com/gin-gonic/gin"

type APIErrorResponse struct {
	Message string              `json:"message"` // struct tag
	Errors  map[string][]string `json:"errors,omitempty"`
	// omitempty: eğer bu alan boş işse JSON'a dahil etme
}

type APISuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SendError (ctx *gin.Context, status int, message string, errs map[string][]string){
	ctx.JSON(status, APIErrorResponse{
		Message: message,
		Errors: errs,
	})
}

func SendSuccess(ctx *gin.Context, status int, message string, data interface{}){
	ctx.JSON(status, APISuccessResponse{
		Message: message,
		Data: data,
	})
}