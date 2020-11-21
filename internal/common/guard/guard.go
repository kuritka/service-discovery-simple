// Package guard panics when error occurs
package guard

import (
	"os"

	"github.com/kuritka/k8gb-discovery/internal/common/log"
)

// FailOnError panics when error occurs.
func FailOnError(err error, message string, a ...interface{}) {
	if err != nil {
		log.Logger().Error(err, message, a)
		os.Exit(1)
	}
}
