GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-w -s" 

PARKINGLOT_EXE=parkinglot

all: build

build: format
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GOBUILD) -o $(PARKINGLOT_EXE) cmd/http/run.go

format:
	$(GOCMD) fmt ./...

clean:
	rm -f $(PARKINGLOT_EXE)
