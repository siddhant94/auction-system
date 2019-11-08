package handlers

import (
	bidderModels "auction-system/bidder/bidderModels"
	bidderUtils "auction-system/bidder/utils"
	commonUtils "auction-system/commonUtils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	//"sync"
	"time"
)

var allBidders bidderModels.AppState
var allottedIds map[int]struct{}

func init() {
	allBidders = bidderModels.AppState{}
	allottedIds = map[int]struct{}{}
}

// TODO: This is currently create and register. Split registering to different function handler
func CreateAndRegisterBidderHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var bidder bidderModels.BidderStruct
	err := decoder.Decode(&bidder)
	if err != nil {
		fmt.Println(err)
		// TODO: Add http status code
		commonUtils.SendJSONResponse(w, map[string]string{"success": "false", "error": "true", "message": "Unable to decode Request body."})
		return
	}

	if len(allBidders.BidderList) >= 1000 {
		commonUtils.SendJSONResponse(w, map[string]string{"success": "false", "message": "BidderList at max capacity."})
		return
	}

	// Check against allotted id's if any repeat ID is given.
	// TODO: handle infinite loop condition
	for {
		newId := commonUtils.GetRandomInt()
		if _, found := allottedIds[newId]; !found {
			bidder.Id = newId
			break;
		}

	}
	err = bidderUtils.StartBidderServer(bidder, BidderNotificationHandler)
	if err != nil {
		commonUtils.SendJSONResponse(w, map[string]string{"success": "false", "error": "true", "message": err.Error()})
		return
	}
	allBidders.Lock()
	defer allBidders.Unlock()

	allBidders.BidderList = append(allBidders.BidderList, bidder)
	commonUtils.SendJSONResponse(w, map[string]string{"success": "true", "bidder_id": strconv.Itoa(bidder.Id)})
	return
}

func ListBiddersHandler(w http.ResponseWriter, r *http.Request) {
	biddersList := GetBiddersList()
	commonUtils.SendJSONResponse(w, biddersList)
}

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


func GetBiddersList() bidderModels.AppState {
	return allBidders
}