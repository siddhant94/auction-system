package auctionModels

import "sync"

type (
	AuctionStruct struct {
		Id   int
		Name string
	}
	AppState struct {
		sync.Mutex  // <-- this mutex protects
		AuctionList []AuctionStruct
	}
)