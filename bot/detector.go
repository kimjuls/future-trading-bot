package bot

import (
	"time"

	"github.com/kimjuls/future-trading-bot/client"
	"github.com/sirupsen/logrus"
)

const targetMarket = "KRW"

type detector struct {
	d chan map[string]interface{}
}

func newDetector() *detector {
	return &detector{d: make(chan map[string]interface{})}
}

func (d *detector) run(bot *Bot, predicate func(b *Bot, t map[string]interface{}) bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Logger <- log.Log{Msg: err, Fields: logrus.Fields{"role": "Detector"}, Level: logrus.ErrorLevel}
		}
	}()

	markets, err := bot.QuotationClient.Call("/market/all", struct {
		IsDetail bool `url:"isDetail"`
	}{false})
	if err != nil {
		panic(err)
	}

	targetMarkets := getMarketNames(markets.([]map[string]interface{}), targetMarket)

	wsc, err := client.NewWebsocketClient("ticker", targetMarkets, true, false)
	if err != nil {
		panic(err)
	}

	defer wsc.Ws.Close()

	log.Logger <- log.Log{Msg: "Detector started...", Level: logrus.DebugLevel}

	for {
		if err := wsc.Ws.WriteJSON(wsc.Data); err != nil {
			panic(err)
		}

		for range targetMarkets {
			var r map[string]interface{}

			if err := wsc.Ws.ReadJSON(&r); err != nil {
				panic(err)
			}

			if predicate(bot, r) {
				d.d <- r
			}

			time.Sleep(time.Millisecond * 300)
		}

	}
}
