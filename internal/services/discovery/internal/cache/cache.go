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
	info     Info
	yamlURL  *url.URL
	cache    map[string]K8gb
	duration time.Duration
}

type Info struct {
	ValidFrom    time.Time
	RefreshCount int
}

type K8gb struct {
	LastHit          time.Time
	ValidFrom        time.Time
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
		cache:    make(map[string]K8gb),
		duration: time.Hour,
		info: Info{
			RefreshCount: 0,
			ValidFrom:    time.Now(),
		},
		yamlURL: yaml,
	}
}

func (s *Cache) Get(key string) (k8gb K8gb, err error) {
	var ok bool
	if k8gb, ok = s.cache[key]; ok {
		k8gb.HitCount++
		k8gb.LastHit = time.Now()
		k8gb.ValidFrom = s.info.ValidFrom
		s.cache[key] = k8gb
		return
	}
	err = fmt.Errorf("%s not found", key)
	return
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

func (s *Cache) Info() Info {
	return s.info
}
