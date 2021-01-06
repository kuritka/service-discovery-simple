package discovery

import "net/url"

// Settings
type Settings struct {
	// Port where disco service is running
	Port int
	// URL to yaml file in GitHub repo
	YamlURL *url.URL
	// Sealed secret
	SealedSecret string
}
