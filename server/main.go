package main

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}
var connectionsUsers map[*websocket.Conn]string

type Message struct {
    Author string `json:"author"`
    Message string `json:"message"`
}

func closeAndRemoveConn(conn *websocket.Conn) {
    conn.Close()
    newConnections := make(map[*websocket.Conn]string, 0)

    for k, v := range connectionsUsers {
        if k != conn {
            newConnections[k] = v
        }
    }

    connectionsUsers = newConnections
}

func broadcast(p []byte, sender *websocket.Conn) {

    for c, _ := range connectionsUsers {
        if c != sender {
            var payload Message

            json.Unmarshal(p, &payload)

            if err := c.WriteJSON(payload); err != nil {
                log.Println(err)
                return
            }
        }
    }
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "this is the home route")
}


func handleWs(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    // remember defer is LIFO
    //defer conn.Close()
    defer closeAndRemoveConn(conn)

    if err != nil {
        log.Println(err)
        return
    }

    receivedUsername := false

    for !receivedUsername {
        _, p, err := conn.ReadMessage() // first return value is the messageType
        if err != nil {
            log.Println(err)
            break
        }

        connectionsUsers[conn] = string(p)
        receivedUsername = true
    }

    fmt.Printf("Amount of connections: %v\n", len(connectionsUsers))

    for {
        _, p, err := conn.ReadMessage() // first return value is the messageType
        if err != nil {
            log.Println(err)
            break
        }

        broadcast(p, conn)
    }
}

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    connectionsUsers = make(map[*websocket.Conn]string)

    http.HandleFunc("/", handleHome)
    http.HandleFunc("/chat", handleWs)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
