# k8gb discovery service

multi-cluster k8gb discovery service

## Project Health

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview
Service provides configuration to particular GSLB instances during GSLB startup.
This solution is useful if you can't provide various configurations during deployment. 

![](https://github.com/kuritka/trash/blob/master/k8gb-discovery-service.png?raw=true)

### Environment variables

| Name | Description | Default |
| --- | --- | --- |
| `K8GB_DISCOVERY_YAML_URL` | (Required) URL to raw yaml configuration | |
| `K8GB_DISCOVERY_EXPOSED_PORT` | (Optional) Service listener port | `8080` |
| `K8GB_DISCOVERY_DURATION` | (Optional) Duration in case you decide to poll yaml configuration <`3m`; `24h`> |  |


### REST-API

| Name | Description |
| --- | --- |
| `/healthy` | In case you establish liveness probe |
| `/discover/:key` | GSLB hits that endpoint to get configuration where key is unique value provided by GSLB |
| `/restore` | Restores cache from raw YAML (`K8GB_DISCOVERY_YAML_URL`) |
| `/metrics` | simple metrics |

### example YAML configuration
```yaml
test-gslb-us: #can I use unique key for particular k8gb instances ? In the worst case I can combine <cluster>:<namespace>:<instance>
  clusterGeoTag: us
  extGslbClustersGeoTags:
    - eu
  dnsZone: cloud.example.com
  ingressNamespace: k8gb
  edgeDNSZone: example.com
  edgeDNSServer: 1.1.1.1
test-gslb-eu:
  cluster: test-gslb1 # do I need this? isn't enough key e.g. test-gslb-eu
  clusterGeoTag: eu
  extGslbClustersGeoTags:
    - us
  dnsZone: cloud.example.com
  ingressNamespace: k8gb
  edgeDNSZone: example.com
  edgeDNSServer: 1.1.1.1
```

## local playground
update docker and install local [k3d](https://k3d.io/) 
```
echo "127.0.0.1 disco.example.com" >> /etc/hosts 
make start
make redeploy
curl http://disco.example.com:8081/healthy
curl https://disco.example.com:8081/healthy
make stop
```

## TODO
 - [ ] RBAC
 - [ ] HELM chart 
 - [ ] tests coverage
 - [ ] wait for `done->` channel in cache 

