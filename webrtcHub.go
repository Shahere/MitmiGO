package main

import (
	"errors"
	"slices"

	"golang.org/x/exp/maps"
)

var webRTCHubs = make([]*WebRTCHub, 0)

type WebRTCHub struct {
	name       string
	clients    map[*WebRTCClient]bool
	register   chan *WebRTCClient
	unregister chan *WebRTCClient
}

func getOrCreateWebRTCHub(hubName string) (*WebRTCHub, error) {
	if hubName == "" {
		return nil, errors.New("Empty webRTCHub name")
	}
	var webRTCHub *WebRTCHub
	iHub := slices.IndexFunc(webRTCHubs, func(webRTCHub *WebRTCHub) bool {
		return webRTCHub.name == hubName
	})
	if iHub == -1 {
		webRTCHub = newWebRTCHub(hubName)
	} else {
		webRTCHub = webRTCHubs[len(hubs)-1]
	}

	return webRTCHub, nil
}

func newWebRTCHub(hubName string) *WebRTCHub {
	webRTCHub := WebRTCHub{
		name:       hubName,
		clients:    make(map[*WebRTCClient]bool),
		register:   make(chan *WebRTCClient),
		unregister: make(chan *WebRTCClient),
	}
	go webRTCHub.run()
	webRTCHubs = append(webRTCHubs, &webRTCHub)
	return &webRTCHub
}

func (webRTCHub *WebRTCHub) getClients() []*WebRTCClient {
	keys := maps.Keys(webRTCHub.clients)
	return keys
}

func (webRTCHub *WebRTCHub) run() {
	for {
		select {
		case webRTCClient := <-webRTCHub.register:
			webRTCHub.clients[webRTCClient] = true
		case webRTCClient := <-webRTCHub.unregister:
			delete(webRTCHub.clients, webRTCClient)
		}
	}
}

func (webRTCHub *WebRTCHub) getClientsById(id string) (*WebRTCClient, error) {
	if id == "" {
		return nil, errors.New("Empty ID")
	}
	webRTCClients := webRTCHub.getClients()
	iClientToSend := slices.IndexFunc(webRTCClients, func(webRTCClient *WebRTCClient) bool {
		return webRTCClient.id == id
	})
	if iClientToSend == -1 {
		return nil, errors.New("No client found")
	}

	return webRTCClients[iClientToSend], nil
}
