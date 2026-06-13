package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	id     string
	name   string
	hub    *Hub
	conn   *websocket.Conn
	toSend chan []byte
}

func (client *Client) write(conn *websocket.Conn) {
	defer func() {
		client.hub.unregister <- client
		client.conn.Close()
	}()

	for {
		select {
		case message := <-client.toSend:
			//fmt.Printf("Sending message %s\n", message)
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Printf("Error while writing message : %s", message)
				break
			}
		}
	}
}

func (client *Client) read(conn *websocket.Conn) {
	defer func() {
		if client.hub != nil {
			client.hub.unregister <- client
		}
		client.conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error while reading message\n")
			break
		}
		//fmt.Printf("Message Receive : %s\n", message)

		newMessage, err := DecodeMessage(message)
		if err != nil {
			break
		}

		if !client.verifyIdentity(*newMessage) {
			fmt.Printf("Identity is not correct !")
			break
		}

		if client.hub == nil {
			hub, _ := getOrCreateHub(newMessage.Payload.HubName)
			client.hub = hub
			client.hub.register <- client
		}

		if newMessage.Payload.Action == "join" && newMessage.Payload.Mesh == false {
			//TODO Connect with server instead of other client
			fmt.Println("SFU is asked")
			continue
		}

		if newMessage.Target == "" {
			client.hub.broadcast <- BroadcastType{message, client}
			continue
		}
		clientToSend, err := client.hub.getClientsById(newMessage.Target)
		if err != nil {
			fmt.Printf("No client to send %s", newMessage.Target)
		}
		clientToSend.toSend <- message

	}
}

func (client *Client) verifyIdentity(newMessage Message) bool {
	if newMessage.From.Name == "" {
		return false
	}
	if newMessage.Payload.Action == "join" {
		if client.name != "" {
			return false
		}
		client.name = newMessage.From.Name

		if client.id != "" {
			return false
		}
		client.id = newMessage.From.Id
	}
	if client.name != newMessage.From.Name {
		return false
	}
	return true
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade connection to : %s\n", err.Error())
		return
	}
	fmt.Printf("New client : %s\n", conn.RemoteAddr())
	client := &Client{
		hub:    nil,
		conn:   conn,
		toSend: make(chan []byte),
	}

	go client.read(conn)
	go client.write(conn)
}
