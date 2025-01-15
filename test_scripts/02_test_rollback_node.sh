#!/bin/bash
set -e

# Function to check if Docker containers are running
check_docker_containers() {
    if [ $(docker ps -q -f name=api -f name=db | wc -l) -eq 0 ]; then
        echo "Docker containers are not running. Start containers first!"
    else
        echo "Docker containers available."
    fi
}

# Check Docker containers
check_docker_containers

API_URL="http://localhost:8080/workflow/node"

# Rollback the workflow
# Assuming the start node ID is passed as an argument to the script
if [ -z "$1" ]; then
    echo "Usage: $0 <start_node_id>"
    exit 1
fi

start_node_id=$1

# Rollback the workflow
curl -s -X POST "$API_URL/$start_node_id/rollback" -H "Content-Type: application/json"
if [ $? -ne 0 ]; then
    echo "Error rolling back the workflow"
    exit 1
fi

echo "Workflow rollback completed successfully"
