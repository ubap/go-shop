package main

import (
	"go-shop/backend/basket"
	"log"
	"net/http"
)

func main() {
	memoryBasket := basket.NewInMemoryBasket()
	hub := NewHub(memoryBasket)
	go hub.Run()

	fs := http.FileServer(http.Dir("../frontend/dist"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//basketId := r.URL.Query().Get("basket")

		ServeWs(hub, w, r, hub.protocol)
	})

	log.Println("HTTP server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
