package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type WebRTCClient struct {
	id             string
	name           string
	webRTCHub      *WebRTCHub
	conn           *websocket.Conn
	peerConnection []*webrtc.PeerConnection
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

	webRTCClient.peerConnection = append(webRTCClient.peerConnection, peerConnection)

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			fmt.Printf("Fail to marshal candidate to JSON : %v", err)
			return
		}
		fmt.Printf("Send candidate to client %s", candidateString)
		writeError := conn.WriteJSON(&Message{
			From: ContactInfo{
				Id:   webRTCClient.id,
				Name: webRTCClient.name,
			},
			Target: webRTCClient.id,
			Payload: PayloadType{
				Action:    "ice",
				Candidate: candidateString,
				HubName:   webRTCClient.webRTCHub.name,
			},
		})
		if writeError != nil {
			fmt.Printf("Error writing Ice candidate")
		}
	})
}
