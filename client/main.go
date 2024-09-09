package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var (
    username string
    deleteLastLine bool
)

type Message struct {
    Author string `json:"author"`
    Message string `json:"message"`
}

func handleIO(conn *websocket.Conn) {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("-> ")
        deleteLastLine = true
        input, err := reader.ReadString('\n')

        if err != nil {
            log.Println(err)
        }

        // CRLF to LF

        input = strings.Replace(input, "\n", "", -1)

        conn.WriteJSON(Message{Author: username, Message: input})
        deleteLastLine = false
    }
}

func handleIncoming(conn *websocket.Conn) {
    for {
        _, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            break
        }

        var payload Message
        unmarshalErr := json.Unmarshal(p, &payload)

        if unmarshalErr != nil {
            fmt.Println(err)
        }

        //fmt.Printf("last line content: %v", lastLine)
        //if deleteLastLine {
            //fmt.Printf("\033[1A\033[1G\033[K")
        //}
        deleteLastLine = false
        fmt.Printf("\n%v: %v\n", payload.Author, payload.Message)
        fmt.Print("-> ")
    }
}


func main() {
    //fmt.Printf("Command line arguments (%v): %v", len(os.Args), os.Args)

    if len(os.Args) < 2 {
        fmt.Println("No url provided. Exiting program...")
        return
    }

    // first handle username
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter a username: ")
    input, err := reader.ReadString('\n')

    if err != nil {
        log.Println(err)
    }

    // CRLF to LF and store username
    username = strings.Replace(input, "\n", "", -1)

    url := os.Args[1]
    header := http.Header{}
    conn, _, err := websocket.DefaultDialer.Dial(url, header)

    if err != nil {
        log.Println(err)

    }

    fmt.Printf("Connected to: %v\n", os.Args[1])

    // send who is the client and then start handling incoming msgs and IO
    conn.WriteMessage(websocket.TextMessage, []byte(username))

    go handleIO(conn)
    handleIncoming(conn)

}
