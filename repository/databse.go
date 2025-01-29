package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var sql_commands commands

func init() {
	sql_commands = commands{
		create: `CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	order_id VARCHAR(12) NOT NULL,
	name VARCHAR(50) NOT NULL,
	status VARCHAR(50) NOT NULL,
	createt_at timestamp NOT NULL DEFAULT NOW(),
	updatet_at timestamp NOT NULL DEFAULT NOW()
	)`,
		update:      "UPDATE orders SET WHERE id= ?",
		insertOrder: "INSERT INTO orders(name,status,data) VALUES (?,?,?,?);",
		inserMOrder: "INSERT INTO orders(name,status,data) VALUES %s;",
		delete:      "DELETE FROM orders WHERE id_key = ?",
	}
}

type RepositoryHandler interface {
	GetItem(id string) (IOrderFromDB, error)
	FindItem(id Query) ([]IOrderFromDB, error)
	AddItem(data IOrder) (int64, error)
	AddItems(data []IOrder) ([]int64, error)
	UpdateItem(data IOrder) (bool, error)
	DeleteItem(id string) (bool, error)
}

type dbHandler struct {
	db *sql.DB
}

// AddItem implements DBHandler.
func (dh dbHandler) AddItem(data IOrder) (int64, error) {

	insertStatement := fmt.Sprintf(sql_commands.inserMOrder, "(?, ?, ?, ?)")
	result, err := dh.db.Exec(insertStatement, data.OrderId(), data.Name(), data.Status(), data.Data())
	if err != nil {
		return 0, fmt.Errorf("AddItem: %v", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("AddItem: %v", err)
	}
	return id, nil
}

// AddItems implements DBHandler.
func (dh dbHandler) AddItems(data []IOrder) ([]int64, error) {
	fmt.Println(data)
	/*
		addedIds := make([]int64, 0)
		placeholders := make([]string, 0)
		valueArgs := []interface{}{}
		for _, d := range data {
			placeholders = append(placeholders, "(?, ?, ?, ?)")
			valueArgs = append(valueArgs, d.Name())
			valueArgs = append(valueArgs, d.Status())
			valueArgs = append(valueArgs, d.Data())
		}

		insertStatement := fmt.Sprintf(sql_commands.inserMOrder, strings.Join(placeholders, ","))
		dh.db.Exec(insertStatement, valueArgs...)*/
	return []int64{1, 2}, nil
}

// DeleteItem implements DBHandler.
func (dh dbHandler) DeleteItem(id string) (bool, error) {
	_, err := dh.db.Exec(sql_commands.delete, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// FindItem implements DBHandler.
func (dh dbHandler) FindItem(query Query) ([]IOrderFromDB, error) {
	var orders = make([]IOrderFromDB, 0)

	rows, err := dh.db.Query("SELECT * FROM orders WHERE")
	if err != nil {
		return nil, fmt.Errorf("FindITem %q: %v", query, err)
	}
	defer rows.Close()

	for rows.Next() {
		var orderFromDB orderFromDB
		if err := rows.Scan(
			&orderFromDB.id,
			&orderFromDB.name,
			&orderFromDB.status,
			&orderFromDB.created_at,
			&orderFromDB.updated_at); err != nil {

			return nil, fmt.Errorf("FintItem: %q: %v", query, err)
		}

		orders = append(orders, orderFromDB)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindItem %q: %v", query, err)
	}
	return orders, nil
}

// GetItem implements DBHandler.
func (dh dbHandler) GetItem(id string) (IOrderFromDB, error) {
	var order orderFromDB

	row := dh.db.QueryRow("SELECT * FROM orders WHERE id= ?", id)
	if err := row.Scan(&order.id, &order.orderId, &order.name, &order.status, &order.created_at, &order.updated_at); err != nil {
		if err == sql.ErrNoRows {
			return &order, fmt.Errorf("GetIem: %d: no such order", id)

		}
		return &order, fmt.Errorf("GetIem: %d: %v", id, err)
	}
	return order, nil
}

// UpdateItem implements DBHandler.
func (dh dbHandler) UpdateItem(data IOrder) (bool, error) {
	dh.db.Exec("")
	return false, nil
}

func ConnectToDB(connectionstring string) (RepositoryHandler, error) {
	database, err := sql.Open("postgres", connectionstring)
	defer database.Close()

	if err != nil {
		return nil, err
	}

	if err = database.Ping(); err != nil {
		return nil, err
	}

	if err = createOrderTable(database); err != nil {
		return nil, err
	}
	dbh := dbHandler{
		db: database,
	}

	return dbh, nil
}

/*
	name       string
	status     int
	data       []byte
	created_at string
	updated_at string
*/

func createOrderTable(db *sql.DB) error {

	_, err := db.Exec(sql_commands.create)
	if err != nil {
		return err
	}
	return nil
}
