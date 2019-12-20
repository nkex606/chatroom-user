package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var addr = os.Getenv("addr")
var port = os.Getenv("port")
var forever chan struct{}

type user struct {
	conn *websocket.Conn
}

func newUser() *user {
	u := &user{}
	return u
}

func (u *user) dial() {
	var err error
	chatroomURL := fmt.Sprintf("ws://%s:%s/ws", addr, port)
	u.conn, _, err = websocket.DefaultDialer.Dial(chatroomURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected!")

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Println("Quit chatroom.")
		u.conn.Close()
		close(forever)
	}()

	go u.read()
	go u.send()
}

func (u *user) read() {
	for {
		_, p, err := u.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Print("client:", string(p))
	}
}

func (u *user) send() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if err := u.conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			break
		}
	}
}

func main() {
	forever = make(chan struct{})
	u := newUser()
	go u.dial()
	<-forever
}
