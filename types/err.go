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
	Success              = 0
	SystemError          = 9999
	InvalidParamError    = 9998
	RepeatRequest        = 9997
	NeedAuthError        = 9996
	DbError              = 9995
	NotFound             = 5001
	ExceedRateLimit      = 5002
	DataDirty            = 5003
	PayChannelNotSupport = 5004
)

func InvalidParamErr(msg string) error {
	return NewErr(InvalidParamError, msg)
}

func SystemErr(msg string) error {
	return NewErr(SystemError, msg)
}

func RepeatRequestErr(msg string) error {
	return NewErr(RepeatRequest, msg)
}

func NeedAuthErr(msg string) error {
	return NewErr(NeedAuthError, msg)
}

func NotFoundErr(msg string) error {
	return NewErr(NotFound, msg)
}

func DbErr(err error) error {
	return NewErr(DbError, err.Error())
}

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
