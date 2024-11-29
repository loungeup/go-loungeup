SELECT
  id,
  error_message,
  result,
  progress,
  started_at,
  ended_at
FROM
  tasks
WHERE
  id = $1;