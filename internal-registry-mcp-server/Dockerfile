FROM golang:1.25.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -o /go/bin/mcp ./cmd/mcp

FROM alpine:latest AS api

WORKDIR /root/

COPY --from=builder /go/bin/app .

RUN apk add --no-cache tzdata

CMD ["./app"]

FROM alpine:latest AS mcp

WORKDIR /root/

COPY --from=builder /go/bin/mcp .

RUN apk add --no-cache tzdata

CMD ["./mcp"]
