package main

import (
	"github.com/kimjuls/future-trading-bot/bot"
	"github.com/kimjuls/future-trading-bot/strategy"
	"github.com/sirupsen/logrus"
)

func main() {
	b := bot.New([]strategy.Strategy{})
	acc, err := bot.NewAccounts(b)
	if err != nil {
		logrus.Fatal(err)
	}
	b.SetAccounts(acc)
	logrus.Panic(b.Run())
}
