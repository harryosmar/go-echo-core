package presentation

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"os"
)

func (r Response) WriteJson(c echo.Context) error {
	return WriteJson(r, c)
}

func WriteJson(r Response, c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	if r.Headers != nil {
		for k, v := range r.Headers {
			c.Response().Header().Set(k, v)
		}
	}
	c.Response().WriteHeader(r.StatusCode)

	if os.Getenv("APP_DEBUG") == "true" {
		r.Content.Error = ""
	}

	return json.NewEncoder(c.Response()).Encode(r.Content)
}
