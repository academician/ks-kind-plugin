package main

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	kindlog "sigs.k8s.io/kind/pkg/log"
)

type logWrapper struct {
	logger hclog.Logger
}

func (l *logWrapper) Warn(message string)               { l.logger.Warn(message) }
func (l *logWrapper) Warnf(format string, args ...any)  { l.Warn(fmt.Sprintf(format, args...)) }
func (l *logWrapper) Error(message string)              { l.logger.Error(message) }
func (l *logWrapper) Errorf(format string, args ...any) { l.Error(fmt.Sprintf(format, args...)) }
func (l *logWrapper) V(kindlevel kindlog.Level) kindlog.InfoLogger {
	level := hclog.Debug
	switch {
	case kindlevel <= 0:
		level = hclog.Info
	case kindlevel >= 3:
		level = hclog.Trace
	}
	return infoLogger{
		l.logger,
		level,
	}
}

// infoLogger implements [kindlog.InfoLogger] for [hclog.Logger]
type infoLogger struct {
	logger hclog.Logger
	level  hclog.Level
}

func (i infoLogger) Enabled() bool                    { return i.logger.GetLevel() <= i.level }
func (i infoLogger) Infof(format string, args ...any) { i.Info(fmt.Sprintf(format, args...)) }
func (i infoLogger) Info(message string) {
	if !i.Enabled() {
		return
	}
	switch i.level {
	case hclog.Info:
		i.logger.Info(message)
	case hclog.Trace:
		i.logger.Trace(message)
	default:
		i.logger.Debug(message)
	}
}
