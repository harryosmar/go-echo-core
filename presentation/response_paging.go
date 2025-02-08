package presentation

import "github.com/labstack/echo/v4"

func WritePaging[T any](c echo.Context, statusCode int, list []T, paginator *Paginator, metas ...map[string]interface{}) error {
	result := make([]interface{}, len(list))
	for i, v := range list {
		result[i] = v
	}

	responseEntity := NewResponseEntity().
		WithStatusCode(statusCode).
		WithContentStatus(true).
		WithData(TransformListAny(result)).
		WithPaginator(paginator).
		WithMetas(metas...)

	return responseEntity.WriteJson(c)
}
