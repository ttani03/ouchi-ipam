# Stage 1: Build the binary
FROM golang:1.20-alpine AS build

# Set the working directory to /app
WORKDIR /app

# Copy the source code into the container
COPY . .

# Install any needed packages
RUN go mod download

# Build the Go app
RUN go build -o ouchi-ipam .

# Stage 2: Create the final image
FROM alpine:latest

# Install any needed packages
RUN apk --no-cache add ca-certificates

# Set the working directory to /app
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/ouchi-ipam .

# Expose port 1323 for the server to listen on
EXPOSE 8080

# Start the server
CMD ["/app/ouchi-ipam"]
