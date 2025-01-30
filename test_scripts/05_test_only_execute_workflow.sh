#!/bin/bash
set -e

API_URL="http://localhost:8080/workflow"
workflow_id="7f52c328-3325-47f5-8cd7-e5817f002e1d"

# Execute the workflow
curl -s -X POST $API_URL/$workflow_id/execute -H "Content-Type: application/json"
if [ $? -ne 0 ]; then
    echo "Error executing the workflow"
    exit 1
fi

echo "Workflow executed successfully"
