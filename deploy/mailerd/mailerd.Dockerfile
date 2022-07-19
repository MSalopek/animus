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
RUN go build -o mailerd ./cmd/mailerd

FROM alpine:3.14 as production
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/mailerd .
COPY --from=builder /app/deploy/mailerd/mailerd.yaml /etc/mailerd.conf
LABEL service=mailerd
LABEL type=daemon
# Exec built binary
CMD ./mailerd -config /etc/mailerd.conf
