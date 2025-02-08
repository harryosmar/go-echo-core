package presentation

import (
	"fmt"
)

type (
	Content struct {
		Status    bool                   `json:"status"`
		Message   string                 `json:"message,omitempty"`
		ErrorCode string                 `json:"code,omitempty"`
		Error     string                 `json:"error,omitempty"`
		Data      interface{}            `json:"data,omitempty"`
		MetaData  map[string]interface{} `json:"metadata,omitempty"`
	}

	Response struct {
		Headers    map[string]string `json:"headers"`
		Content    Content           `json:"content"`
		StatusCode int               `json:"status"`
	}

	Paginator struct {
		Page    int
		PerPage int
		Total   int64
	}
)

func (r Response) Error() string {
	return fmt.Sprintf("%+v", map[string]interface{}{
		"message":    r.Content.Message,
		"error":      r.Content.Error,
		"error_code": r.Content.ErrorCode,
	})
}

func NewResponseEntity() *Response {
	return &Response{
		Content: Content{
			MetaData: map[string]interface{}{},
		},
		Headers: map[string]string{},
	}
}

func (r Response) WithHeaders(headers map[string]string) Response {
	for k, v := range headers {
		r.Headers[k] = v
	}
	return r
}

func (r Response) WithData(data interface{}) Response {
	if data != nil {
		r.Content.Data = data
	}
	return r
}

func (r Response) WithStatusCode(statusCode int) Response {
	r.StatusCode = statusCode
	return r
}

func (r Response) WithContentStatus(status bool) Response {
	r.Content.Status = status
	return r
}

func (r Response) WithMessage(message string, args ...interface{}) Response {
	if args == nil {
		r.Content.Message = message
		return r
	}
	r.Content.Message = fmt.Sprintf(message, args...)
	return r
}

func (r Response) WithErrorCode(errorCode string) Response {
	r.Content.ErrorCode = errorCode
	return r
}

func (r Response) WithError(err error) Response {
	r.Content.Error = err.Error()
	return r
}

func (r Response) WithMetaData(meta map[string]interface{}) Response {
	for k, v := range meta {
		r.Content.MetaData[k] = v
	}
	return r
}

func (r Response) WithPaginator(paginator *Paginator) Response {
	r.Content.MetaData["page"] = paginator.Page
	r.Content.MetaData["per_page"] = paginator.PerPage
	r.Content.MetaData["total"] = paginator.Total

	return r
}

func (r Response) WithMetas(metas ...map[string]interface{}) Response {
	if metas == nil || len(metas) == 0 {
		return r
	}

	for _, meta := range metas {
		for k, v := range meta {
			r.Content.MetaData[k] = v
		}
	}
	return r
}
