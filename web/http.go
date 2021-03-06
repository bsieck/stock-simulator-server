package web

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"github.com/stock-simulator-server/client"
	"os"
	"fmt"
)

var clients = make(map[*websocket.Conn]http.Client) // connected clients

func StartHandlers() {
	shareDir := os.Getenv("FILE_SERVE")
	if shareDir == ""{
		shareDir = "static"
	}
	port := os.Getenv("PORT")
	if port == ""{
		port = "8000"
	}
	port = ":" + port
	fmt.Println(shareDir)
	var fs = http.FileServer(http.Dir(shareDir))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	err := http.ListenAndServe(port, nil)
	if err != nil{
		log.Fatal("ListenAndServe:", err)
	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	//first upgrade the connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		return
	}
	socketRX := make(chan string)
	socketTX := make(chan string)
	// Gate Keeper
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}
		loginErr := client.Login(string(msg), socketTX, socketRX)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(loginErr.Error()))
		} else {
			break
		}

	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	go runTxSocket(ws, socketTX)
	rxSocket(ws, socketRX)
}

func runTxSocket(conn *websocket.Conn, tx chan string){
	for str := range tx{
		conn.WriteMessage(websocket.TextMessage, []byte(str))
	}
}

func rxSocket(conn *websocket.Conn, rx chan string){
	for{
		_, msg, err := conn.ReadMessage()
		if err != nil{
			break
		}
		rx <- string(msg)
	}
}
