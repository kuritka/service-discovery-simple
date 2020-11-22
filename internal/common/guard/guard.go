// Package guard panics when error occurs
package guard

import (
	"fmt"
	"html"
	"net/http"
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

// FailOnError panics when error occurs.
func HandleErrorWithInternalServerError(w http.ResponseWriter, r *http.Request, err error, message string) {
	if err != nil {
		_, _ = fmt.Fprintf(w, message)
		w.WriteHeader(http.StatusInternalServerError)
		log.Logger().Errorf("%s: %s", html.EscapeString(r.URL.Path), err.Error())
	}
}

func HandleErrorWithNotFound(w http.ResponseWriter, r *http.Request, err error, message string) {
	if err != nil {
		_, _ = fmt.Fprintf(w, message)
		w.WriteHeader(http.StatusNotFound)
		log.Logger().Errorf("%s: %s", html.EscapeString(r.URL.Path), err.Error())
	}
}
