package entity

// OrderBooksResponse is struct of /api/OrderBooks/{assetPairId}
type OrderBooksResponse []struct {
	AssetPair string      `json:"AssetPair"`
	IsBuy     bool        `json:"IsBuy"`
	Timestamp string      `json:"Timestamp"`
	Prices    []OrderUnit `json:"Prices"`
}
