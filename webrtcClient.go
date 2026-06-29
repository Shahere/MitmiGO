package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

type WebRTCClient struct {
	id             string
	name           string
	webRTCHub      *WebRTCHub
	conn           *websocket.Conn
	peerConnection []*webrtc.PeerConnection
	stream         *Stream
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

	peerConnection.OnICEConnectionStateChange(func(is webrtc.ICEConnectionState) {
		fmt.Printf("ICE connection state changed: %s", is)
	})

	peerConnection.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		fmt.Printf("Peer connection has changed %s", pcs)

		switch pcs {
		case webrtc.PeerConnectionStateFailed:
			err := peerConnection.Close()
			if err != nil {
				fmt.Printf("Cant close peer connection : %v", err)
			}
		case webrtc.PeerConnectionStateClosed:
			//Leaving ?
		}
	})

	peerConnection.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		fmt.Printf("Got remote track: Kind=%s, ID=%s, PayloadType=%d", tr.Kind(), tr.ID(), tr.PayloadType())
		localTrack := webRTCClient.addTrack(tr)

		buffer := make([]byte, 1500)
		rtpPkt := &rtp.Packet{}

		for {
			i, _, err := tr.Read(buffer)
			if err != nil {
				return
			}

			err = rtpPkt.Unmarshal(buffer[:i])
			if err != nil {
				fmt.Printf("Failed to unmarshal incoming RTP packet: %v", err)
				return
			}

			rtpPkt.Extension = false
			rtpPkt.Extensions = nil

			// Send back if normal
			err = localTrack.WriteRTP(rtpPkt)
			if err != nil {
				return
			}
		}
	})

	//TODO Signaling peers and send client message all this stuff
}

//*******************************************

func (webRTCClient *WebRTCClient) addTrack(track *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	localTrack, err := webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability, track.ID(), track.StreamID())
	if err != nil {
		fmt.Printf("Fail to get local track")
	}

	webRTCClient.stream = newStream(localTrack.ID(), localTrack.Kind())
	//TODO defer remove stream and release pointer?
	return localTrack
}
