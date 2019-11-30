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
	SystemError       = 9999
	InvalidParamError = 8001
	BusinessError     = 5001
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
