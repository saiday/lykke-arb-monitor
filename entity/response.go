package entity

// OrderBooksResponse is struct of /api/OrderBooks/{assetPairId}
type OrderBooksResponse []OrderBooksData

// OrderBooksData contains actual response data
type OrderBooksData struct {
	AssetPair string      `json:"AssetPair"`
	IsBuy     bool        `json:"IsBuy"`
	Timestamp string      `json:"Timestamp"`
	Prices    []OrderUnit `json:"Prices"`
}

// SellUnit retuns sell data
func (r *OrderBooksResponse) SellUnit() *OrderBooksData {
	for _, data := range *r {
		if !data.IsBuy {
			return &data
		}
	}

	return nil
}

// BuyUnit returns buy data
func (r *OrderBooksResponse) BuyUnit() *OrderBooksData {
	for _, data := range *r {
		if data.IsBuy {
			return &data
		}
	}
	return nil
}

// IsBuyDisplayMessage returns BUY if IsBuy is true, otherwise SELL
func (r *OrderBooksData) IsBuyDisplayMessage() string {
	if r.IsBuy {
		return "BUY"
	}
	return "SELL"
}
