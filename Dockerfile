# Step 1: Use the official cosmtrek/air image
FROM cosmtrek/air:latest

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Step 4: Download dependencies
RUN go mod download

# Step 5: Copy the rest of the application source code to the container
COPY . .

EXPOSE 2468:2468

# Step 6: Start the app with air for hot-reloading
CMD ["air"]
