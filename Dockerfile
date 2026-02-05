FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /envdoc ./cmd/envdoc

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /envdoc /usr/local/bin/envdoc
ENTRYPOINT ["envdoc"]
