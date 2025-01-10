#!/bin/bash
set -e

# Function to check if Docker containers are running
check_docker_containers() {
    if [ $(docker ps -q -f name=api -f name=db | wc -l) -eq 0 ]; then
        echo "Docker containers are not running. Starting containers..."
        docker-compose up -d
        sleep 10 # Wait for containers to start
    else
        echo "Docker containers are already running."
    fi
}

# Check Docker containers
check_docker_containers

API_URL="http://localhost:8080/workflow/node"

# Create nodes
start_node_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"title": "start", "type": "Start", "description": "Start node"}' | jq -r '.id')
input_new_user_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"title": "input_new_user", "type": "Task", "description": "Input new user"}' | jq -r '.id')
check_user_personal_info_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"title": "check_user_personal_info", "type": "Task", "description": "Check user personal info"}' | jq -r '.id')
save_user_data_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"title": "save_user_data", "type": "Task", "description": "Save user data"}' | jq -r '.id')
end_node_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"title": "end", "type": "End", "description": "End node"}' | jq -r '.id')

# Create relationships
curl -s -X POST $API_URL/$start_node_id/relationship -H "Content-Type: application/json" -d "{\"ancestor\": \"$start_node_id\", \"descendant\": \"$input_new_user_id\"}"
curl -s -X POST $API_URL/$input_new_user_id/relationship -H "Content-Type: application/json" -d "{\"ancestor\": \"$input_new_user_id\", \"descendant\": \"$check_user_personal_info_id\"}"
curl -s -X POST $API_URL/$check_user_personal_info_id/relationship -H "Content-Type: application/json" -d "{\"ancestor\": \"$check_user_personal_info_id\", \"descendant\": \"$save_user_data_id\"}"
curl -s -X POST $API_URL/$save_user_data_id/relationship -H "Content-Type: application/json" -d "{\"ancestor\": \"$save_user_data_id\", \"descendant\": \"$end_node_id\"}"

# Execute the workflow
curl -s -X POST $API_URL/$start_node_id/execute -H "Content-Type: application/json"

echo "Workflow created and executed successfully"
