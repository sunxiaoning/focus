package types

import "encoding/json"

type FocusError struct {
	Code    int
	Message string
}

const (
	ErrCode = "responseCode"
	ErrMsg  = "responseMessage"
)

const (
	Success           = 0
	SystemError       = 9999
	InvalidParamError = 9998
	RepeatRequest     = 9997
	NeedAuthError     = 9996
	UserNotFound      = 5001
	ExceedRateLimit   = 5002
)

func NewErr(code int, msg string) error {
	return &FocusError{
		Code:    code,
		Message: msg,
	}
}

func (focusError *FocusError) Error() string {
	errContent := make(map[string]interface{})
	errContent[ErrCode] = focusError.Code
	errContent[ErrMsg] = focusError.Message
	msg, err := json.Marshal(errContent)
	if err != nil {
		return err.Error()
	}
	return string(msg)
}
