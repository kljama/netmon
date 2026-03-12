#!/bin/sh

# Check if certificates exist, if not generate them
if [ ! -f /etc/nginx/certs/cert.pem ] || [ ! -f /etc/nginx/certs/key.pem ]; then
    echo "Certificates not found in /etc/nginx/certs. Generating self-signed certificates..."
    mkdir -p /etc/nginx/certs
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout /etc/nginx/certs/key.pem \
        -out /etc/nginx/certs/cert.pem \
        -subj "/CN=localhost"
    echo "Certificates generated successfully."
fi

# Execute the main command
exec "$@"
