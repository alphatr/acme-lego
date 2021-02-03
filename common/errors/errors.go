package errors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

// ErrorNum 错误代码类型
type ErrorNum uint

// ErrorContent 错误调用信息
type ErrorContent struct {
	Log   string
	Errno ErrorNum
}

// Error 错误信息
type Error struct {
	Content ErrorContent
	Parent  *Error
	Origin  error
	Data    []interface{}
	Level   logrus.Level
}

// NewError 返回错误
func NewError(errno ErrorNum, parent error, argus ...interface{}) *Error {
	content, ok := ErrorMap[errno]
	if !ok {
		return NewError(UnknowErrno, nil, errno)
	}

	content.Errno = errno
	resultError := &Error{Content: content, Data: argus, Level: logrus.WarnLevel}

	ErrorParent, ok := parent.(*Error)
	if ok {
		resultError.Parent = ErrorParent
		resultError.Level = ErrorParent.Level
	}

	ErrorOrigin, ok := parent.(error)
	if ok {
		resultError.Origin = ErrorOrigin
	}

	return resultError
}

func (err *Error) Error() string {
	output := []string{}
	var errno ErrorNum

	var current = err
	for {
		output = append(output, Format(current.Content.Log, current.Data...))
		if current.Parent != nil {
			current = current.Parent
		} else {
			errno = current.Content.Errno
			if current.Origin != nil {
				output = append(output, current.Origin.Error())
			}

			break
		}
	}

	return fmt.Sprintf("[%d] %s;", uint(errno), strings.Join(output, "; "))
}

// SetLevel 设置错误等级
func (err *Error) SetLevel(level logrus.Level) *Error {
	err.Level = level
	return err
}

// Format 格式化
func Format(tpl string, params ...interface{}) string {
	pattern := regexp.MustCompile(`(^|[^%])%[A-Za-z]`)
	all := pattern.FindAll([]byte(tpl), -1)

	// TODO: len(params) < len(all) 的错误处理
	return fmt.Sprintf(tpl, params[0:len(all)]...)
}
