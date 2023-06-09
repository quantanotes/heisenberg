package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

type M map[string]interface{}

type Level uint32

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func SetOutput(w io.Writer) {
	logrus.SetOutput(w)
}

func SetLevel(level Level) {
	var l logrus.Level
	switch level {
	case LevelTrace:
		l = logrus.TraceLevel
	case LevelInfo:
		l = logrus.InfoLevel
	case LevelWarn:
		l = logrus.WarnLevel
	case LevelError:
		l = logrus.ErrorLevel
	case LevelFatal:
		l = logrus.FatalLevel
	case LevelPanic:
		l = logrus.PanicLevel
	}
	logrus.SetLevel(l)
}

func Info(msg string, args M) {
	logrus.WithFields(logrus.Fields(args)).Info(msg)
}

func Debug(msg string, args M) {
	logrus.WithFields(logrus.Fields(args)).Debug(msg)
}

func Warn(msg string, args M) {
	logrus.WithFields(logrus.Fields(args)).Warn(msg)
}

func Error(msg string, args M) {
	logrus.WithFields(logrus.Fields(args)).Error(msg)
}

func Trace(msg string, args M) {
	logrus.WithFields(logrus.Fields(args)).Trace(msg)
}

func Fatal(msg string, args M) {
	logrus.WithFields(logrus.Fields(args)).Fatal(msg)
}

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}
