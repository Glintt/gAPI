package sockets

import (
	"fmt"
	"gAPIManagement/api/config"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var SocketsConnected []*gosocketio.Channel

func SocketListen() {
	port := os.Getenv("SOCKET_PORT")

	if port == "" {
		port = config.SOCKET_PORT_DEFAULT
	}

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		SocketsConnected = append(SocketsConnected, c)
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		var SocketsConnectedTemp []*gosocketio.Channel
		for _, element := range SocketsConnected {
			if element.Id() != c.Id() {
				SocketsConnectedTemp = append(SocketsConnectedTemp, element)
			}
			SocketsConnected = SocketsConnectedTemp
		}
	})
	//error catching handler
	server.On(gosocketio.OnError, func(c *gosocketio.Channel) {
		var SocketsConnectedTemp []*gosocketio.Channel
		for _, element := range SocketsConnected {
			if element.Id() != c.Id() {
				SocketsConnectedTemp = append(SocketsConnectedTemp, element)
			}
			SocketsConnected = SocketsConnectedTemp
		}
	})

	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins:   config.GApiConfiguration.Cors.AllowedOrigins,
		AllowCredentials: config.GApiConfiguration.Cors.AllowCredentials,
	})

	handler := c.Handler(mux)
	mux.Handle("/socket.io/", server)

	fmt.Println("WS PORT = " + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
