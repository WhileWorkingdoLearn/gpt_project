-- name: CreateTask :one
INSERT INTO Task(id,created_at,updated_at,order_id,name,data,status) VALUES(
      gen_random_uuid(),
      NOW(),
      NOW(),
      $1,
      $2,
      $3,
      $4
) RETURNING *;