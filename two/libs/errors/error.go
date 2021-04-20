package _err

import "fmt"

type AppError struct {
	Code   int32  // 错误码(给前端)
	Detail string // 错误信息
}

func (e *AppError) Error() string {
	return e.Detail
}

func New(code int32, msg ...interface{}) *AppError {
	err := new(AppError)
	err.Code = code

	if len(msg) > 0 {
		err.Detail = format(msg...)
	}

	return err
}

func format(msg ...interface{}) string {
	if fm, ok := msg[0].(string); ok {
		return fmt.Sprintf(fm, msg[1:]...)
	}
	return ""
}
