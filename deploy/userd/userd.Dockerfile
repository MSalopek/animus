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
RUN go build -o userd ./cmd/userd

FROM alpine:3.14 as production
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/userd .
COPY --from=builder /app/deploy/userd/userd.yaml /etc/userd.conf
LABEL service=userd
LABEL type=daemon
EXPOSE 8083
# Exec built binary
CMD ./userd -config /etc/userd.conf
