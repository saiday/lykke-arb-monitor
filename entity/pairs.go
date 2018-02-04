package entity

// OrderUnit contains Volume and Price pair
type OrderUnit struct {
	Price  float64 `json:"Price"`
	Volume float64 `json:"Volume"`
}
