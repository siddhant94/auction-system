package handlers

import (
	auctionModels "auction-system/auctioneer/auctionModels"
	"auction-system/auctioneer/utils"
	"auction-system/bidder/bidderModels"
	bidderHandlers "auction-system/bidder/handlers"
	commonUtils "auction-system/commonUtils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var allAuctions auctionModels.AppState

func init() {
	allAuctions = auctionModels.AppState{} // Initialize Global App State only once
}

func BidEndpointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome to bid endpoint")
}

// RegisterAuctionHandler : Get's optional Auction Name and creates an auction and updates Global App State for auctions list.
func RegisterAuctionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var auction auctionModels.AuctionStruct
	err := decoder.Decode(&auction)
	if err != nil {
		fmt.Println(err)
		commonUtils.SendJSONResponse(w, map[string]string{"success": "false", "error": "true", "message": "Unable to decode Request body."})
	}
	// Start lock for allAuctions i.e. AppState
	allAuctions.Lock()
	defer allAuctions.Unlock()

	auction.Id = commonUtils.GetRandomInt()

	allAuctions.AuctionList = append(allAuctions.AuctionList, auction)
	log.Printf("%+v", allAuctions)
	commonUtils.SendJSONResponse(w, map[string]string{"success": "true", "auction_id": strconv.Itoa(auction.Id)})
}

func ListEndpointHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", allAuctions)
	commonUtils.SendJSONResponse(w, allAuctions)
}

// BidRoundHandler : Checks Auctions List and starts 1st Entry of the list.
func BidRoundHandler(w http.ResponseWriter, r *http.Request) {
	listLen := len(allAuctions.AuctionList)
	var resp interface{}
	// Start an auction Round if list of auctions has entry
	if listLen <= 0 {
		resp = map[string]string{"success": "false", "message": "No Auctions Listed. Register an auction first before starting bid round."}
		commonUtils.SendJSONResponse(w, resp)
		return
	}
	allAuctions.AuctionList, allAuctions.LiveAuction = utils.Remove(allAuctions.AuctionList, 0)
	resp = map[string]string{"success": "true", "auction_name": allAuctions.LiveAuction.Name, "auction_id": strconv.Itoa(allAuctions.LiveAuction.Id)}
	// TODO: Notify all bidders about Auction
	// Create a channel to collect bid notification responses. Default timeout - 200ms
	bidEntriesChannel := make(chan bidderModels.BidResponse, 10)
	biddersObj := bidderHandlers.GetBiddersList()
	var bidEntries []bidderModels.BidResponse
	select {
	case bid := <-bidEntriesChannel:
		bidEntries = append(bidEntries, bid)
		fmt.Println("received bid", bid)

	default:
		fmt.Println("no bids received")
	}

	go sendAuctionNotification(biddersObj, bidEntriesChannel)
	timer := time.NewTimer(2000 * time.Millisecond)
	<-timer.C
	close(bidEntriesChannel)
	fmt.Println("Timer finished")

	highestBidder :=  bidderModels.BidResponse{}
	for i := range bidEntriesChannel {
		if i.Price > highestBidder.Price {
			highestBidder.BidderId = i.BidderId
			highestBidder.Price = i.Price
		}
	}
	fmt.Println("Highest bidder")
	fmt.Printf("%+v", highestBidder)

	resp = map[string]string{"auction_id": strconv.Itoa(allAuctions.LiveAuction.Id), "price" : fmt.Sprintf("%f", highestBidder.Price),
	"highest_bidder_id": strconv.Itoa(highestBidder.BidderId)}
	commonUtils.SendJSONResponse(w, resp)
}

func sendAuctionNotification(biddersObj bidderModels.AppState, bidEntriesChannel chan bidderModels.BidResponse) {
	// Create an http client for making requests
	client := http.Client{Timeout: 200 * time.Millisecond}
	for _, v := range biddersObj.BidderList {
		url := "http://127.0.0.1:" + strconv.Itoa(v.Port) + "/register-notification"
		go sendRequests(client, url, bidEntriesChannel)
	}
}

func sendRequests(client http.Client, url string, channel chan bidderModels.BidResponse) bidderModels.BidResponse {
	var bidResp bidderModels.BidResponse
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Response error: ", err)
		return bidResp
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&bidResp)
	select {
	case channel <- bidResp:
		fmt.Println("sent message", bidResp)
	default:
		fmt.Println("no message sent")
	}
	return bidResp
}