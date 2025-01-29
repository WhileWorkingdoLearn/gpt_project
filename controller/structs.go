package controller

type OrderImpl struct {
	OrderIdField string `json:"order_id" validate:"required,min=3"`
	NameField    string `json:"name" validate:"required,min=3"`
	StatusField  string `json:"status" validate:"required,oneof='Neu' 'In Bearbeitung' 'Abgeschlossen'"`
	DataField    string `json:"data" validate:"required,gt=0"` // Muss größer als 0 sein
}

func (o OrderImpl) OrderId() string {
	return o.OrderIdField
}

func (o OrderImpl) Name() string {
	return o.NameField
}

func (o OrderImpl) Status() string {
	return o.StatusField
}

func (o OrderImpl) Data() string {
	return o.DataField
}

type PostOrderInput struct {
	Orders []OrderImpl `json:"orders"`
}
