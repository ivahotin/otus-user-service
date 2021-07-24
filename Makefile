.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: run
run:
	go build -v ./cmd/apiserver
	./apiserver

.PHONY: clean
clean:
	rm apiserver

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