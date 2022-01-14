FROM golang:1.17-buster as builder

COPY go.mod /go/src/github.com/ddvkid/learngo/
COPY go.sum /go/src/github.com/ddvkid/learngo/

COPY internal /go/src/github.com/ddvkid/learngo/internal/
COPY services/order/ /go/src/github.com/ddvkid/learngo/services/order/

WORKDIR /go/src/github.com/ddvkid/learngo/services/order

ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o order

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/github.com/ddvkid/learngo/services/order .

EXPOSE 8080
# Run the web service on container startup.
CMD ["/app/order"]
