package auctionModels

import "sync"

type (
	AuctionStruct struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	AppState struct {
		sync.Mutex                  // <-- mutex protection
		AuctionList []AuctionStruct `json:"auction_list"`
		LiveAuction AuctionStruct   `json:"live_auction"`
	}
)