# discovery service demo
We need to build security into the architecture from day one. Sensitive information must be encrypted and test. 
Following demo presents usage of [cert manager](https://cert-manager.io/docs/) and [sealed secrets](https://github.com/bitnami-labs/sealed-secrets)
(maybe [medium article](https://medium.com/better-programming/encrypting-kubernetes-secrets-with-sealed-secrets-fe363149a211) is better).
The demo runs web application on the top of [k3d](https://k3d.io/) and intentionally provides functionality on http and https.

## Overview
Service provides configuration to particular [GSLB instances](https://github.com/AbsaOSS/k8gb) during GSLB startup.
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
bump docker version to the latest and install local [k3d](https://k3d.io/).
Certificate manager generates self-signed certificate `*.example.com` the transfered data is still sent encrypted, 
but `curl https://...` will require `-k/--insecure` argument which will "only make" curl skip certificate validation, 
it will not turn off SSL all together. depending on a browser you will need to skip 
[NET::ERR_CERT_INVALID](https://www.pandasecurity.com/en/mediacenter/panda-security/your-connection-is-not-private/)
error.
```
echo "127.0.0.1 disco.example.com" >> /etc/hosts 
make reset
curl http://disco.example.com:8080/healthy
curl https://disco.example.com:8443/healthy
make stop
```

To manipulate with sealed-secrets run :
```shell script
make sealed-secrets
```

## TODO
 - [ ] RBAC


