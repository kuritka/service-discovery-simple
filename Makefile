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
	#new instance of cluster requires to regenerate sealed secret, TODO: should be executed within start
	kubeseal --format yaml <sealed-secrets/disco-secret.yaml >app/base/sealed-secret.yaml
	kubectl apply -k ./app/base

.PHONY: start
start:
	k3d cluster create $(CLUSTER_NAME) --agents 1 -p "8443:443@loadbalancer" -p "8080:80@loadbalancer" --k3s-server-arg "--no-deploy=metrics-server"
	#k3d cluster create --config=k3d-config.yaml
	kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.4/controller.yaml
	kubectl create ns cert-manager
	kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.1.0/cert-manager.yaml
	kubectl -n cert-manager wait --for=condition=Ready pod -l app.kubernetes.io/instance=cert-manager --timeout=50s


.PHONY: stop
stop:
	k3d cluster delete $(CLUSTER_NAME)

.PHONY: reset
reset: stop start redeploy

.PHONY: test-api
test-api:
	kubectl run -it --rm busybox --restart=Never --image=busybox -- sh -c \
	"wget -qO - k8gb-discovery.nonprod.bcp.absa.co.za/metrics"

.PHONY: certificates
install-seal-secret:
	openssl genrsa -out sealed-secrets/certs/sealed-disco.example.com.pem 2048
	chmod 400 sealed-secrets/certs/sealed-disco.example.com.pem
	openssl req -new -key sealed-secrets/certs/sealed-disco.example.com.pem -out sealed-secrets/certs/sealed-disco.example.com.csr -config sealed-secrets/certs/sealed-disco.example.com.cnf
	openssl x509 -req -days 3650 -in sealed-secrets/certs/sealed-disco.example.com.csr -signkey sealed-secrets/certs/sealed-disco.example.com.pem -out sealed-secrets/certs/sealed-disco.example.com.crt

.PHONY: sealed-secrets
sealed-secrets:
	kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.4/controller.yaml
	echo "This is a secret!" | kubectl create secret generic disco-secret -n k8gb-discovery --dry-run=client --from-file=secret=/dev/stdin -o yaml > sealed-secrets/disco-secret.yaml
	# sealed-secret depends on cluster instance. You change cluster, so you need refresh sealed secret
	kubeseal --format yaml <sealed-secrets/disco-secret.yaml >sealed-secrets/sealed-secret.yaml
	kubectl apply -f sealed-secrets/sealed-secret.yaml
	@echo only sealed-secret goes to github, no secret
	@echo writing installed secret down:
	kubectl get secret disco-secret -n k8gb-discovery -o jsonpath="{.data.secret}" | base64 --decode

define lint
	golangci-lint run
	go test  ./...
endef

define docker-build
	time docker build -t k8gb-discovery:$(VERSION) .
endef
