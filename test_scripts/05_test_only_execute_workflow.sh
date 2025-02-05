#!/bin/bash
set -e

API_URL="http://localhost:8080/workflow"
workflow_id="ec40715c-5b80-4bdb-af4c-e1abc7d3a728"

# Execute the workflow
curl -s -X POST $API_URL/$workflow_id/execute -H "Content-Type: application/json"
if [ $? -ne 0 ]; then
    echo "Error executing the workflow"
    exit 1
fi

echo "Workflow executed successfully"
