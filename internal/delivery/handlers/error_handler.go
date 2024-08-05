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

type ErrorPair struct {
	Message   string
	ErrorCode int
}

const (
	ReadHTTPBodyError = iota + 1
	UnmarshalHTTPBodyError
	CreateHouseError
	MarshalHTTPBodyError
	ParseURLError
	GetFlatsByHouseIDError
	ReadCookieError
	NotAuthorizedError
	RegisterUserError
	LoginUserError
	DummyLoginError
	CreateFlatError
	UpdateFlatError
)

const (
	ReadHTTPBodyMsg           = "can't read request"
	UnmarshalHTTPBodyMsg      = "can't unmarshal request"
	CreateHouseErrorMsg       = "can't create house"
	MarshalHTTPBodyErrorMsg   = "can't marshal response"
	ParseURLErrorMsg          = "can't parse url"
	GetFlatsByHouseIDErrorMsg = "can't get flats by house id"
	ReadCookieErrorMsg        = "can't read cookie"
	NotAuthorizedErrorMsg     = "not authorized"
	RegisterUserErrorMsg      = "can't register user"
	LoginUserErrorMsg         = "can't login user"
	DummyLoginErrorMsg        = "can't simple login"
	CreateFlatErrorMsg        = "can't create flat"
	UpdateFlatErrorMsg        = "can't update flat"
)

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
