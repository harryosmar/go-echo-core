package presentation

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func ResponseOk(statusCode int, data interface{}, metas ...map[string]interface{}) error {
	response := NewResponseEntity().
		WithStatusCode(statusCode).
		WithContentStatus(true).
		WithData(data).
		WithMetas(metas...)

	if transformer, ok := data.(Transformable); ok {
		response = response.WithData(transformer.Transform())
	}

	return response
}

func WriteResponseCreated(c echo.Context, statusCode int, data interface{}) error {
	err := ResponseOk(http.StatusCreated, data)
	responseEntity, _ := err.(Response)
	return responseEntity.WriteJson(c)
}

func WriteResponseOk(c echo.Context, data interface{}) error {
	err := ResponseOk(http.StatusOK, data)
	responseEntity, _ := err.(Response)
	return responseEntity.WriteJson(c)
}
