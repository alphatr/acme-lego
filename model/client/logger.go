package client

import (
	"strings"

	"github.com/alphatr/acme-lego/common/bootstrap"
)

type clientLogger struct{}

func (log *clientLogger) Fatal(args ...interface{}) {
	bootstrap.Log.Fatal(args)
}

func (log *clientLogger) Fatalln(args ...interface{}) {
	bootstrap.Log.Fatalln(args)
}

func (log *clientLogger) Fatalf(format string, args ...interface{}) {
	bootstrap.Log.Fatalf(format, args)
}

func (log *clientLogger) Print(args ...interface{}) {
	bootstrap.Log.Warn(args)
}

func (log *clientLogger) Println(args ...interface{}) {
	bootstrap.Log.Warningln(args)
}

func (log *clientLogger) Printf(format string, args ...interface{}) {
	if strings.HasPrefix(format, "[INFO] ") {
		format = strings.Replace(format, "[INFO]", "[lego]", -1)
		bootstrap.Log.Debugf(format, args)
		return
	}

	if strings.HasPrefix(format, "[WARN] ") {
		format = strings.Replace(format, "[WARN]", "[lego]", -1)
		bootstrap.Log.Warnf(format, args)
		return
	}

	bootstrap.Log.Warnf("[lego] "+format, args)
}
