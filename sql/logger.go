package sql

import (
	"github.com/axonfibre/fibre.go/log"
)

type sqlLogger struct {
	log.Logger
}

func newLogger(logger log.Logger) *sqlLogger {
	return &sqlLogger{
		Logger: logger,
	}
}

func (l *sqlLogger) Printf(t string, args ...interface{}) {
	l.LogWarnf(t, args...)
}
