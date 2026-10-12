
# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS build

WORKDIR /app

# Copy Go module files
COPY container_src/go.mod container_src/go.sum ./
RUN go mod download

# Copy all Go source files and directories
COPY container_src/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /server ./cmd/tier-maker

FROM scratch

# Add CA certificates for HTTPS requests
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary
COPY --from=build /server /server

EXPOSE 8080

CMD ["/server"]
