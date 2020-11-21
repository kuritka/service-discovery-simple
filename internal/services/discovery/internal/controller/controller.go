package controller

import (
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/kuritka/k8gb-discovery/internal/services/discovery/internal/disco"
)

var router *httprouter.Router

func init() {
	router = httprouter.New()
}

type DiscoController struct {
	d *disco.Disco
}

func Startup(yamlURL *url.URL) *DiscoController {
	ctrl := &DiscoController{
		d: disco.NewDisco(yamlURL),
	}
	ctrl.register()
	return ctrl
}

func (c *DiscoController) Router() *httprouter.Router {
	return router
}

func (c *DiscoController) register() {
	// register routes here
	router.GET("/healthy", c.handleHealthy)
	router.GET("/disco", c.handleDiscovery)
}

func (c *DiscoController) handleHealthy(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func (c *DiscoController) handleDiscovery(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.NotFound(w, r)
}
