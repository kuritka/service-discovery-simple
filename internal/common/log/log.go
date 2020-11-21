// Package log provides logger in global scope
// see: https://github.com/sirupsen/logrus
package log

import (
	"github.com/sirupsen/logrus"
)

// log is the global logger.
var log *logrus.Logger

// Logger returns the global logger instance
func Logger() *logrus.Logger {
	return log
}

// run once on startup
func init() {
	log = logrus.New()
}
