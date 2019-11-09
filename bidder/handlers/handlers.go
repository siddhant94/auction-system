package handlers

import (
	"auction-system/bidder/bidderModels"
	bidderUtils "auction-system/bidder/utils"
	commonUtils "auction-system/commonUtils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	allBidders  bidderModels.AppState
	allottedIds map[int]struct{}
)

func init() {
	allBidders = bidderModels.AppState{}
	allottedIds = map[int]struct{}{}
}

// TODO: This is currently create and register. Split registering to different function handler
func CreateAndRegisterBidderHandler(w http.ResponseWriter, r *http.Request) {
	check := commonUtils.VerifyHTTPMethod(w, r, "POST")
	if check == false {
		w = commonUtils.SendMethodNotAllowed(w)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var bidder bidderModels.BidderStruct
	err := decoder.Decode(&bidder)
	if err != nil {
		fmt.Println(err)
		commonUtils.SendJSONResponse(w, map[string]string{"success": "false", "error": "true", "message": "Unable to decode Request body."})
		return
	}

	if len(allBidders.BidderList) >= 1000 {
		commonUtils.SendJSONResponse(w, map[string]string{"success": "false", "message": "BidderList at max capacity."})
		return
	}

	// Check against allotted id's if any repeat ID is given.
	// Random Id creation would run 100 times, and if still no unique Id, then drop the creation
	i := 0
	for i <= 100 {
		newId := commonUtils.GetRandomInt(1, 1000)
		if _, found := allottedIds[newId]; !found {
			bidder.Id = newId
			allottedIds[newId] = struct{}{}
			break;
		} else {
			commonUtils.SendJSONResponse(w, map[string]string{"success": "false", "error": "true", "message": "Could not generate unique Bidder"})
			return
		}

	}
	if !bidderUtils.IsTCPPortAvailable(bidder.Port) {
		commonUtils.SendJSONResponse(w, map[string]string{"success": "false", "error": "true", "message": "Port in use."})
		return
	}

	go bidderUtils.StartBidderServer(bidder, BidderNotificationHandler)

	allBidders.Lock()
	defer allBidders.Unlock()

	allBidders.BidderList = append(allBidders.BidderList, bidder)
	commonUtils.SendJSONResponse(w, map[string]string{"success": "true", "bidder_id": strconv.Itoa(bidder.Id)})
	return
}

func ListBiddersHandler(w http.ResponseWriter, r *http.Request) {
	check := commonUtils.VerifyHTTPMethod(w, r, "GET")
	if check == false {
		w = commonUtils.SendMethodNotAllowed(w)
		return
	}
	biddersList := GetBiddersList()
	commonUtils.SendJSONResponse(w, biddersList)
}

//BidderNotificationHandler : Takes in a delay and Id, which it passes on to func that returns RequestHandlerFunction
func BidderNotificationHandler(t time.Duration, id int) bidderModels.RequestHandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(t * time.Millisecond)

		var bidResp bidderModels.BidResponse
		bidResp.BidderId = id
		bidResp.Price = commonUtils.GetRandomFloat()
		fmt.Println("BidCreated")
		fmt.Printf("%+v", bidResp)
		commonUtils.SendJSONResponse(w, bidResp)
	}
}

// GetBiddersList : Returns the app state for bidders module
func GetBiddersList() bidderModels.AppState {
	return allBidders
}
