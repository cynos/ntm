# Start from golang base latest image
FROM golang:latest

# Set maintainer
LABEL maintainer="Setia Budi"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy all resource into container
COPY . .

# Download and install dependencies
RUN go get -d -v ./...

# Install hot reload module for development purpose
RUN go install github.com/githubnemo/CompileDaemon@latest

# Command build and run app with hot reload mode
ENTRYPOINT CompileDaemon --build="go build -a -installsuffix cgo -o main ." --command="./main"