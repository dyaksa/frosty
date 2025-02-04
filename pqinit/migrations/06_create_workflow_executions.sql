CREATE TABLE workflow_executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL,
    node_id UUID,
    task_id UUID,
    status VARCHAR(50) NOT NULL,
    message TEXT,
    executed_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    FOREIGN KEY (workflow_id) REFERENCES workflows(id),
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);
