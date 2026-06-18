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

	// Accept one audio and one video track incoming
	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		_, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionSendrecv,
		})

		if err != nil {
			fmt.Errorf("Failed to add transceiver: %v", err)
			return
		}
	}

}
