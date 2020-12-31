VERSION ?= 0.0.1
CLUSTER_NAME ?=k3s-k8gb-disco

.PHONY: lint
lint:
	$(call lint)

.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

.PHONY: docker-build
docker-build:
	$(call docker-build)

.PHONY: docker-deploy
docker-deploy:
	$(call docker-build)
	docker push docker.io/kuritka/k8gb-discovery:$(VERSION)

.PHONY: check
check:
	goimports -l -w ./
	go generate ./...
	go mod tidy
	$(call lint)
	$(call docker-build)

.PHONY: redeploy
redeploy:
	docker build -t docker.io/kuritka/k8gb-discovery:$(VERSION) .
	docker push docker.io/kuritka/k8gb-discovery:$(VERSION)
	kubectl delete ns k8gb-discovery
	kubectl apply -k ./app/base

.PHONY: start
start:
	k3d cluster create $(CLUSTER_NAME) --api-port 6550 -p "8080:80@loadbalancer"  -p "8443:443@loadbalancer" --agents 1 --k3s-server-arg "--no-deploy=traefik,metrics-server"
	kubectl create ns k8gb-discovery
	kubectl create ns cert-manager
	kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.1.0/cert-manager.yaml
	kubectl -n cert-manager wait --for=condition=Ready pod -l app.kubernetes.io/instance=cert-manager --timeout=30s

.PHONY: stop
stop:
	k3d cluster delete $(CLUSTER_NAME)

.PHONY: reset
reset: stop start

.PHONY: test-api
test-api:
	kubectl run -it --rm busybox --restart=Never --image=busybox -- sh -c \
	"wget -qO - k8gb-discovery.nonprod.bcp.absa.co.za/metrics"

define lint
	golangci-lint run
	go test  ./...
endef

define docker-build
	time docker build -t k8gb-discovery:$(VERSION) .
endef
