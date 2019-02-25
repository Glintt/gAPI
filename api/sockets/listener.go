package sockets

import (
	"gAPIManagement/api/config"
	"log"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/cors"
)

var SocketsConnected []socketio.Conn

func SocketListen() {
	port := os.Getenv("SOCKET_PORT")

	if port == "" {
		port = config.SOCKET_PORT_DEFAULT
	}

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(so socketio.Conn) error {
		SocketsConnected = append(SocketsConnected, so)
		return nil
	})

	server.OnDisconnect("/", func(so socketio.Conn, msg string) {
		var SocketsConnectedTemp []socketio.Conn
		for _, element := range SocketsConnected {
			if element.ID() != so.ID() {
				SocketsConnectedTemp = append(SocketsConnectedTemp, element)
			}
			SocketsConnected = SocketsConnectedTemp
		}
	})

	server.OnError("/", func(err error) {
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
