ARG GO_VERSION=1.24.5

# Development stage
FROM golang:${GO_VERSION}-alpine AS development

RUN apk add --no-cache git curl && \
    go install github.com/air-verse/air@latest && \
    go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY .air.toml .air.toml
COPY . .

RUN make generate-lazy

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

# Production stage
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache git ca-certificates && \
    go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN templ generate && \
    sqlc generate && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/app/main.go

# Final production stage
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]
