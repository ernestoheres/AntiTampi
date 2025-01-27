package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		
		for {
			//read message
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			//print the above message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
			
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.ListenAndServe(":6974", nil)
}
