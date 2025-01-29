package repository

import "fmt"

type Query struct {
}

type commands struct {
	create      string
	insertOrder string
	inserMOrder string
	update      string
	delete      string
	selectTable string
}

type orderFromDB struct {
	id         int64
	orderId    string
	name       string
	status     string
	data       string
	created_at string
	updated_at string
}

func (o orderFromDB) Id() int64          { return o.id }
func (o orderFromDB) OrderId() string    { return o.orderId }
func (o orderFromDB) Name() string       { return o.name }
func (o orderFromDB) Status() string     { return o.status }
func (o orderFromDB) Data() string       { return o.data }
func (o orderFromDB) Created_at() string { return o.created_at }
func (o orderFromDB) Updated_at() string { return o.updated_at }
func (o orderFromDB) ToString() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v\n", o.id, o.orderId, o.name, o.status, o.created_at, o.updated_at)
}
func (o orderFromDB) IsValid() bool {
	if o.id < 0 {
		return false
	}
	if len(o.orderId) < 6 {
		return false
	}
	if len(o.name) == 0 {
		return false
	}
	if len(o.status) < 1 {
		return false
	}
	return true
}

func NewOrder(
	id,
	name string,
	status string,
	data string,
	created_at,
	updated_at string) IOrderFromDB {

	return orderFromDB{name: name, status: status, data: data}
}
