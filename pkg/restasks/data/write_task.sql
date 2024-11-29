INSERT INTO
  tasks (
    id,
    error_message,
    result,
    progress,
    started_at,
    ended_at
  )
VALUES
  ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO
UPDATE
SET
  error_message = EXCLUDED.error_message,
  result = EXCLUDED.result,
  progress = EXCLUDED.progress,
  ended_at = EXCLUDED.ended_at;