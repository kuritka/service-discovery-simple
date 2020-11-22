package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kuritka/k8gb-discovery/internal/common/guard"

	"github.com/kuritka/k8gb-discovery/internal/common/log"
	"github.com/kuritka/k8gb-discovery/internal/services/discovery/internal/cache"

	"github.com/julienschmidt/httprouter"
)

var router *httprouter.Router

func init() {
	router = httprouter.New()
}

type DiscoController struct {
	cache *cache.Cache
}

func Startup(yamlURL *url.URL) (ctrl *DiscoController, err error) {
	ctrl = &DiscoController{
		cache: cache.NewCache(yamlURL),
	}
	ctrl.registerRoutes()
	log.Logger().Infof("restoring cache from %s", yamlURL.String())
	err = ctrl.cache.RestoreCache()
	if err != nil {
		return nil, err
	}
	return
}

func (c *DiscoController) Router() *httprouter.Router {
	return router
}

func (c *DiscoController) registerRoutes() {
	// registerRoutes routes here
	router.GET("/healthy", c.handleHealthy)
	router.GET("/restore", c.handleRestore)
	router.GET("/discover", c.handleDiscovery)
	router.GET("/metrics", c.handleMetrics)

}

func (c *DiscoController) handleRestore(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := c.cache.RestoreCache()
	guard.HandleErrorWithInternalServerError(w, r, err, "Server error; contact your system administrator")
	s := fmt.Sprintf("cache restored %s", c.cache.Info().ValidFrom.Format(time.RFC822))
	_, _ = fmt.Fprintf(w, s)
	log.Logger().Info(s)
}

func (c *DiscoController) handleHealthy(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprintf(w, "healthy")
}

func (c *DiscoController) handleDiscovery(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.NotFound(w, r)
}

func (c *DiscoController) handleMetrics(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(c.cache.Info())
	guard.HandleErrorWithInternalServerError(w, r, err, "retrieving info")
}
