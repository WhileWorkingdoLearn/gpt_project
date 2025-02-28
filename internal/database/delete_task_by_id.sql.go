// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: delete_task_by_id.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const deleteTaskById = `-- name: DeleteTaskById :exec
DELETE FROM  Task WHERE id = $1
`

func (q *Queries) DeleteTaskById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteTaskById, id)
	return err
}
