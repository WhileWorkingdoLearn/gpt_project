package repository

type IOrder interface {
	OrderId() string
	Name() string
	Status() string
	Data() string
}

type IOrderFromDB interface {
	Id() int64
	IOrder
	Updated_at() string
	ToString() string
	IsValid() bool
}
