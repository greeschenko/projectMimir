# ---------- Runtime Image ----------
FROM alpine:latest

# Set working directory inside container
WORKDIR /app

# Add CA certificates for HTTPS support
RUN apk add --no-cache ca-certificates

# Copy prebuilt Go binary from local build
COPY ./auth-service/app ./app

# Copy database migrations into container
COPY ./auth-service/migrations ./migrations

# Expose gRPC port
EXPOSE 50051

# Run the binary
CMD ["./app"]
