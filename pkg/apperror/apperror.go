package apperror

import "errors"

// 应用错误类型，包含一个月湖错误码和对外提示
type AppError struct {
	Code string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code, msg string) *AppError {
	return &AppError{
		Code: code,
		Message: msg,
	}
}

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}