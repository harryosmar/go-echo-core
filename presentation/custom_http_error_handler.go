package presentation

import (
	"fmt"
	coreCtx "github.com/harryosmar/go-echo-core/context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	// if it's presentation error
	if re, ok := err.(Response); ok {
		_ = re.WriteJson(c)
		return
	}

	ctx := c.Request().Context()
	code := http.StatusInternalServerError

	defer func() {
		coreCtx.NewContextBuilder(ctx).GetLogger().WithField("status", code).Error(err)
	}()

	// general http error
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msgStr := ""
		if he.Message != nil {
			msgStr = fmt.Sprintf("%s", he.Message)
		}
		_ = NewResponseEntity().
			WithStatusCode(code).
			WithMessage(msgStr).
			WriteJson(c)
		return
	}

	// custom err
	responseErr := ResponseErr(err)
	responseEntity, _ := responseErr.(Response)
	code = responseEntity.StatusCode
	_ = responseEntity.WriteJson(c)
}
