# Start by specifying the base image with Go 1.19 or later
FROM golang:1.19-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Copy the .env file to the container (if you need it inside the container)
# COPY .env .env

# Set the working directory to the location of your main.go file
WORKDIR /app/cmd/server

# Build the Go app
RUN go build -o /app/main .

# Expose port 50051 to the outside world
EXPOSE 50051

# Command to run the executable
CMD ["/app/main"]
