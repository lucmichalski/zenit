SHELL := /bin/bash
BUILD_DATE := `date +%Y%m%d%H%M`
DOCKER_COMPOSE_FILE := docker-compose --file docker/docker-compose.yml
DOCKER_EXEC := docker exec -i -t -u root

.PHONY: help

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	go get -u github.com/aws/aws-sdk-go
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/go-yaml/yaml
	go get -u github.com/hpcloud/tail
	go get -u github.com/influxdata/influxdb1-client/v2
	go get -u github.com/kardianos/service
	go get -u github.com/shirou/gopsutil
	go get -u github.com/sirupsen/logrus
	go get -u golang.org/x/sys/unix

deps-devel:
	brew install jq

tests: ## Run tests
	go generate ./...
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic ./...

build: ## Build binary for local operating system
	go generate ./...
	go build -ldflags "-s -w -X github.com/swapbyt3s/zenit/command.BuildTime=$(BUILD_DATE)" -o zenit main.go

build_linux: ## Build binary for Linux
	go generate ./...
	GOOS=linux go build -ldflags "-s -w -X github.com/swapbyt3s/zenit/command.BuildTime=$(BUILD_DATE)" -o zenit main.go

build-and-run: ## Build binary for Linux and run
	go generate ./...
	GOOS=linux go build -ldflags "-s -w -X github.com/swapbyt3s/zenit/command.BuildTime=$(BUILD_DATE)" -o zenit main.go
	docker cp zenit zenit_percona_server_primary:/usr/bin/
	docker exec -i -t -u root zenit_percona_server_primary /usr/bin/zenit

release: ## Create release
	scripts/release.sh

docker_build: ## Build docker images
	$(DOCKER_COMPOSE_FILE) build

docker_build_clickhouse: ## Build docker images for clickhouse
	$(DOCKER_COMPOSE_FILE) build clickhouse

docker_build_percona_server_primary: ## Build docker images for percona-server-primary
	$(DOCKER_COMPOSE_FILE) build percona-server-primary

docker_build_proxysql: ## Build docker images for proxysql
	$(DOCKER_COMPOSE_FILE) build proxysql

docker_up: ## Run docker-compose
	$(DOCKER_COMPOSE_FILE) --project-name=zenit up

docker_ps: ## Show status for all containers
	$(DOCKER_COMPOSE_FILE) --project-name=zenit ps

docker_down: ## Down docker-compose
	$(DOCKER_COMPOSE_FILE) --project-name=zenit down

docker_clickhouse: ## Enter into ClickHouse Client
	$(DOCKER_EXEC) zenit_clickhouse /usr/bin/clickhouse-client

docker_percona_primary: ## Enter in Percona Server Primary Console
	$(DOCKER_EXEC) zenit_percona_server_primary /usr/bin/mysql

docker_percona_primary_bash: ## Enter in Percona Server Primary bash console
	$(DOCKER_EXEC) zenit_percona_server_primary /bin/bash

docker_percona_secondary: ## Enter in Percona Server Secondary Console
	$(DOCKER_EXEC) zenit_percona_server_secondary /usr/bin/mysql

docker_percona_secondary_bash: ## Enter in Percona Server Secondary bash console
	$(DOCKER_EXEC) zenit_percona_server_secondary /bin/bash

docker_proxysql: ## Enter in ProxySQL Console
	$(DOCKER_EXEC) zenit_proxysql /usr/bin/mysql --socket=/tmp/proxysql_admin.sock -u proxysql -padmin  --prompt='ProxySQLAdmin> '

docker_proxysql_bash: ## Enter in ProxySQL bash console
	$(DOCKER_EXEC) zenit_proxysql /bin/bash

docker_influxdb_bash: ## Enter in InfluxDB bash console
	$(DOCKER_EXEC) influxdb /usr/bin/influx

docker_zenit_build: build_linux ## Build binary and copy to container
	docker cp zenit zenit_percona_server_primary:/usr/bin/
	docker cp zenit zenit_percona_server_secondary:/usr/bin/
	docker cp zenit zenit_proxysql:/usr/bin/
