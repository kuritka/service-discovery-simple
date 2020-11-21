package data

import "time"

type K8gb struct {
	LastHit         time.Time
	HitCount        int
	GeoTag          string
	ExternalGeoTags []string
	//TODO: fill the rest
}
