package handlers

import (
	"auction-system/bidder/bidderModels"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const urlPrepend, auctionNotifURL = "http://127.0.0.1:", "/auction-notification"

func SendAuctionNotification(biddersObj bidderModels.AppState, bidEntriesChannel chan bidderModels.BidResponse) {
	// Create an http client for making requests.
	client := http.Client{Timeout: 200 * time.Millisecond}
	for _, v := range biddersObj.BidderList {
		url := urlPrepend + strconv.Itoa(v.Port) + auctionNotifURL
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
