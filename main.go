package main

import (
	"encoding/json"
	"flag"
	"fmt"
	cm "github.com/reteps/go-coinmarketcap"
	"sort"
	"strings"
)

type SortedCoin struct {
	CoinName      string
	Difference    float64
	BestPrice     float64
	WorstPrice    float64
	BestExchange  string
	WorstExchange string
	BestPair      string
	WorstPair     string
	BestVolume    float64
	WorstVolume   float64
}

var pairs string
var exchanges string

func reverse(coins []SortedCoin) []SortedCoin {
	for i := 0; i < len(coins)/2; i++ {
		j := len(coins) - i - 1
		coins[i], coins[j] = coins[j], coins[i]
	}
	return coins
}
func main() {
	flag.StringVar(&exchanges, "exchanges", "", "Exchanges the coin can trade on")
	flag.StringVar(&pairs, "pairs", "", "Pairs the coin can trade into")
	var minimumVolume = flag.Float64("min_vol", 1.0, "Minimum Percent volume per exchange")
	var minimumRank = flag.Int("min_rank", 100, "Lowest rank a coin can have")
	var coinsShown = flag.Int("coins_shown", 10, "Number of results shown")
	flag.Parse()
	data, err := cm.GetAllCoinData(*minimumRank)
	if err != nil {
		panic(err)
	}
	var information []SortedCoin
	i := 1
	for coin_name, _ := range data {
		markets, err := cm.CoinMarkets(coin_name)
		if err != nil {
			panic(err)
		}
		best_price := 0.0
		worst_price := 99999999.0
		var best_exchange, worst_exchange, best_pair, worst_pair string
		var best_volume, worst_volume float64
		fits_criteria := false
		for _, market := range markets {
			trades_into := strings.Replace(strings.Replace(market.Pair, data[coin_name].Symbol, "", -1), "/", "", -1)
			if market.PercentVolume >= *minimumVolume && market.Updated && (strings.Contains(pairs, trades_into) || (pairs == "")) && (strings.Contains(exchanges, market.Exchange) || (exchanges == "")) {
				fits_criteria = true
				if market.Price >= best_price {
					best_price = market.Price
					best_exchange = market.Exchange
					best_pair = trades_into
					best_volume = market.PercentVolume
				} else if market.Price <= worst_price {
					worst_price = market.Price
					worst_exchange = market.Exchange
					worst_pair = trades_into
					worst_volume = market.PercentVolume
				}
			}
		}
		if fits_criteria {
			information = append(information, SortedCoin{CoinName: coin_name, Difference: float64(int64((best_price/worst_price-1)*10000+0.5) / 100), WorstPrice: worst_price, BestPrice: best_price, WorstExchange: worst_exchange, BestExchange: best_exchange, WorstPair: worst_pair, BestPair: best_pair, BestVolume: best_volume, WorstVolume: worst_volume})
		}
		fmt.Printf("\r[%d/%d]", i, *minimumRank)
		i += 1
	}
	sort.Slice(information, func(i, j int) bool { return information[i].Difference > information[j].Difference })
	for _, v := range reverse(information[:*coinsShown]) {
		pretty_v, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(pretty_v))
		//fmt.Printf("%s %f% %f %s %s %f %f %s %s %f", v.CoinName, v.Difference, v.BestPrice, v.BestExchange, v.BestPair, v.BestVolume, v.WorstPrice, v.WorstExchange, v.WorstPair, v.WorstVolume)
	}
}
