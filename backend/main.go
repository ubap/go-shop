package main

import (
	"log"
	"net/http"
)

func main() {
	dispatcher := NewMethodDispatcher()
	dispatcher.Register("itemAddedToBuy", handleItemAdded)
	dispatcher.Register("itemBought", handleItemBought)

	hub := NewHub(dispatcher)
	go hub.Run()

	fs := http.FileServer(http.Dir("../frontend/dist"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})

	log.Println("HTTP server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
