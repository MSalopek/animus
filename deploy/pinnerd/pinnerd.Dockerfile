FROM golang:1.18 as builder
# Define build env
ENV GOOS linux
ENV CGO_ENABLED 0
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
# Build app
RUN go build -o pinnerd ./cmd/pinnerd

FROM alpine:3.14 as production
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/pinnerd .
COPY --from=builder /app/deploy/pinnerd/pinnerd.yaml /etc/pinnerd.conf
LABEL service=pinnerd
LABEL type=daemon
# Exec built binary
CMD ./pinnerd -config /etc/pinnerd.conf
