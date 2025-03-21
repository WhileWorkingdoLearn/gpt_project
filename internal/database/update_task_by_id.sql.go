// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: update_task_by_id.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const updateTaskName = `-- name: UpdateTaskName :exec
UPDATE Task 
SET 
updated_at = NOW(),
name = $2
WHERE id = $1
`

type UpdateTaskNameParams struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) UpdateTaskName(ctx context.Context, arg UpdateTaskNameParams) error {
	_, err := q.db.ExecContext(ctx, updateTaskName, arg.ID, arg.Name)
	return err
}
