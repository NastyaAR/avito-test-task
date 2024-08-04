package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
)

type ErrorResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

const (
	ReadHTTPBodyError = iota + 1
	UnmarshalHTTPBodyError
)

const (
	ReadHTTPBodyMsg      = "can't read request"
	UnmarshalHTTPBodyMsg = "can't unmarshal request"
)

//type CustomError struct {
//
//}

//var ErrReflect map[int]string

func CreateErrorResponse(ctx context.Context, errCode int, msg string) []byte {
	var errResponse ErrorResponse
	errResponse.Code = errCode
	errResponse.RequestID = middleware.GetReqID(ctx)
	errResponse.Message = msg

	response, err := json.Marshal(errResponse)
	if err != nil {
		return nil
	}

	return response
}
