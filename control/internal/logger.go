package internal

import (
	"log"
)

type Logger struct {
	Debug bool
}

func (l Logger) Debugf(format string, args ...interface{}) {
	if l.Debug {
		log.Printf(format+"\n", args...)
	}
}

func (l Logger) Infof(format string, args ...interface{}) {
	if l.Debug {
		log.Printf(format+"\n", args...)
	}
}

func (l Logger) Warnf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}
