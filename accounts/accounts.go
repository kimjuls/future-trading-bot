package accounts

import (
	"strconv"

	"github.com/thoas/go-funk"
)

const minimumOrderPrice = 5000

const (
	b = "bid"
	s = "ask"
)

type Accounts interface {
	order(b *Bot, c *coin, side string, volume, price float64) (bool, error)
	accounts() ([]map[string]interface{}, error)
}

func getBalances(accounts []map[string]interface{}) map[string]float64 {
	return funk.Reduce(accounts, func(balances map[string]float64, acc map[string]interface{}) map[string]float64 {
		balance, err := strconv.ParseFloat(acc["balance"].(string), 64)
		if err != nil {
			panic(err)
		}
		balances[acc["currency"].(string)] = balance
		return balances
	}, make(map[string]float64)).(map[string]float64)
}

func getAverageBuyPrice(accounts []map[string]interface{}, coin string) float64 {
	t := funk.Find(accounts, func(acc map[string]interface{}) bool { return acc["currency"].(string) == coin })

	if t, ok := t.(map[string]interface{}); ok {
		avgBuyPrice, err := strconv.ParseFloat(t["avg_buy_price"].(string), 64)
		if err != nil {
			panic(err)
		}
		return avgBuyPrice
	}
	return 0
}

func getTotalBalance(accounts []map[string]interface{}, balances map[string]float64) float64 {
	return funk.Reduce(funk.Keys(balances), func(totalBalance float64, coin string) float64 {
		totalBalance += getAverageBuyPrice(accounts, coin) * balances[coin]
		return totalBalance
	}, balances["KRW"]).(float64)
}
