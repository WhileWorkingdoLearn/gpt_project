-- name: UpdateTaskName :exec
UPDATE Task 
SET 
updated_at = NOW(),
name = $2
WHERE id = $1;