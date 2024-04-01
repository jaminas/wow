# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o main ./cmd/server
FROM scratch
COPY --from=builder /build/main /
EXPOSE 8082
ENTRYPOINT ["/main"]
