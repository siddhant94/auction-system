package bidderModels

import (
	"net/http"
	"sync"
	"time"
)

type (
	BidderStruct struct {
		Id    int           `json:"id"`
		Name  string        `json:"name"`
		Delay time.Duration `json:"delay"`
		Port  int           `json:"port"`
		Url   string        `json:"url"`
	}
	AppState struct {
		sync.Mutex // <-- mutex protection
		BidderList []BidderStruct
	}
)

type BidResponse struct {
	Price    float32 `json:"price"`
	BidderId int `json:"bidder_id"`
}

type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)
