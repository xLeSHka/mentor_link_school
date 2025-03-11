# Use the official Golang Alpine image as the builder stage
FROM golang:1.24.1-alpine3.21 AS builder

# Install necessary dependencies
RUN apk update && apk add --no-cache ca-certificates git

# Set the working directory inside the container
WORKDIR /app

# Copy and download Go module dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the application source code
COPY . .

# Make the swag executable and initialize Swagger documentation
RUN chmod a+x ./deploy/swag && ./deploy/swag init -g ./internal/transport/http/httpServer.go

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/application ./cmd/back/main.go

# Use a minimal Alpine image for the final stage
FROM alpine:3.21 AS runner

# Install necessary runtime dependencies
RUN apk update && apk add --no-cache ca-certificates bash && rm -rf /var/cache/apk/*

# Set the working directory inside the container
WORKDIR /app

# Copy necessary files from the builder stage
COPY --from=builder /app/bin/application ./
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations
COPY ./.env ./

# Expose the port the application will run on (if applicable)
EXPOSE 8080

# Set the command to run the application
CMD ["./application"]