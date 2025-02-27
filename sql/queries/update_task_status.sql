-- name: UpdateTaskStatus :exec
UPDATE Task 
SET 
updated_at = NOW(),
status = $2
WHERE id = $1;