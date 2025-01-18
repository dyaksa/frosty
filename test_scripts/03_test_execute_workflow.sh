#!/bin/bash
set -e

API_URL="http://localhost:8080/workflow"

# Function to check if Docker containers are running
check_docker_containers() {
    if [ $(docker ps -q -f name=api -f name=db | wc -l) -eq 0 ]; then
        echo "Docker containers are not running. Start containers first!"
        exit 1
    else
        echo "Docker containers available."
    fi
}

# Check Docker containers
check_docker_containers

# Create a workflow
workflow_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"name": "Sample Workflow", "description": "A test workflow"}')
if [ $? -ne 0 ]; then
    echo "Error creating workflow"
    exit 1
fi

workflow_id=$(echo $workflow_id | sed 's/"//g')

echo workflow_id: $workflow_id

# Create nodes
node_titles=("start" "input_new_user" "check_user_personal_info" "save_user_data" "end")

for title in "${node_titles[@]}"; do
    node_id=$(curl -s -X POST $API_URL/node -H "Content-Type: application/json" -d "{\"workflow_id\": \"$workflow_id\", \"title\": \"$title\", \"type\": \"Task\", \"description\": \"$title node\"}")
    if [ $? -ne 0 ]; then
        echo "Error creating node $title"
        exit 1
    fi
    node_id=$(echo $node_id | sed 's/"//g')
    nodes[$title]=$node_id
done

# Create relationships
relationships=(
    "start input_new_user"
    "input_new_user check_user_personal_info"
    "check_user_personal_info save_user_data"
    "save_user_data end"
)

for relationship in "${relationships[@]}"; do
    ancestor=$(echo $relationship | cut -d' ' -f1)
    descendant=$(echo $relationship | cut -d' ' -f2)
    curl -s -X POST $API_URL/relationship -H "Content-Type: application/json" -d "{\"ancestor\": \"${nodes[$ancestor]}\", \"descendant\": \"${nodes[$descendant]}\"}"
    if [ $? -ne 0 ]; then
        echo "Error creating relationship between $ancestor and $descendant"
        exit 1
    fi
done

# Execute the workflow
curl -s -X POST $API_URL/$workflow_id/execute -H "Content-Type: application/json"
if [ $? -ne 0 ]; then
    echo "Error executing the workflow"
    exit 1
fi

echo "Workflow executed successfully"
