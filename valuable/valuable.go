package valuable

import "github.com/stock-simulator-server/utils"

var Valuables = make(map[string]Valuable)
var ValuablesLock = utils.NewLock("valuables")
var ValuableUpdateChannel = utils.MakeDuplicator()

type Valuable interface {
	GetID() string
	GetValue() float64
	GetLock() *utils.Lock
	GetUpdateChannel() *utils.ChannelDuplicator
}
