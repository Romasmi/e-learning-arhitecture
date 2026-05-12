BUF_VERSION := v1.45.0

.PHONY: generate clean

generate:
	go run github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) dep update
	go run github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) generate

clean:
	rm -rf gen/go/*
	rm -rf gen/swagger/*
