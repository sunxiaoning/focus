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
	PayOrderTimeout      = 5005
	FileSizeTooLarge     = 5006
	GenUUIDError         = 5007
)

func InvalidParamErr(msg string) *FocusError {
	return NewErr(InvalidParamError, msg)
}

func InvalidParamPanic(msg string) {
	ErrPanic(InvalidParamError, msg)
}

func SystemErr(msg string) *FocusError {
	return NewErr(SystemError, msg)
}

func SystemPanic(msg string) {
	ErrPanic(SystemError, msg)
}

func RepeatRequestErr(msg string) *FocusError {
	return NewErr(RepeatRequest, msg)
}

func RepeatRequestPanic(msg string) {
	ErrPanic(RepeatRequest, msg)
}

func NeedAuthErr(msg string) *FocusError {
	return NewErr(NeedAuthError, msg)
}

func NeedAuthPanic(msg string) {
	ErrPanic(NeedAuthError, msg)
}

func NotFoundErr(msg string) *FocusError {
	return NewErr(NotFound, msg)
}

func NotFoundPanic(msg string) {
	ErrPanic(NotFound, msg)
}

func DbErr(err error) *FocusError {
	return NewErr(DbError, err.Error())
}

func DbPanic(err error) {
	ErrPanic(DbError, err.Error())
}

func NewErr(code int, msg string) *FocusError {
	return &FocusError{
		Code:    code,
		Message: msg,
	}
}

func ErrPanic(code int, msg string) {
	NewErr(code, msg).Throw()
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

func (focusError *FocusError) Throw() {
	panic(focusError)
}
