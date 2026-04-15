PACK            := oneuptime
PROVIDER        := pulumi-resource-$(PACK)
VERSION         ?= 0.0.1-dev
PROVIDER_PATH   := provider/cmd/pulumi-resource-oneuptime
LDFLAGS         := -X github.com/ottercoders/pulumi-oneuptime/provider.Version=$(VERSION)

.PHONY: provider install schema test lint clean

provider:
	go build -o bin/$(PROVIDER) -ldflags "$(LDFLAGS)" ./$(PROVIDER_PATH)

install: provider
	cp bin/$(PROVIDER) $(GOPATH)/bin/

schema: provider
	./bin/$(PROVIDER) schema > schema.json

test:
	go test ./provider/... -v -count=1

test-acceptance:
	go test -tags=acceptance ./tests/... -v -count=1

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ schema.json
