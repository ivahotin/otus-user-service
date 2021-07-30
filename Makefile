ENTRYPOINT?=./cmd/apiserver
BINARY?=./build/apiserver

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