package handlers

import (
	auctionModels "auction-system/auctioneer/auctionModels"
	"auction-system/auctioneer/utils"
	"auction-system/bidder/bidderModels"
	bidderHandlers "auction-system/bidder/handlers"
	commonUtils "auction-system/commonUtils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	allAuctions auctionModels.AppState
	allottedIds map[int]struct{}
)

func init() {
	allAuctions = auctionModels.AppState{} // Initialize Global App State only once
	allottedIds = map[int]struct{}{}
}

// RegisterAuctionHandler : Get's optional Auction Name and creates an auction and updates Global App State for auctions list.
func RegisterAuctionHandler(w http.ResponseWriter, r *http.Request) {
	check := commonUtils.VerifyHTTPMethod(w, r, "POST")
	if check == false {
		w = commonUtils.SendMethodNotAllowed(w)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var auction auctionModels.AuctionStruct
	err := decoder.Decode(&auction)
	if err != nil {
		fmt.Println(err)
		commonUtils.SendJSONResponse(w, map[string]string{"success": "false", "error": "true", "message": "Unable to decode Request body."})
	}

	// Check against allotted id's if any repeat ID is given.
	// Random Id creation would run 100 times, and if still no unique Id, then drop the creation
	i := 0
	for i <= 100 {
		newId := commonUtils.GetRandomInt()
		if _, found := allottedIds[newId]; !found {
			auction.Id = newId
			break;
		} else {
			commonUtils.SendJSONResponse(w, map[string]string{"success": "false", "error": "true", "message": "Could not generate unique Auction"})
			return
		}

	}
	// Start lock for allAuctions i.e. AppState
	allAuctions.Lock()
	defer allAuctions.Unlock()


	allAuctions.AuctionList = append(allAuctions.AuctionList, auction)
	commonUtils.SendJSONResponse(w, map[string]string{"success": "true", "auction_id": strconv.Itoa(auction.Id)})
}

func ListEndpointHandler(w http.ResponseWriter, r *http.Request) {
	check := commonUtils.VerifyHTTPMethod(w, r, "GET")
	if check == false {
		w = commonUtils.SendMethodNotAllowed(w)
		return
	}
	commonUtils.SendJSONResponse(w, allAuctions)
}

// BidRoundHandler : Checks Auctions List and starts 1st Entry of the list.
func BidRoundHandler(w http.ResponseWriter, r *http.Request) {
	check := commonUtils.VerifyHTTPMethod(w, r, "GET")
	if check == false {
		w = commonUtils.SendMethodNotAllowed(w)
		return
	}
	listLen := len(allAuctions.AuctionList)
	var resp interface{}
	// Start an auction Round if list of auctions has entry
	if listLen <= 0 {
		resp = map[string]string{"auction_id": "0", "price": "null", "bidder_id": "null"}
		commonUtils.SendJSONResponse(w, resp)
		return
	}
	allAuctions.AuctionList, allAuctions.LiveAuction = utils.Remove(allAuctions.AuctionList, 0)

	// Create a channel to collect bid notification responses.
	bidEntriesChannel, biddersObj := make(chan bidderModels.BidResponse, 10), bidderHandlers.GetBiddersList()
	select {
	case bid := <-bidEntriesChannel:
		fmt.Println("received bid", bid)

	default:
		fmt.Println("no bids received")
	}

	// SendAuctionNotification: Takes bidders object and concurrently Notifies all bidders.
	go SendAuctionNotification(biddersObj, bidEntriesChannel)

	// Set 200 millisecond timer.
	timer := time.NewTimer(200 * time.Millisecond)
	<-timer.C
	// Close channel as timer is up.
	close(bidEntriesChannel)

	resp = map[string]string{"auction_id": strconv.Itoa(allAuctions.LiveAuction.Id), "price": "null", "bidder_id": "null"}
	highestBidder := bidderModels.BidResponse{}
	// If no bids received return
	if len(bidEntriesChannel) <= 0 {
		commonUtils.SendJSONResponse(w, resp)
	}
	for i := range bidEntriesChannel {
		if i.Price > highestBidder.Price {
			highestBidder.BidderId = i.BidderId
			highestBidder.Price = i.Price
		}
	}

	resp = map[string]string{"auction_id": strconv.Itoa(allAuctions.LiveAuction.Id), "price": fmt.Sprintf("%f", highestBidder.Price),
		"bidder_id": strconv.Itoa(highestBidder.BidderId)}
	commonUtils.SendJSONResponse(w, resp)
}
