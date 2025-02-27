-- name: GetTaskByMonth :many
SELECT * FROM Task WHERE created_at >= date_trunc('month', $1::TIMESTAMP)
      AND created_at <  date_trunc('month', $1::TIMESTAMP) + INTERVAL '1 month' ORDER BY created_at DESC;