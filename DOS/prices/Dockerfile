# Building stage
FROM golang:1.23.2-alpine AS builder

WORKDIR /go/src/blaucorp.com/prices

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY pkg pkg

COPY internal internal

COPY cmd cmd

RUN CGO_ENABLED=0 \
    GOOS=linux    \
    GOARCH=amd64  \
    go build -ldflags="-w -s" -o /sidecar-prices cmd/http/run.go

# Final image
FROM debian:bookworm

ENV APP_DIR=/
COPY --from=builder /sidecar-prices $APP_DIR
COPY scripts/run.sh $APP_DIR
WORKDIR $APP_DIR

EXPOSE 10100

CMD ["/run.sh"]
