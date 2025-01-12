#!/bin/bash
set -e

DROP_DB_FLAG=false

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --drop-db)
            DROP_DB_FLAG=true
            shift
            ;;
        *)
            shift
            ;;
    esac
done

# Stop and remove Docker containers
docker-compose down

if [ "$DROP_DB_FLAG" = true ]; then
    # Remove the database data directory
    docker volume rm frosty_db-data
fi

echo "Docker containers have been deleted."
