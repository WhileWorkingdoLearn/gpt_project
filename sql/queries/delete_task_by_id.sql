-- name: DeleteTaskById :exec
DELETE FROM  Task WHERE id = $1;