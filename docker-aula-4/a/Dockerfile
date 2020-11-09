# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the golang base image
FROM golang:1.15.4-alpine3.12 as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the Go app
RUN go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o a .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/a .

# Build a small image
FROM scratch AS final
COPY --from=builder /build/a /
COPY --from=builder /build/templates ./templates

# Expose port 9090 to the outside world
EXPOSE 9090

# Run application
ENTRYPOINT ["./a"]
