package util

import (
	"strconv"
)

var (
	ErrorNotLogin         = NewError(1, "请先登录")
	ErrorTopicNotFound    = NewError(-1, "话题不存在")
	ErrorArticleNotFound  = NewError(-2, "文章不存在")
	ErrorTagNotFound      = NewError(-3, "标签不存在")
	ErrorCaptchaWrong     = NewError(1000, "验证码错误")
	ErrorPermissionDenied = NewError(-100, "Permission denied.")
)

func NewError(code int, text string) *CodeError {
	return &CodeError{code, text, nil}
}

func NewErrorMsg(text string) *CodeError {
	return &CodeError{-1, text, nil}
}

func NewErrorData(code int, text string, data interface{}) *CodeError {
	return &CodeError{code, text, data}
}

func FromError(err error) *CodeError {
	if err == nil {
		return nil
	}
	return &CodeError{-1, err.Error(), nil}
}

type CodeError struct {
	Code    int
	Message string
	Data    interface{}
}

func (e *CodeError) Error() string {
	return strconv.Itoa(e.Code) + ": " + e.Message
}
