package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/saiday/lykke-arb-monitor/entity"
)

func main() {
	response, err := http.Get("https://hft-api.lykke.com/api/OrderBooks/LKKUSD")
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		defer response.Body.Close()
		entity := &entity.OrderBooksResponse{}
		json.NewDecoder(response.Body).Decode(entity)
		fmt.Printf("%f", (*entity)[0].Prices[0].Price)
	}
}
