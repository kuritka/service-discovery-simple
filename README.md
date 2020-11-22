# k8gb discovery service

multi-cluster k8gb discovery

## Project Health

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Contents
tbd

QUATIONS:
 - security? Is that secret solution for PROD ? 
 - Can we provide unique key for each GSLB instance ? 

## Overview
tbd

## Usage
tbd

### Environment variables

| Name | Description | Default |
| --- | --- | --- |
| `K8GB_DISCOVERY_YAML_URL` | (Required) URL to raw yaml configuration | |
| `K8GB_DISCOVERY_PORT` | (Optional) Service listener port | `8080` |
| `K8GB_DISCOVERY_DURATION` | (Optional) Duration in case you decide to poll yaml configuration <`3m`; `24h`> |  |


### REST-API

| Name | Description |
| --- | --- |
| `/healthy` | In case you establish liveness probe |
| `/discover/:key` | GSLB hits that endpoint to get configuration where key is unique value provided by GSLB |
| `/restore` | Restores cache from raw YAML (`K8GB_DISCOVERY_YAML_URL`) |
| `/metrics` | simple metrics |
