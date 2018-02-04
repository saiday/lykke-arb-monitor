package entity

// OrderPair contains buy and sell order units
type OrderPair struct {
	Sell OrderUnit
	Buy  OrderUnit
}

// OrderUnit contains Volume and Price pair
type OrderUnit struct {
	Price  float64 `json:"Price"`
	Volume float64 `json:"Volume"`
}

// IsEmpty identiy OrderUnit struct is been set
func (u *OrderUnit) IsEmpty() bool {
	if u.Price == 0 && u.Volume == 0 {
		return true
	}

	return false
}

// NewOrderPair returns default OrderPair struct
func NewOrderPair() *OrderPair {
	return &OrderPair{*NewOrderUnit(), *NewOrderUnit()}
}

// NewOrderUnit returns default OrderUnit struct
func NewOrderUnit() *OrderUnit {
	return &OrderUnit{0, 0}
}
