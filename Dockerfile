# Use an official Golang runtime as a parent image
FROM golang:1.22

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

RUN go install github.com/air-verse/air
RUN go install module github.com/rousseau-romain/round-timing

RUN go mod vendor

# Download and install any required dependencies
RUN go mod download

# Build the Go app
RUN air

# Expose port 8080 for incoming traffic
EXPOSE 8080

# Define the command to run the app when the container starts
CMD ["air"]