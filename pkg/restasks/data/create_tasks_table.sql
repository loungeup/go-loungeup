CREATE TABLE IF NOT EXISTS tasks (
  id UUID PRIMARY KEY,
  error_message TEXT,
  result JSONB,
  progress INTEGER NOT NULL CHECK (
    progress >= 0
    AND progress <= 100
  ),
  started_at TIMESTAMP WITH TIME ZONE NOT NULL,
  ended_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS tasks_ended_at_idx ON tasks (ended_at);