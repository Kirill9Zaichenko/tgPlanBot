ALTER TABLE tasks ADD COLUMN organization_id INTEGER NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_organization_id ON tasks(organization_id);