#!/bin/bash
set -e

BUILD_FLAG=false

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --build)
            BUILD_FLAG=true
            shift
            ;;
        *)
            shift
            ;;
    esac
done

if [ "$BUILD_FLAG" = true ]; then
    # Start Docker containers with build
    docker-compose up -d --build
else
    # Start Docker containers without build
    docker-compose up -d
fi

echo "Docker containers have been initialized and are running."
