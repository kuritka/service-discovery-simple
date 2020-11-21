package controller

import (
	"net/http"
	"net/url"

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
	log.Logger().Infof("fetching configuration from %s", yamlURL.String())
	err = ctrl.cache.RestoreCache()
	if err != nil {
		return nil, err
	}
	log.Logger().Infof("done")
	return
}

func (c *DiscoController) Router() *httprouter.Router {
	return router
}

func (c *DiscoController) registerRoutes() {
	// registerRoutes routes here
	router.GET("/healthy", c.handleHealthy)
	router.GET("/disco", c.handleDiscovery)
}

func (c *DiscoController) handleHealthy(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func (c *DiscoController) handleDiscovery(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.NotFound(w, r)
}
