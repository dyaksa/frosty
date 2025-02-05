#!/bin/bash
set -e

API_URL="http://localhost:8080/workflow"
workflow_id="5b8f46af-5cd2-4d01-a5da-a8fc02c60d85"

# Execute the workflow
curl -s -X POST $API_URL/$workflow_id/execute -H "Content-Type: application/json"
if [ $? -ne 0 ]; then
    echo "Error executing the workflow"
    exit 1
fi

echo "Workflow executed successfully"
