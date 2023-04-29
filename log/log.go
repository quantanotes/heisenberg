package log

import "github.com/golang/glog"

type M = map[string]interface{}

func Info(msg string, args M) {
	glog.Info(msg, args)
}

func Warning(msg string, args M) {
	glog.Warning(msg, args)
}

func Error(msg string, args M) {
	glog.Error(msg, args)
}

func Fatal(msg string, args M) {
	glog.Fatalf(msg, args)
}
