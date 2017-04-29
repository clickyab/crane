package controller

import (
	"io"

	echo "gopkg.in/labstack/echo.v3"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/gommon/log"
)

type myLogger struct {
	logger *logrus.Logger
}

func (m *myLogger) SetLevel(l log.Lvl) {
	switch l {
	case log.DEBUG:
		m.logger.Level = logrus.DebugLevel
	case log.INFO:
		m.logger.Level = logrus.InfoLevel
	case log.WARN:
		m.logger.Level = logrus.WarnLevel
	case log.ERROR:
		m.logger.Level = logrus.ErrorLevel
	case log.OFF:
		m.logger.Level = logrus.PanicLevel
	}
}
func (m *myLogger) Print(args ...interface{}) {
	m.logger.Print(args...)
}
func (m *myLogger) Printf(s string, a ...interface{}) {
	m.logger.Panicf(s, a...)
}
func (m *myLogger) Printj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Print("-")
}

func (m *myLogger) Debug(args ...interface{}) {
	m.logger.Debug(args...)
}
func (m *myLogger) Debugf(s string, args ...interface{}) {
	m.logger.Debugf(s, args...)
}
func (m *myLogger) Debugj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Debug("-")
}
func (m *myLogger) Info(args ...interface{}) {
	m.logger.Info(args...)
}
func (m *myLogger) Infof(s string, args ...interface{}) {
	m.logger.Infof(s, args...)
}
func (m *myLogger) Infoj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Info("-")
}
func (m *myLogger) Warn(args ...interface{}) {
	m.logger.Warn(args)
}
func (m *myLogger) Warnf(s string, args ...interface{}) {
	m.logger.Warnf(s, args...)
}
func (m *myLogger) Warnj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Warn("-")
}
func (m *myLogger) Error(args ...interface{}) {
	m.logger.Error(args...)
}
func (m *myLogger) Errorf(s string, args ...interface{}) {
	m.logger.Errorf(s, args...)
}
func (m *myLogger) Errorj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Error("-")
}
func (m *myLogger) Fatal(args ...interface{}) {
	m.logger.Fatal(args...)
}
func (m *myLogger) Fatalj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Fatal("-")
}
func (m *myLogger) Panicf(s string, args ...interface{}) {
	m.logger.Panicf(s, args...)
}

func (m *myLogger) Panic(args ...interface{}) {
	m.logger.Panic(args...)
}
func (m *myLogger) Panicj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Panic("-")
}
func (m *myLogger) Fatalf(s string, args ...interface{}) {
	m.logger.Fatalf(s, args...)
}
func (m *myLogger) Level() log.Lvl {
	switch m.logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.InfoLevel:
		return log.INFO
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.FatalLevel:
		return log.OFF
	case logrus.PanicLevel:
		return log.OFF
	}

	return log.DEBUG
}
func (m *myLogger) Output() io.Writer {
	return m.logger.Out
}
func (m *myLogger) SetOutput(i io.Writer) {
	m.logger.Out = i
}
func (m *myLogger) Prefix() string {
	return ""
}
func (m *myLogger) SetPrefix(string) {

}

// NewLogger return a dummy logger for echo
func NewLogger() echo.Logger {
	return &myLogger{logger: logrus.New()}
}
