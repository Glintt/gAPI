package sockets

import (
	"github.com/Glintt/gAPI/api/utils"

	//"encoding/json"
	"strconv"
	"time"
)

var RequestsCount = 0

func IncrementRequestCounter() {
	RequestsCount = RequestsCount + 1
}
func ResetCounter() {
	RequestsCount = 0
}

func PropagateToSockets() {
	for _, element := range SocketsConnected {
		element.Emit("logs", string(strconv.Itoa(RequestsCount)))
	}

	ResetCounter()
}

func StartRequestsCounterSender() {
	PreventCrash()

	var ticker = time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				PropagateToSockets()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func PreventCrash() {
	if r := recover(); r != nil {
		utils.LogMessage("recover request counter", utils.DebugLogType)
	}
}
