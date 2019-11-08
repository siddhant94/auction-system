package utils

import (
	bidderModels "auction-system/bidder/bidderModels"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	//"sync"
)

type bidderStateHandler func(time.Duration, int) bidderModels.RequestHandlerFunction

func StartBidderServer(bidder bidderModels.BidderStruct, handlerFunc bidderStateHandler) error {
	fmt.Println("Creating new server")
	m := http.NewServeMux()
	s := http.Server{Addr: ":" + strconv.Itoa(bidder.Port), Handler: m}
	fmt.Println("Starting server at " + s.Addr)
	m.HandleFunc("/register-notification", handlerFunc(bidder.Delay, bidder.Id))
	if err := s.ListenAndServe(); err != nil {
		log.Print("ListenAndServe: ", err)
		return err
	}
	return nil
}