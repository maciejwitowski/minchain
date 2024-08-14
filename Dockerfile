# Start from the official Go image
FROM golang:1.21

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod ./

# Download the dependencies
RUN go mod download

EXPOSE 9555

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the application
RUN go build -o main .

# Run the application
CMD ["./main"]
