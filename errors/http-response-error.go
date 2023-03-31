package errors

import "fmt"

type HttpResponseError struct {
	Code int
	Msg  string
}

func (e *HttpResponseError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Msg != "" {
		return fmt.Sprintf("请求返回错误代码(%d), 错误原因: %s", e.Code, e.Msg)
	} else {
		return fmt.Sprintf("请求返回错误代码(%d)", e.Code)
	}
}
