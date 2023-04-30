package log

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

// Shamelessly copied and pasted directly from AnthonyGG. Thanks!

type M map[string]any

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

func LogErrNilReturn[T interface{}](at string, err error) (*T, error) {
	e := fmt.Errorf("@%s, %v", at, err)
	Error(e.Error(), nil)
	return nil, e
}

func LogErrReturn(at string, err error) error {
	e := fmt.Errorf("@%s, %v", at, err)
	Error(e.Error(), nil)
	return e
}

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}
