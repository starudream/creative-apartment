PROJECT ?= $(shell basename $(CURDIR))
MODULE  ?= $(shell go list -m)

GO      ?= GO111MODULE=on go
VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")
BIDTIME ?= $(shell date +%FT%T%z)

BITTAGS := viper_yaml3
LDFLAGS := -s -w
LDFLAGS += -X "$(MODULE)/config.VERSION=$(VERSION)"
LDFLAGS += -X "$(MODULE)/config.BIDTIME=$(BIDTIME)"

.PHONY: bin

bin:
	@$(MAKE) tidy
	CGO_ENABLED=0 $(GO) build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/app $(MODULE)

run:
	@$(MAKE) tidy
	CGO_ENABLED=1 $(GO) run -race -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' $(MODULE)

test:
	@$(MAKE) tidy
	CGO_ENABLED=1 $(GO) test -race -tags '$(BITTAGS)' -count=1 -cover -v $(MODULE)/internal/...

tidy:
	$(GO) mod tidy

upx:
	upx bin/app

lint:
	golangci-lint run --skip-dirs-use-default
