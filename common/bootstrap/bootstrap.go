package bootstrap

import (
	"github.com/alphatr/acme-lego/common/errors"
)

// Handler 注册额外的启动实例
type Handler interface {
	Init() *errors.Error
	Name() string
}

var extendInit = []Handler{}

// InitBootstrap 初始化 Bootstrap
func InitBootstrap() *errors.Error {
	logger, err := NewLogger()
	if err != nil {
		return errors.NewError(errors.BootstrapInitLoggerErrno, err)
	}

	Log = logger
	for _, handler := range extendInit {
		if err := handler.Init(); err != nil {
			return errors.NewError(errors.BootstrapInitHandlerErrno, err, handler.Name())
		}
	}

	return nil
}

// RegisterHandle 注册启动项
func RegisterHandle(handler Handler) {
	extendInit = append(extendInit, handler)
}
