# Use an official Golang runtime as a parent image
FROM golang:1.22

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

RUN go install github.com/air-verse/air
RUN go build -o /server .

FROM scratch
COPY --from=build /server /server
EXPOSE 3000

# Define the command to run the app when the container starts
CMD ["air"]