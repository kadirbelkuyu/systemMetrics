# Start from the latest golang base image
FROM golang:latest

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build server.go

# Expose port 3000 to the outside
EXPOSE 3000

# Command to run the executable
CMD ["./server"]