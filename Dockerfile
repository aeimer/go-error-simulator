FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-error-simulator

# Use a minimal base image for the final build
FROM scratch
COPY --from=builder /app/go-error-simulator /go-error-simulator

ENTRYPOINT ["/go-error-simulator"]

# Expose the port the app runs on by default
EXPOSE 8080

# Set the environment variable to specify the port
ENV PORT=8080

LABEL org.opencontainers.image.source=https://github.com/aeimer/go-error-simulator
LABEL org.opencontainers.image.licenses=MIT
