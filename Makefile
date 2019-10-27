.PHONY: lint build push

lint:
	@if [ -z `which golangci-lint 2> /dev/null` ]; then \
		GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1; \
	fi
	@golangci-lint run --tests --disable-all --enable=goimports --enable=golint --enable=govet --enable=errcheck --enable=staticcheck
	@gofmt -s -w .

build:
	docker build -t read-replicas-outbox-miner .

tag:
	docker tag read-replicas-outbox-miner:latest $(REPOSITORY):latest

push: tag
	docker push $(REPOSITORY):latest

