FROM alpine:latest as builder

RUN apk --no-cache add ca-certificates make musl-dev go

COPY . /app

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
RUN go mod download

# Build the binary.
RUN go build -v -o /app/server

# Run the web service on container startup.
CMD ["/app/server"]


