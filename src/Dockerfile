# Build stage
FROM golang:1.21-alpine AS builder

# Required for go mod downloads depending on network setup
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod ./
RUN go mod tidy

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o netmon .

# Run stage
FROM alpine:3.19

# Add tzdata and ca-certificates
RUN apk --no-cache add ca-certificates tzdata libcap \
    && adduser -D -H -u 1000 netmon

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/netmon /app/netmon
RUN chown netmon:netmon /app/netmon

# Explicitly grant the required capability to the binary so we can execute it as non-root
RUN setcap 'cap_net_raw=+ep' /app/netmon

USER netmon

ENTRYPOINT ["/app/netmon"]
