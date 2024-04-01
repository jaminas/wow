FROM golang:1.21 AS builder
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o main ./cmd/client
FROM scratch
COPY --from=builder /build/main /
ENTRYPOINT ["/main"]
