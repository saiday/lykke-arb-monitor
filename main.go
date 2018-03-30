package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/everdev/mack"
	"github.com/saiday/lykke-arb-monitor/entity"
)

func main() {
	fetchOrderEvery(time.Second*15, monitorOrderBooks)
}

func fetchOrderEvery(d time.Duration, f func(string) *entity.OrderPair) {
	for range time.Tick(d) {
		lkk := f("LKKUSD")
		btclkk := f("BTCLKK")
		btc := f("BTCUSD")
		arbChances(lkk, btclkk, btc, "BTC")
		ethlkk := f("ETHLKK")
		eth := f("ETHUSD")
		arbChances(lkk, ethlkk, eth, "ETH")
	}
}

// arbChances caculate is there's chance to arb
// pair1 represents target1 order
// pair2 represetns target2 glove order
// pair12 represents  target2 / target1 order
func arbChances(pair1 *entity.OrderPair, pair12 *entity.OrderPair, pair2 *entity.OrderPair, entity string) {
	// arbPair1Detail(pair1, pair12, pair2)
	arbPair1(pair1, pair12, pair2, entity)
	arbPair2(pair1, pair12, pair2, entity)
}

// e.g. LKK buy: 0.195  sell: 0.207
//      ETH buy: 908    sell: 909
//   ETHLKK buy: 4360   sell: 4675
func arbPair1(pair1 *entity.OrderPair, pair12 *entity.OrderPair, pair2 *entity.OrderPair, entity string) {
	sellPair1ToPair2Rate := pair2.Sells[0].Price / pair1.Buys[0].Price
	if sellPair1ToPair2Rate < pair12.Buys[0].Price {
		message := fmt.Sprintf("LKK to %s tourment", entity)
		mack.Notify(message)
		mack.Say(message)
		fmt.Printf("!!!!!!!!!!!!!!!!!!!!! chance to arb: sold LKK to usd, buy %s, sold to LKK: %f >>>> %f\n", entity, sellPair1ToPair2Rate, pair12.Buys[0].Price)
	} else {
		fmt.Printf("no chance to arb LKK: %f > %f\n", sellPair1ToPair2Rate, pair12.Buys[0].Price)
	}
}

func arbPair2(pair1 *entity.OrderPair, pair12 *entity.OrderPair, pair2 *entity.OrderPair, entity string) {
	sellPair2ToPair1Rate := pair2.Buys[0].Price / pair1.Sells[0].Price
	if sellPair2ToPair1Rate > pair12.Sells[0].Price {
		message := fmt.Sprintf("%s to LKK tourment", entity)
		mack.Notify(message)
		mack.Say(message)
		fmt.Printf("!!!!!!!!!!!!!!!!!!!!! chance to arb: sold %s to usd, buy LKK, sold to %s : %f >>>> %f\n", entity, entity, sellPair2ToPair1Rate, pair12.Sells[0].Price)
	} else {
		fmt.Printf("no chance to arb %s: %f < %f\n", entity, sellPair2ToPair1Rate, pair12.Sells[0].Price)
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

	findArbMaximum(&pair1Sells, &pair2Buys, &pair12Buys)

	return 0.0
}

func findArbMaximum(sells *[]entity.OrderUnit, buys *[]entity.OrderUnit, arbMarket *[]entity.OrderUnit) {

	*sells = append((*sells)[:0], (*sells)[1:]...)

	// sells[0] = sells[len(sells)-1]
	// sells = sells[:len(sells)-1]

	fmt.Printf("\n\n\n\n remove")
	fmt.Printf("%v, len is: %d\n", (*sells)[0], len(*sells))
}
