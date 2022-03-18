package logging

import (
	"go.uber.org/zap/zapcore"
	"sync"
)

type (
	// hookable logger container
	loggerWithHookImpl struct {
		loggerImpl
		hookerLock sync.RWMutex
		hooker     LoggerHook
	}
)

func (l *loggerWithHookImpl) hook(entry zapcore.Entry) (err error) {
	l.hookerLock.RLock()
	defer l.hookerLock.RUnlock()
	if l.hooker != nil {
		err = l.hooker(entry.Level, entry.LoggerName, entry.Message, entry.Time)
	}
	return
}

func (l *loggerWithHookImpl) SetHook(hook LoggerHook) (err error) {
	l.hookerLock.Lock()
	defer l.hookerLock.Unlock()
	l.hooker = hook
	if hook == nil {
		l.Info("log hook cleared")
	} else {
		l.Info("log hook set")
	}
	return
}
