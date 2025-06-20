FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev ca-certificates git

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

ARG SERVICE_NAME

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o /app/service \
    ./cmd/${SERVICE_NAME}

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/service /app/service
COPY --from=builder --chown=appuser:appgroup /app/internal/sqlc/migrations /app/internal/sqlc/migrations

USER appuser

EXPOSE 8080

CMD ["/app/service"]