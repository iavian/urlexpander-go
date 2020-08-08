FROM alpine:latest as builder

RUN apk add --no-cache ca-certificates memcached git make musl-dev go

COPY . /app

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.

RUN go mod download

# Build the binary.
RUN go build -mod=readonly -v -o /app/server

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

ENV MEMCACHED_SERVER localhost:11211

# Run the web service on container startup.
CMD ["/app/server"]
