-- +goose Up
CREATE TABLE Task(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    order_id UUID NOT NULL,
    name Text NOT NULL,
    data Text NOT NULL,
    status Text NOT NULL CHECK(status IN ('New', 'In process','Completed'))
);

-- +goose Down
DROP TABLE Task;