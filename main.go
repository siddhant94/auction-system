package main

import (
	auctioneerHandlers "auction-system/auctioneer/handlers"
	bidderHandlers "auction-system/bidder/handlers"
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
	http.HandleFunc("/", handleRoot)
	// Auctioneer Routes
	http.HandleFunc("/bid-round/start", auctioneerHandlers.BidRoundHandler) // Makes an  auction live.
	http.HandleFunc("/list-auctions", auctioneerHandlers.ListEndpointHandler) // Lists all registered auctions.
	http.HandleFunc("/register-auction", auctioneerHandlers.RegisterAuctionHandler) // Register a new auction
	// Bidder Routes
	http.HandleFunc("/list-bidders", bidderHandlers.ListBiddersHandler)
	http.HandleFunc("/create-bidder", bidderHandlers.CreateAndRegisterBidderHandler) // Creates a bidder and registers it to Global bidders list.
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome to a diligent server!!")
}