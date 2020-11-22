// package contains Disco cache
package cache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"gopkg.in/yaml.v2"
)

type Cache struct {
	info     CacheInfo
	yamlURL  *url.URL
	cache    map[string]K8gb
	duration time.Duration
}

type CacheInfo struct {
	ValidFrom    time.Time
	RefreshCount int
}

type K8gb struct {
	LastHit          time.Time
	ValidFrom        time.Time
	HitCount         int
	GeoTag           string   `yamlURL:"clusterGeoTag"`
	ExternalGeoTags  []string `yamlURL:"extGslbClustersGeoTags"`
	DNSZone          string   `yamlURL:"dnsZone"`
	EdgeDNSZone      string   `yamlURL:"edgeDNSZone"`
	EdgeDNSServer    string   `yamlURL:"edgeDNSServer"`
	IngressNamespace string   `yamlURL:"ingressNamespace"`
}

func NewCache(yaml *url.URL) *Cache {
	return &Cache{
		cache:    make(map[string]K8gb),
		duration: time.Hour,
		info: CacheInfo{
			RefreshCount: 0,
			ValidFrom:    time.Now(),
		},
		yamlURL: yaml,
	}
}

func (s *Cache) RestoreCache() (err error) {
	var t = time.Now()
	bytes, err := s.getYAML()
	err = yaml.Unmarshal(bytes, &s.cache)
	for _, k8gb := range s.cache {
		k8gb.ValidFrom = t
	}
	s.info.RefreshCount++
	s.info.ValidFrom = t
	return
}

func (s *Cache) getYAML() (b []byte, err error) {
	var client http.Client
	var resp *http.Response
	resp, err = client.Get(s.yamlURL.String())
	if err != nil {
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("%s: %s returns %d", http.MethodGet, s.yamlURL.String(), resp.StatusCode)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	return
}

func (s *Cache) Info() CacheInfo {
	return s.info
}
