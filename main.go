package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.Host)
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
	fmt.Printf("Start server...\n")

	hub := newHub()
	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s\n", r.Host)
		fmt.Printf("got / request\n")
		clients := hub.getClients()
		ipAddrClient := make([]string, len(clients))

		for i := 0; i < len(clients); i++ {
			ipAddrClient[i] = clients[i].conn.RemoteAddr().String()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ipAddrClient)
	})
	http.HandleFunc("/hello", getHello)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		serveWs(hub, w, r)
	})

	port := 3333
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closes")
	} else if err != nil {
		fmt.Printf("Error server starting... %s\n", err.Error())
		os.Exit(1)
	}
}
