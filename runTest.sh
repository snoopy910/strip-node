#!/bin/bash

# Start test containers
sudo docker compose -f docker-compose.test.yaml up -d

# Wait for database to be ready
echo "Waiting for database to be ready..."
until sudo docker compose -f docker-compose.test.yaml exec test-db pg_isready -U test_user -d test_db; do
    echo "Database is unavailable - sleeping"
    sleep 1
done
echo "Database is ready!"

# Run the tests
go test -v ./sequencer/

# Capture the exit code
exit_code=$?

# Always try to clean up
sudo docker compose down --remove-orphans

# Exit with the test exit code
exit $exit_code