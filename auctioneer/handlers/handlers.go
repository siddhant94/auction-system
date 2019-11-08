package handlers

import (
	auctionModels "auction-system/auctioneer/models/auctionModels"
	"auction-system/auctioneer/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var allAuctions auctionModels.AppState

func init() {
	allAuctions = auctionModels.AppState{} // Initialize Global App State only once
}

func BidEndpointHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome to bid endpoint")
}

func RegisterAuctionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var auction auctionModels.AuctionStruct
	err := decoder.Decode(&auction)
	if err != nil {
		panic(err)
	}
	// Start lock for allAuctions i.e. AppState
	allAuctions.Lock()
	defer allAuctions.Unlock()

	auction.Id = utils.GetRandomID()
	/* TODO: Check new Id against  a list of all Id's currently used in State and if repeated get repeated Number. Also,
	check count of auctions present, return error when size limit reaches, currently its 1000 (max - min) present in utils.
	*/
	allAuctions.AuctionList = append(allAuctions.AuctionList, auction)
	log.Printf("%+v", allAuctions)
	utils.SendJSONResponse(w, map[string]string{"success": "true", "auction_id": strconv.Itoa(auction.Id)})
}

func ListEndpointHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", allAuctions)
	utils.SendJSONResponse(w, allAuctions)
}
