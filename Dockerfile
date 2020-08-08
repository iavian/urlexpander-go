FROM golang:1.13 as builder

COPY . /app

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.

RUN go mod download

# Build the binary.
RUN go build -mod=readonly -v -o /app/server

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine
RUN apk add --no-cache ca-certificates memcached

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

CMD ["ls  /app/server"]

CMD ["ls  /"]

CMD ["ls  /app"]

ENV MEMCACHED_SERVER localhost:11211

# Run the web service on container startup.
CMD ["/app/server"]
