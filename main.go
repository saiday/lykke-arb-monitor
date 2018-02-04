package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/saiday/lykke-arb-monitor/entity"
)

func main() {
	fetchOrderEvery(time.Second, monitorOrderBooks)
}

func fetchOrderEvery(d time.Duration, f func(string) *entity.OrderPair) {
	for range time.Tick(d) {
		lkk := f("LKKUSD")
		ethlkk := f("ETHLKK")
		eth := f("ETHUSD")
		arbChances(lkk, ethlkk, eth)
	}
}

// arbChances caculate is there's chance to arb
// pair1 represents target1 order
// pair2 represetns target2 glove order
// pair12 represents  target2 / target1 order
func arbChances(pair1 *entity.OrderPair, pair12 *entity.OrderPair, pair2 *entity.OrderPair) {
	arbPair1(pair1, pair12, pair2)
	arbPair2(pair1, pair12, pair2)
}

// e.g. LKK buy: 0.195  sell: 0.207
//      ETH buy: 908    sell: 909
//   ETHLKK buy: 4360   sell: 4675
func arbPair1(pair1 *entity.OrderPair, pair12 *entity.OrderPair, pair2 *entity.OrderPair) {
	sellPair1ToPair2Rate := pair2.Sell.Price / pair1.Buy.Price
	if sellPair1ToPair2Rate < pair12.Buy.Price {
		fmt.Printf("!!!!!!!!!!!!!!!!!!!!! chance to arb: sold LKK to usd, buy ETH, sold to LKK: %f >>>> %f\n", sellPair1ToPair2Rate, pair12.Buy.Price)
	} else {
		fmt.Printf("no chance to arb LKK: %f > %f\n", sellPair1ToPair2Rate, pair12.Buy.Price)
	}
}

func arbPair2(pair1 *entity.OrderPair, pair12 *entity.OrderPair, pair2 *entity.OrderPair) {
	sellPair2ToPair1Rate := pair2.Buy.Price / pair1.Sell.Price
	if sellPair2ToPair1Rate > pair12.Sell.Price {
		fmt.Printf("!!!!!!!!!!!!!!!!!!!!! chance to arb: sold ETH to usd, buy LKK, sold to ETH : %f >>>> %f\n", sellPair2ToPair1Rate, pair12.Sell.Price)
	} else {
		fmt.Printf("no chance to arb ETH: %f < %f\n", sellPair2ToPair1Rate, pair12.Sell.Price)
	}
}

func monitorOrderBooks(pairID string) *entity.OrderPair {
	endpoint := "https://hft-api.lykke.com/api/OrderBooks/" + pairID
	response, err := http.Get(endpoint)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		defer response.Body.Close()
		orderBooksResponse := &entity.OrderBooksResponse{}
		json.NewDecoder(response.Body).Decode(orderBooksResponse)
		buyOrder := orderBooksResponse.BuyUnit().Prices[0]
		sellOrder := orderBooksResponse.SellUnit().Prices[0]

		pair := &entity.OrderPair{sellOrder, buyOrder}
		return pair
	}

	return entity.NewOrderPair()
}
