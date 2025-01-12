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

# Create nodes
start_node_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"title": "start", "type": "Start", "description": "Start node"}')
if [ $? -ne 0 ]; then
    echo "Error creating start node"
    exit 1
fi

start_node_id=$(echo $start_node_id | sed 's/"//g')

echo $start_node_id

input_new_user_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"title": "input_new_user", "type": "Task", "description": "Input new user"}')
if [ $? -ne 0 ]; then
    echo "Error creating input_new_user node"
    exit 1
fi

echo $input_new_user_id

input_new_user_id=$(echo $input_new_user_id | sed 's/"//g')

check_user_personal_info_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"title": "check_user_personal_info", "type": "Task", "description": "Check user personal info"}')
if [ $? -ne 0 ]; then
    echo "Error creating check_user_personal_info node"
    exit 1
fi

check_user_personal_info_id=$(echo $check_user_personal_info_id | sed 's/"//g')


save_user_data_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"title": "save_user_data", "type": "Task", "description": "Save user data"}')
if [ $? -ne 0 ]; then
    echo "Error creating save_user_data node"
    exit 1
fi

save_user_data_id=$(echo $save_user_data_id | sed 's/"//g')

end_node_id=$(curl -s -X POST $API_URL -H "Content-Type: application/json" -d '{"title": "end", "type": "End", "description": "End node"}')
if [ $? -ne 0 ]; then
    echo "Error creating end node"
    exit 1
fi

end_node_id=$(echo $end_node_id | sed 's/"//g')

# Create relationships
curl -s -X POST $API_URL/$start_node_id/relationship -H "Content-Type: application/json" -d "{\"ancestor\": \"$start_node_id\", \"descendant\": \"$input_new_user_id\"}"
if [ $? -ne 0 ]; then
    echo "Error creating relationship between start node and input_new_user node"
    exit 1
fi

curl -s -X POST $API_URL/$input_new_user_id/relationship -H "Content-Type: application/json" -d "{\"ancestor\": \"$input_new_user_id\", \"descendant\": \"$check_user_personal_info_id\"}"
if [ $? -ne 0 ]; then
    echo "Error creating relationship between input_new_user node and check_user_personal_info node"
    exit 1
fi

curl -s -X POST $API_URL/$check_user_personal_info_id/relationship -H "Content-Type: application/json" -d "{\"ancestor\": \"$check_user_personal_info_id\", \"descendant\": \"$save_user_data_id\"}"
if [ $? -ne 0 ]; then
    echo "Error creating relationship between check_user_personal_info node and save_user_data node"
    exit 1
fi

curl -s -X POST $API_URL/$save_user_data_id/relationship -H "Content-Type: application/json" -d "{\"ancestor\": \"$save_user_data_id\", \"descendant\": \"$end_node_id\"}"
if [ $? -ne 0 ]; then
    echo "Error creating relationship between save_user_data node and end node"
    exit 1
fi

# Execute the workflow
# curl -s -X POST $API_URL/$start_node_id/execute -H "Content-Type: application/json"
# if [ $? -ne 0 ]; then
#     echo "Error executing the workflow"
#     exit 1
# fi

echo "Node configuration created successfully"
