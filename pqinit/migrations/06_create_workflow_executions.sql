CREATE TABLE workflow_executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL,
    node_id UUID,
    last_executed_node_id UUID,
    last_executed_task_id UUID,
    reference_number VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    message TEXT,
    last_node_executed_at TIMESTAMPTZ DEFAULT NOW(),
    last_node_completed_at TIMESTAMPTZ,
    last_task_executed_at TIMESTAMPTZ DEFAULT NOW(),
    last_task_completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (workflow_id) REFERENCES workflows(id),
    FOREIGN KEY (node_id) REFERENCES nodes(id)
);
