// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: get_all_tasks.sql

package database

import (
	"context"
)

const getAllTasks = `-- name: GetAllTasks :many
SELECT id, created_at, updated_at, order_id, name, data, status FROM Task ORDER BY created_at DESC LIMIT
  $1
OFFSET
  $2
`

type GetAllTasksParams struct {
	Limit  int32
	Offset int32
}



func DefaultGetParam() GetAllTasksParams {
	param := GetAllTasksParams{}
	param.Offset = 0
	param.Limit = 1000
	return param
}

func (q *Queries) GetAllTasks(ctx context.Context, arg GetAllTasksParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getAllTasks, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.OrderID,
			&i.Name,
			&i.Data,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
