package server

import (
	"backend/api"
	"backend/pkg/jwt"
	"backend/service/auth"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrUnableToParseJSON   = errors.New("unable to parse json")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrVerifyMissingParams = errors.New("verify missing params")
)

var errStatusCode = map[error]int{
	// bad request
	auth.ErrAccountNotFound:     http.StatusBadRequest,
	auth.ErrAccountNotVerified:  http.StatusBadRequest,
	auth.ErrWrongPassword:       http.StatusBadRequest,
	ErrVerifyMissingParams:      http.StatusBadRequest,
	auth.ErrAccountAlreadyExist: http.StatusBadRequest,

	ErrUnableToParseJSON: http.StatusBadRequest,

	// unauthorized
	ErrUnauthorized: http.StatusUnauthorized,

	// internal server error
	auth.ErrFailedToLogin:    http.StatusInternalServerError,
	auth.ErrFailedToVerify:   http.StatusInternalServerError,
	auth.ErrFailedToRegister: http.StatusInternalServerError,
	jwt.ErrFailToSignedToken: http.StatusInternalServerError,
}

// possible to translate this to Bahasa
var errMessage = map[error]string{
	auth.ErrFailedToLogin:       "failed to login",
	auth.ErrAccountNotFound:     "account not found",
	auth.ErrAccountNotVerified:  "account hasn't been verified yet",
	auth.ErrWrongPassword:       "wrong password",
	auth.ErrFailedToVerify:      "failed to verify",
	auth.ErrFailedToRegister:    "failed to register",
	auth.ErrAccountAlreadyExist: "account already exist",

	ErrVerifyMissingParams:   "verify missing params",
	ErrUnableToParseJSON:     "unable to parse json",
	jwt.ErrFailToSignedToken: "fail to sign token",
	ErrUnauthorized:          "unauthorized",
}

func httpError(w http.ResponseWriter, err error) {
	statusCode := errStatusCode[err]

	message := "something went wrong" // default message
	v, ok := errMessage[err]
	if ok {
		message = v
	}

	w.WriteHeader(statusCode)
	if statusCode == http.StatusUnauthorized {
		json.NewEncoder(w).Encode(api.ErrorUnAuthorizedResponse{
			Error: api.ErrorMessage(message),
		})
		return
	} else if statusCode == http.StatusBadRequest {
		json.NewEncoder(w).Encode(api.ErrorBadRequestResponse{
			Error: api.ErrorMessage(message),
		})
		return
	} else {
		json.NewEncoder(w).Encode(api.ErrorInternalServerResponse{
			Error: api.ErrorMessage(message),
		})
		return
	}
}
