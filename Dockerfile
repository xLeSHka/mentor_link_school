FROM golang:1.24-alpine3.21 AS builder

# Setup base software for building an app.
RUN apk update && apk add ca-certificates git

WORKDIR /app

# Download dependencies.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy application source.
COPY . .

RUN chmod a+x ./deploy/swag

RUN ./deploy/swag init -g ./internal/transport/http/httpServer.go
# Build the application.
RUN go build -o bin/application ./cmd/back/main.go

# Prepare executor image.
FROM alpine:3.21 AS runner

RUN apk update && apk add ca-certificates bash && rm -rf /var/cache/apk/*

WORKDIR /app

COPY migrations migrations

COPY --from=builder /app/docs ./docs
COPY --from=builder /app/bin/application ./
COPY ./.env ./

# Run the application.
CMD ["./application"]