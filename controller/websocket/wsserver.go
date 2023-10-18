package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Data struct {
	Id      *string `json:"id"`
	Jsonrpc *string `json:"jsonrpc"`
	Result  *string `json:"result"`
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Check the message type (TextMessage in this case)
		if messageType == websocket.TextMessage {
			// Parse the JSON data
			var data Data
			err := json.Unmarshal(p, &data)
			if err != nil {
				log.Println(err)
				continue
			}

			// Process the received data
			log.Printf("Received Data: id=%s jsonrpc=%s, result=%s", *data.Id, *data.Jsonrpc, *data.Result)
		}

		if string(p) == "close-connection" {
			// Close the connection when a specific message is received
			return
		}

		// Process the message as needed
	}

}
