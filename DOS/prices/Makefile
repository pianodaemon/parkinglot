GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-w -s" 

PRICES_EXE=./fiscal

all: build

build: fmt
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD) -o $(PRICES_EXE) cmd/xxx/run.go

fmt:
	$(GOCMD) fmt ./...

clean:
	rm -f $(PRICES_EXE)
