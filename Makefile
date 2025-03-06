generate-api:
	go generate cmd/api/main.go

generate-worker:
	go generate cmd/worker/main.go

api:
	@set -a && . env/local.env && set +a && go run cmd/api/*

worker:
	@set -a && . env/local.env && set +a && go run cmd/worker/*

sql-dev:
	./cloud_sql_proxy -instances=example-dev:us-central1:example=tcp:5432

sql-prod:
	./cloud_sql_proxy -instances=example:us-central1:example=tcp:5432

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

DATE_TAG=$(shell date +'%Y.%m.%d')
HASH=$(shell git rev-parse --short HEAD)
FULL_TAG=$(DATE_TAG)-$(HASH)

.PHONY: tag
tag:
	@echo "Pinning version tag..."
	git tag $(FULL_TAG)
	git push origin $(FULL_TAG)
