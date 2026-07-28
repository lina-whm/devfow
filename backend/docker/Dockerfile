FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/devflow ./cmd/devflow

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -u 1001 app
COPY --from=builder /app/devflow /usr/local/bin/devflow
USER app
EXPOSE 8080
CMD ["devflow"]
