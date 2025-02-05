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

# Create nodes
node_titles=("start" "input_new_user" "check_user_personal_info" "save_user_data" "end")

# Declare associative array in a way that works in both bash and zsh
if [ -n "$ZSH_VERSION" ]; then
    typeset -A nodes
elif [ -n "$BASH_VERSION" ]; then
    declare -A nodes
else
    echo "Unsupported shell"
    exit 1
fi

for title in "${node_titles[@]}"; do
    if [ "$title" == "start" ]; then
        node_id=$(curl -s -X POST $API_URL/node -H "Content-Type: application/json" -d "{\"title\": \"$title\", \"type\": \"Start\", \"description\": \"$title node\"}")
    elif [ "$title" == "end" ]; then
        node_id=$(curl -s -X POST $API_URL/node -H "Content-Type: application/json" -d "{\"title\": \"$title\", \"type\": \"End\", \"description\": \"$title node\"}")
    else
        node_id=$(curl -s -X POST $API_URL/node -H "Content-Type: application/json" -d "{\"title\": \"$title\", \"type\": \"Task\", \"description\": \"$title node\"}")
    fi

    if [ $? -ne 0 ]; then
        echo "Error creating node $title"
        exit 1
    fi
    node_id=$(echo $node_id | sed 's/"//g')
    nodes[$title]=$node_id
done

# Create tasks for nodes
task_titles=("input_new_user_task" "check_user_personal_info_task" "save_user_data_task")
# task_actions=("http://example.com/api/input_new_user" "http://example.com/api/check_user_personal_info" "http://example.com/api/save_user_data")
task_actions=("https://webhook-test.com/4b56515b4c61c9f98dc48c494fb6c334" "https://webhook-test.com/19afedc91e785e45c4104dceedf36e4e" "https://webhook-test.com/832e7916239a45f3f56574d6d1f5389a")

for i in "${!task_titles[@]}"; do
    task_id=$(curl -s -X POST $API_URL/task -H "Content-Type: application/json" -d "{\"title\": \"${task_titles[$i]}\", \"type\": \"API\", \"http_method\": \"POST\", \"action\": \"${task_actions[$i]}\", \"params\": \"{}\", \"max_retries\": 3}")
    if [ $? -ne 0 ]; then
        echo "Error creating task ${task_titles[$i]}"
        exit 1
    fi

    task_id=$(echo $task_id | sed 's/"//g')
    curl -s -X POST $API_URL/node/task -H "Content-Type: application/json" -d "{\"task_id\": \"$task_id\", \"node_id\": \"${nodes[${node_titles[$i+1]}]}\", \"task_order\": $i}"
    if [ $? -ne 0 ]; then
        echo "Error adding task ${task_titles[$i]} to node ${node_titles[$i+1]}"
        exit 1
    fi
done

# Create workflow
start_node_id=${nodes["start"]}
echo start_node_id=$start_node_id
workflow_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d "{\"name\": \"Sample Workflow\", \"description\": \"A test workflow\", \"starting_node_id\": \"$start_node_id\"}")
if [ $? -ne 0 ]; then
    echo "Error creating workflow"
    exit 1
fi

workflow_id=$(echo $workflow_id | sed 's/"//g')

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
    curl -s -X POST $API_URL/node/${nodes[$ancestor]}/relationship -H "Content-Type: application/json" -d "{\"ancestor\": \"${nodes[$ancestor]}\", \"descendant\": \"${nodes[$descendant]}\"}"
    if [ $? -ne 0 ]; then
        echo "Error creating relationship between $ancestor and $descendant"
        exit 1
    fi
done

# Execute the workflow
echo $workflow_id

echo "Workflow setup successfully"
