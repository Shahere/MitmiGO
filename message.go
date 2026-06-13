package main

import (
	"encoding/json"
	"fmt"
)

type ActionType string

const (
	Close ActionType = "close"
	Join  ActionType = "join"
	GetId ActionType = "getid"
)

type Message struct {
	Target  string      `json:"target"`
	From    ContactInfo `json:"from"`
	Payload PayloadType `json:"payload"`
}

type ContactInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RTCSessionDescriptionInit struct {
	Type string `json:"name"`
	Sdp  string `json:"sdp"`
}

type PayloadType struct {
	Action     ActionType                `json:"action"`
	HubName    string                    `json:"hubname"`
	Message    string                    `json:"message"`
	Disconnect string                    `json:"disconnect"`
	Sdp        RTCSessionDescriptionInit `json:"sdp"`
	Candidate  any                       `json:"candidate"`
	Mesh       bool                      `json:"mesh"`
}

func DecodeMessage(data []byte) (*Message, error) {
	var message Message
	err := json.Unmarshal(data, &message)
	if err != nil {
		fmt.Printf("Error while decoding message\n")
		fmt.Printf("%s", err)
		return nil, err
	}
	return &message, nil
}

func (message *Message) encodeMessage() ([]byte, error) {
	s, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error while encoding message")
		fmt.Printf("%s", err)
		return nil, err
	}
	return s, nil
}
