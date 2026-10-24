# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build

WORKDIR /src

ENV CGO_ENABLED=0 GOOS=linux

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S -g 10001 app \
    && adduser  -S -u 10001 -G app app

WORKDIR /app

COPY --from=build /out/server ./server
COPY openapi.yaml ./openapi.yaml
COPY data ./data

USER 10001:10001

EXPOSE 8080

ENTRYPOINT ["./server"]
