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

RUN templ generate

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
