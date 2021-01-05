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

// nolint[:lll] See https://longurl.org/path/to/thing/with/commit/has/1232u4t98u98fdugdsgeawefe/now/a/file.txt
const blankFavicon = "\"data:image/x-icon;base64,AAABAAEAEBAAAAAAAABoBQAAFgAAACgAAAAQAAAAIAAAAAEACAAAAAAAAAEAAAAAAAAAAAAAAAEAAAAAAAD///8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=\\n\""

var router *httprouter.Router

func init() {
	router = httprouter.New()
}

type DiscoController struct {
	cache *cache.Cache
	sealedSecret string
}

func Startup(yamlURL *url.URL, sealedSecret string) (ctrl *DiscoController, err error) {
	ctrl = &DiscoController{
		cache: cache.NewCache(yamlURL),
		sealedSecret: sealedSecret,
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
	router.GET("/sealed-secret", c.handleSealedSecret)
	router.GET("/healthy", c.handleHealthy)
	router.GET("/restore", c.handleRestore)
	router.GET("/discover/:key", c.handleDiscovery)
	router.GET("/metrics", c.handleMetrics)
	router.GET("/favicon.ico", c.handleFavicon)
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

func (c *DiscoController) handleSealedSecret(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprintf(w, c.sealedSecret)
}

func (c *DiscoController) handleDiscovery(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	key := p.ByName("key")
	k8gb, err := c.cache.Get(key)
	if err != nil {
		guard.HandleErrorWithNotFound(w, r, err, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(k8gb)
	guard.HandleErrorWithInternalServerError(w, r, err, "retrieving info")
}

func (c *DiscoController) handleMetrics(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(c.cache.Info())
	guard.HandleErrorWithInternalServerError(w, r, err, "retrieving info")
}

func (c *DiscoController) handleFavicon(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	_, _ = fmt.Fprintln(w, blankFavicon)
}
