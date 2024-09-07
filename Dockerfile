# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application source code
COPY . .

# Download and install the Go dependencies
RUN go mod download

# Build the Go application
RUN go build -o go-auth

# Set the entry point for the container
CMD ["./go-auth"]