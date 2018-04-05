package sockets


import (
	"time"
	"strconv"
)

var RequestsCount = 0

func PropagateToSockets(){
	for _, element := range SocketsConnected {
		element.Emit("logs", strconv.Itoa(RequestsCount))
	}
	RequestsCount = 0
}

func StartRequestsCounterSender(){
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
		select {
			case <- ticker.C:
				PropagateToSockets()
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
}