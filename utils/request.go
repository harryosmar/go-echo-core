package utils

import (
	"bytes"
	"io"
	"net/http"
)

func GetCopyPayloadFromRequest(req *http.Request) ([]byte, error) {
	var err error
	// Duplicate the request body
	var bodyCopy bytes.Buffer
	if req.Body == nil {
		return nil, err
	}
	_, err = bodyCopy.ReadFrom(req.Body)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(bodyCopy.Bytes()))
	return bodyCopy.Bytes(), nil
}
