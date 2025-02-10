# Use Golang base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy Go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Expose port 8080 for the Go API
EXPOSE 8080

CMD ["go", "get", "github.com/go-sql-driver/mysql"]

# Run the Go application
CMD ["go", "run", "./cmd/explore/main.go"]