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
RUN go build -o clientd ./cmd/clientd

FROM alpine:3.14 as production
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/clientd .
COPY --from=builder /app/deploy/clientd/clientd.yaml /etc/clientd.conf
LABEL service=clientd
LABEL type=daemon
EXPOSE 8084
# Exec built binary
CMD ./clientd -config /etc/clientd.conf
