ENTRYPOINT?=./cmd/apiserver
BINARY?=./build/apiserver

.PHONY: install
install:
	kubectl create ns user-service
	helm upgrade --install -n user-service -f infra/user-service/values.yaml user-service infra/user-service/.
	kubectl create ns auth-service
	helm upgrade --install -n auth-service -f infra/auth-service/values.yaml auth-service infra/auth-service/.
	kubectl apply -f infra/api-gateway/ingress.yaml

.PHONY: uninstall
uninstall:
	helm uninstall user-service -n user-service
	helm uninstall auth-service -n auth-service
	kubectl delete ns user-service
	kubectl delete ns auth-service

.PHONY: build
build: clean
	cd user-service; go build -o ${BINARY} -v ${ENTRYPOINT}

.PHONY: run
run: build
	cd user-service; go build -o ${BINARY} -v ${ENTRYPOINT}
	cd user-service; ./build/apiserver

.PHONY: clean
clean:
	cd user-service; rm -f ${BINARY}

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