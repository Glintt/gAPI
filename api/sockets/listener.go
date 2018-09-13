package sockets

import (
	"gAPIManagement/api/config"
	"log"
	"net/http"
	"os"

	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"
)

var SocketsConnected []socketio.Socket

func SocketListen() {
	port := os.Getenv("SOCKET_PORT")

	if port == "" {
		port = config.SOCKET_PORT_DEFAULT
	}

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		SocketsConnected = append(SocketsConnected, so)

		so.On("disconnection", func() {
			var SocketsConnectedTemp []socketio.Socket
			for _, element := range SocketsConnected {
				if element.Id() != so.Id() {
					SocketsConnectedTemp = append(SocketsConnectedTemp, element)
				}
				SocketsConnected = SocketsConnectedTemp
			}
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins:   config.GApiConfiguration.Cors.AllowedOrigins,
		AllowCredentials: config.GApiConfiguration.Cors.AllowCredentials,
	})

	handler := c.Handler(mux)

	mux.Handle("/socket.io/", server)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
