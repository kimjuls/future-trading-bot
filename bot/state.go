package bot

import "github.com/thoas/go-funk"

const (
	staged = iota
	untracked
)

var states = make(map[string]int)

func getMarketsFromStates(marketState int) []string {
	return funk.Filter(funk.Keys(states), func(market string) bool {
		return states[market] == marketState
	}).([]string)
}
