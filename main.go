package main

import (
	handlers "auction-system/auctioneer/handlers"
	"fmt"
	"log"
	"net/http"
)

type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

func main() {
	setRoutes()
	if err := http.ListenAndServe("localhost:5000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//http.ListenAndServe(":5000", nil)
}

func setRoutes() {
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/bid", handlers.BidEndpointHandler) // Accepts ad-request
	http.HandleFunc("/list-auctions", handlers.ListEndpointHandler) // Accepts ad-request
	http.HandleFunc("/register-bidder", handlers.BidEndpointHandler) // Accepts ad-request

	http.HandleFunc("/register-auction", handlers.RegisterAuctionHandler) // Accepts ad-request
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome to a diligent server!!")
}