PACK            := oneuptime
PROVIDER        := pulumi-resource-$(PACK)
VERSION         ?= 0.0.1-dev
PROVIDER_PATH   := provider/cmd/pulumi-resource-oneuptime
LDFLAGS         := -X github.com/ottercoders/pulumi-oneuptime/provider.Version=$(VERSION)

.PHONY: provider install schema sdk test lint clean

provider:
	go build -o bin/$(PROVIDER) -ldflags "$(LDFLAGS)" ./$(PROVIDER_PATH)

install: provider
	cp bin/$(PROVIDER) $(GOPATH)/bin/

schema: provider
	pulumi package get-schema ./bin/$(PROVIDER) > schema.json

sdk: schema
	pulumi package gen-sdk ./bin/$(PROVIDER) --language go --out sdk/go
	pulumi package gen-sdk ./bin/$(PROVIDER) --language nodejs --out sdk/nodejs
	pulumi package gen-sdk ./bin/$(PROVIDER) --language python --out sdk/python
	pulumi package gen-sdk ./bin/$(PROVIDER) --language dotnet --out sdk/dotnet

test:
	go test ./provider/... -v -count=1

test-acceptance:
	go test -tags=acceptance ./tests/... -v -count=1

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ schema.json sdk/
