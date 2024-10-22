# Stage 1: Build
FROM golang:latest AS build
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main ./cmd/docker/main.go

# Stage 2: Runtime image
FROM alpine:3.18

# Install necessary packages
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary and assets
COPY --from=build /app/main .
COPY --from=build /app/assets ./assets
RUN chmod +x ./main

# Expose port 8080
EXPOSE 8080

# Start the application
CMD ["./main"]
