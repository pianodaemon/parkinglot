GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-w -s" 

PRICES_EXE=./sidecar-prices

all: build

build: fmt
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD) -o $(PRICES_EXE) cmd/http/run.go

test: build
	$(GOCMD) test -v ./...

fmt:
	$(GOCMD) fmt ./...

clean:
	rm -f $(PRICES_EXE)
