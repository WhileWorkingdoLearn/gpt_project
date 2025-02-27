-- name: GetAllTasks :many
SELECT * FROM Task ORDER BY created_at DESC LIMIT
  $1
OFFSET
  $2;