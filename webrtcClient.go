package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type WebRTCClient struct {
	id   string
	name string
	hub  *Hub
	conn *websocket.Conn
}

func (webRTCClient *WebRTCClient) createWebRTCConnection(conn *websocket.Conn) {
	defer conn.Close()

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		fmt.Println("Error while creating peer connection %v", err)
		return
	}

	defer peerConnection.Close()

}
