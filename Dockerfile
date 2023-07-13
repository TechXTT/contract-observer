# Use an official Golang runtime as a parent image
FROM golang:1.20-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any dependencies required by your application
RUN apk add --no-cache git

# Build the Go application
RUN go build -o cmd/main cmd/main.go

# Run the application when the container starts
CMD ["./main"]