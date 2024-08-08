package main

import (
	"fmt"
	"net/http"

	"github.com/kimjuls/future-trading-bot/binanceapi"
	"github.com/kimjuls/future-trading-bot/static"
)

func main() {
	fmt.Println("Started")
	c := &binanceapi.Client{
		Client:    http.DefaultClient,
		AccessKey: static.Config.KeyPair.AccessKey,
		SecretKey: static.Config.KeyPair.SecretKey,
	}
	params := struct {
		Symbol   string `url:"symbol"`
		Interval string `url:"interval"`
	}{"BTCUSDT", "1m"}
	_, err := c.Request("GET", "/fapi/v1/klines", params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Ended")
}
