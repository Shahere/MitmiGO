package main

import (
	"errors"
	"slices"

	"golang.org/x/exp/maps"
)

var hubs = make([]*Hub, 0)

type BroadcastType struct {
	message []byte
	client  *Client
}

type Hub struct {
	name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastType
}

func getOrCreateHub(hubName string) (*Hub, error) {
	if hubName == "" {
		return nil, errors.New("Empty hub name")
	}
	var hub *Hub
	iHub := slices.IndexFunc(hubs, func(hub *Hub) bool {
		return hub.name == hubName
	})
	if iHub == -1 {
		hub = newHub(hubName)
	} else {
		hub = hubs[len(hubs)-1]
	}

	return hub, nil
}

func newHub(hubName string) *Hub {
	hub := Hub{
		name:       hubName,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastType),
	}
	go hub.run()
	hubs = append(hubs, &hub)
	return &hub
}

func (hub *Hub) getClients() []*Client {
	keys := maps.Keys(hub.clients)
	return keys
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			delete(hub.clients, client)
		case broadcastType := <-hub.broadcast:
			//fmt.Printf("Broadcast message : %s\n", broadcastType.message)
			for k := range hub.clients {
				if k == broadcastType.client {
					continue
				}
				k.toSend <- broadcastType.message
			}
		}
	}
}

func (hub *Hub) getClientsById(id string) (*Client, error) {
	if id == "" {
		return nil, errors.New("Empty ID")
	}
	clients := hub.getClients()
	iClientToSend := slices.IndexFunc(clients, func(client *Client) bool {
		return client.id == id
	})
	if iClientToSend == -1 {
		return nil, errors.New("No client found")
	}

	return clients[iClientToSend], nil
}
