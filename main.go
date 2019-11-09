package main

import (
	auctioneerHandlers "auction-system/auctioneer/handlers"
	bidderHandlers "auction-system/bidder/handlers"
	driverUtils "auction-system/driverUtils"
	"fmt"
	"net/http"
	"time"
)

type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

const baseURL = "http://127.0.0.1"
const Port = "5000"
const createAuctionsAPI, listAuctionsAPI, createBidderAPI, listBiddersAPI, startBidRoundAPI = "/register-auction", "/list-auctions", "/create-bidder", "/list-bidders", "/bid-round/start"

func main() {
	m := http.NewServeMux()
	m = setRoutes(m)
	s := http.Server{Addr: ":" + Port,
		Handler: m,
		//ReadTimeout:  5 * time.Second,
		//WriteTimeout: 10 * time.Second,
		//IdleTimeout:  60 * time.Second,
	}
	go func() {
		fmt.Println("Starting server at port" + Port)
		err := s.ListenAndServe()
		if err != nil {
			fmt.Println("ListenAndServe: " + err.Error())
		}
	}()
	time.Sleep(5 * time.Second) // To Synchronize

	fmt.Println("Main server up !") // Check via ping

	urlPrepend := baseURL + ":" + Port

	// Create Auctions
	var url string
	url = urlPrepend + createAuctionsAPI
	driverUtils.CreateAuctions(2, url)

	// List Auctions
	url = urlPrepend + listAuctionsAPI
	driverUtils.GetAuctionsList(url)

	// Create Bidders
	url = urlPrepend + createBidderAPI
	driverUtils.CreateBidders(8, url)

	// List Bidders
	url = urlPrepend + listBiddersAPI
	driverUtils.GetBiddersList(url)

	// Start Auction Round
	url = urlPrepend + startBidRoundAPI
	driverUtils.StartBidRound(url)
	select {} // This will forever wait
}

func setRoutes(m *http.ServeMux)  *http.ServeMux{
	m.HandleFunc("/", handleRoot)
	// Auctioneer Routes
	m.HandleFunc("/bid-round/start", auctioneerHandlers.BidRoundHandler)         // Makes an  auction live.
	m.HandleFunc("/list-auctions", auctioneerHandlers.ListEndpointHandler)       // Lists all registered auctions.
	m.HandleFunc("/register-auction", auctioneerHandlers.RegisterAuctionHandler) // Register a new auction
	// Bidder Routes
	m.HandleFunc("/list-bidders", bidderHandlers.ListBiddersHandler)
	m.HandleFunc("/create-bidder", bidderHandlers.CreateAndRegisterBidderHandler) // Creates a bidder and registers it to Global bidders list.
	return m
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome to a diligent server!!")
}
