FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -o /envdoc ./cmd/envdoc

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /envdoc /usr/local/bin/envdoc
ENTRYPOINT ["envdoc"]
