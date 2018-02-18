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
	arbPair1Detail(pair1, pair12, pair2)
	arbPair1(pair1, pair12, pair2)
	arbPair2(pair1, pair12, pair2)
}

// e.g. LKK buy: 0.195  sell: 0.207
//      ETH buy: 908    sell: 909
//   ETHLKK buy: 4360   sell: 4675
func arbPair1(pair1 *entity.OrderPair, pair12 *entity.OrderPair, pair2 *entity.OrderPair) {
	sellPair1ToPair2Rate := pair2.Sells[0].Price / pair1.Buys[0].Price
	if sellPair1ToPair2Rate < pair12.Buys[0].Price {
		fmt.Printf("!!!!!!!!!!!!!!!!!!!!! chance to arb: sold LKK to usd, buy ETH, sold to LKK: %f >>>> %f\n", sellPair1ToPair2Rate, pair12.Buys[0].Price)
	} else {
		fmt.Printf("no chance to arb LKK: %f > %f\n", sellPair1ToPair2Rate, pair12.Buys[0].Price)
	}
}

func arbPair2(pair1 *entity.OrderPair, pair12 *entity.OrderPair, pair2 *entity.OrderPair) {
	sellPair2ToPair1Rate := pair2.Buys[0].Price / pair1.Sells[0].Price
	if sellPair2ToPair1Rate > pair12.Sells[0].Price {
		fmt.Printf("!!!!!!!!!!!!!!!!!!!!! chance to arb: sold ETH to usd, buy LKK, sold to ETH : %f >>>> %f\n", sellPair2ToPair1Rate, pair12.Sells[0].Price)
	} else {
		fmt.Printf("no chance to arb ETH: %f < %f\n", sellPair2ToPair1Rate, pair12.Sells[0].Price)
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
		buyOrder := orderBooksResponse.BuyUnit().Prices
		sellOrder := orderBooksResponse.SellUnit().Prices

		pair := &entity.OrderPair{Sells: sellOrder, Buys: buyOrder}
		return pair
	}

	return entity.NewOrderPair()
}

func arbPair1Detail(pair1 *entity.OrderPair, pair12 *entity.OrderPair, pair2 *entity.OrderPair) float64 {
	// pair1 as currency base unit
	pair1Sells := make([]entity.OrderUnit, len(pair1.Sells))
	copy(pair1Sells, pair1.Sells)

	pair2Buys := make([]entity.OrderUnit, len(pair2.Buys))
	copy(pair2Buys, pair2.Buys)

	pair12Buys := make([]entity.OrderUnit, len(pair12.Buys))
	copy(pair12Buys, pair12.Buys)

	fmt.Printf("\n\n\n\n origin")
	fmt.Printf("%v, len is: %d\n", pair1Sells[0], len(pair1Sells))

	findArbMaximum(&pair1Sells, &pair2Buys, &pair12Buys)

	fmt.Printf("\n\n\n\n final")
	fmt.Printf("%v, len is: %d\n", pair1Sells[0], len(pair1Sells))

	return 0.0
}

func findArbMaximum(sells *[]entity.OrderUnit, buys *[]entity.OrderUnit, arbMarket *[]entity.OrderUnit) {

	*sells = append((*sells)[:0], (*sells)[1:]...)

	// sells[0] = sells[len(sells)-1]
	// sells = sells[:len(sells)-1]

	fmt.Printf("\n\n\n\n remove")
	fmt.Printf("%v, len is: %d\n", (*sells)[0], len(*sells))
}
