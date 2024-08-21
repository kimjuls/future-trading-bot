package bot

import (
	"net/http"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"

	"github.com/kimjuls/future-trading-bot/accounts"
	"github.com/kimjuls/future-trading-bot/client"
	"github.com/kimjuls/future-trading-bot/static"
	"github.com/kimjuls/future-trading-bot/strategy"
)

type Bot struct {
	*client.Client
	*client.QuotationClient
	accounts   accounts.Accounts
	strategies []strategy.Strategy
}

func New(strategies []strategy.Strategy) *Bot {
	c := &client.Client{
		Client:    http.DefaultClient,
		AccessKey: static.Config.KeyPair.AccessKey,
		SecretKey: static.Config.KeyPair.SecretKey,
	}
	qc := &client.QuotationClient{Client: http.DefaultClient}

	return &Bot{c, qc, nil, strategies}
}

func (b *Bot) SetAccounts(accounts accounts.Accounts) {
	b.accounts = accounts
}

func (b *Bot) Run() error {

	log.Logger <- log.Log{Msg: "Bot started...", Level: logrus.WarnLevel}

	for _, strategy := range b.strategies {
		log.Logger <- log.Log{
			Msg:    "Register strategy...",
			Fields: logrus.Fields{"strategy": reflect.TypeOf(strategy).String()},
			Level:  logrus.DebugLevel,
		}
		if err := strategy.register(b); err != nil {
			return err
		}
	}

	if err := b.inHands(); err != nil {
		return err
	}

	d := newDetector()
	go d.run(b, predicate)

	for tick := range d.d {
		market := tick["code"].(string)

		if b.trackable(market) {
			log.Logger <- log.Log{
				Msg:    "Detected",
				Fields: logrus.Fields{"market": market},
				Level:  logrus.DebugLevel,
			}
			if err := b.launch(market); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bot) trackable(market string) bool {
	if static.Config.MaxTrackedMarket > 0 {
		if len(getMarketsFromStates(staged)) >= static.Config.MaxTrackedMarket {
			return false
		}
	}

	excluded := funk.Contains(static.Config.Blacklist, market)
	included := funk.Contains(static.Config.Whitelist, market)

	if excluded {
		return false
	}

	if len(static.Config.Whitelist) < 1 || included {
		if s, ok := states[market]; ok {
			return s == untracked
		} else {
			return true
		}
	}

	return false
}

func (b *Bot) inHands() error {
	acc, err := b.accounts.accounts()
	if err != nil {
		return err
	}

	balances := getBalances(acc)
	delete(balances, "KRW")

	for coin := range balances {
		if err := b.launch(targetMarket + "-" + coin); err != nil {
			return err
		}
	}
	log.Logger <- log.Log{Msg: "Run strategy for coins in hands.", Level: logrus.DebugLevel}

	return nil
}

func (b *Bot) launch(market string) error {
	coin, err := newCoin(b.accounts, market[4:], static.Config.TradableBalanceRatio)
	if err != nil {
		return err
	}

	states[market] = staged

	go b.tick(coin)

	for _, strategy := range b.strategies {
		if err := strategy.boot(b, coin); err != nil {
			return err
		}
		log.Logger <- log.Log{
			Msg:    "Booting strategy...",
			Fields: logrus.Fields{"strategy": reflect.TypeOf(strategy).String()},
			Level:  logrus.DebugLevel,
		}
		go b.strategy(coin, strategy)
	}

	return nil
}

func (b *Bot) strategy(c *coin, strategy Strategy) {
	defer func() {
		if err := recover(); err != nil {
			log.Logger <- log.Log{
				Msg: err,
				Fields: logrus.Fields{
					"role": "Strategy", "strategy": reflect.TypeOf(strategy).String(), "coin": c.name,
				},
				Level: logrus.ErrorLevel,
			}
		}
	}()

	log.Logger <- log.Log{
		Msg:    "STARTED",
		Fields: logrus.Fields{"strategy": reflect.TypeOf(strategy).String(), "coin": c.name},
		Level:  logrus.DebugLevel,
	}

	stat, ok := states[targetMarket+"-"+c.name]

	for ok && stat == staged {
		t := <-c.t

		acc, err := b.accounts.accounts()
		if err != nil {
			panic(err)
		}

		balances := getBalances(acc)

		if balances["KRW"] >= minimumOrderPrice && balances["KRW"] > c.onceOrderPrice && c.onceOrderPrice > minimumOrderPrice {
			if _, err := strategy.run(b, c, t); err != nil {
				panic(err)
			}
		}
	}

	log.Logger <- log.Log{
		Msg:    "CLOSED",
		Fields: logrus.Fields{"strategy": reflect.TypeOf(strategy).String(), "coin": c.name},
		Level:  logrus.DebugLevel,
	}
}

func (b *Bot) tick(c *coin) {
	defer func() {
		if err := recover(); err != nil {
			log.Logger <- log.Log{Msg: err, Fields: logrus.Fields{"role": "Tick", "coin": c.name}, Level: logrus.ErrorLevel}
		}
	}()

	m := targetMarket + "-" + c.name

	wsc, err := client.NewWebsocketClient("ticker", []string{m}, true, false)
	if err != nil {
		panic(err)
	}

	for states[m] == staged {
		var r map[string]interface{}

		if err := wsc.Ws.WriteJSON(wsc.Data); err != nil {
			panic(err)
		}

		if err := wsc.Ws.ReadJSON(&r); err != nil {
			panic(err)
		}

		log.Logger <- log.Log{
			Msg: c.name,
			Fields: logrus.Fields{
				"change-rate": r["signed_change_rate"].(float64),
				"price":       r["trade_price"].(float64),
			},
			Level: logrus.TraceLevel,
		}

		for range b.strategies {
			c.t <- r
		}

		time.Sleep(time.Second * 1)
	}
	if err := wsc.Ws.Close(); err != nil {
		panic(err)
	}

	log.Logger <- log.Log{
		Msg:    "CLOSED",
		Fields: logrus.Fields{"role": "Tick", "coin": c.name},
		Level:  logrus.DebugLevel,
	}
}
