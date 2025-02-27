-- name: GetTaskByID :one
SELECT * FROM Task WHERE id = $1;