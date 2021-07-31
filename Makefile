ENTRYPOINT?=./cmd/apiserver
BINARY?=./build/apiserver

.PHONY: install
install:
	kubectl create ns monitoring
	kubectl config set-context --current --namespace=monitoring
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm upgrade --install -f infra/prometheus/prometheus.yaml prometheus prometheus-community/kube-prometheus-stack --atomic
	helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
	helm upgrade --install -f infra/prometheus/nginx-ingress.yaml nginx ingress-nginx/ingress-nginx --atomic
	kubectl create ns user-service
	kubectl config set-context --current --namespace=user-service
	helm upgrade --install -f infra/user-service/values.yaml user-service infra/user-service/.

.PHONY: build
build: clean
	go build -o ${BINARY} -v ${ENTRYPOINT}

.PHONY: run
run: build
	go build -o ${BINARY} -v ${ENTRYPOINT}
	./build/apiserver

.PHONY: clean
clean:
	rm -f ${BINARY}

.PHONY: local-run
local-run:
	docker-compose up -d

.PHONY: local-stop
local-down:
	docker-compose down

.PHONY: start-deps
start-deps:
	docker-compose -f docker-compose.deps.yaml up -d

.PHONY: stop-deps
stop-deps:
	docker-compose -f docker-compose.deps.yaml down

.DEFAULT_GOAL := build