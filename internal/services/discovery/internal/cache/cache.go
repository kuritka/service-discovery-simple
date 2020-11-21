// package contains Disco cache
package cache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ghodss/yaml"
	"github.com/kuritka/k8gb-discovery/internal/common/log"
)

type Cache struct {
	yaml         *url.URL
	cache        map[string]K8gb
	validFrom    time.Time
	duration     time.Duration
	refreshCount int
}

type K8gb struct {
	LastHit          time.Time
	HitCount         int
	GeoTag           string   `yaml:"clusterGeoTag"`
	ExternalGeoTags  []string `yaml:"extGslbClustersGeoTags"`
	DNSZone          string   `yaml:"dnsZone"`
	EdgeDNSZone      string   `yaml:"edgeDNSZone"`
	EdgeDNSServer    string   `yaml:"edgeDNSServer"`
	IngressNamespace string   `yaml:"ingressNamespace"`
}


func NewCache(yaml *url.URL) *Cache {
	return &Cache{
		yaml:         yaml,
		cache: make(map[string]K8gb),
		validFrom:    time.Now(),
		duration:     time.Hour,
		refreshCount: 0,
	}
}

func (s *Cache) RestoreCache() (err error) {
	var client http.Client
	var resp *http.Response
	resp, err = client.Get(s.yaml.String())
	if err != nil {
		return
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s: %s returns %d", http.MethodGet, s.yaml.String(), resp.StatusCode)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Logger().Errorf("unable to close response body %s", err.Error())
		}
	}()

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(bodyBytes, &s.cache)
	if err != nil {
		return
	}
	return fmt.Errorf("skmfkmnsdf")

	// TODO: refactor large function
	// TODO: fill rest of fields in cache items
	//return
}
