CREATE TABLE workflow_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID,
    node_id UUID,
    status VARCHAR(50) NOT NULL,
    message TEXT,
    executed_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    error_message TEXT,
    action_type VARCHAR(50),
    metadata TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (workflow_id) REFERENCES workflows(id),
    FOREIGN KEY (node_id) REFERENCES nodes(id)
);

CREATE TABLE node_task_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_log_id UUID,
    node_task_id UUID,
    status VARCHAR(50) NOT NULL,
    message TEXT,
    executed_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    error_message TEXT,
    action_type VARCHAR(50),
    metadata TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (workflow_log_id) REFERENCES workflow_logs(id),
    FOREIGN KEY (node_task_id) REFERENCES node_tasks(id)
);