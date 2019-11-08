package utils

import (
	bidderModels "auction-system/bidder/bidderModels"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	minTCPPort         = 0
	maxTCPPort         = 65535
	maxReservedTCPPort = 1024
	maxRandTCPPort     = maxTCPPort - (maxReservedTCPPort + 1)
)

type bidderStateHandler func(time.Duration, int) bidderModels.RequestHandlerFunction

func StartBidderServer(bidder bidderModels.BidderStruct, handlerFunc bidderStateHandler) error {
	m := http.NewServeMux()
	s := http.Server{Addr: ":" + strconv.Itoa(bidder.Port), Handler: m}

	fmt.Println("Starting server at " + s.Addr)
	m.HandleFunc("/auction-notification", handlerFunc(bidder.Delay, bidder.Id))
	if err := s.ListenAndServe(); err != nil {
		log.Print("ListenAndServe: ", err)
		return err
	}
	return nil
}

func IsTCPPortAvailable(port int) bool {
	if port < minTCPPort || port > maxTCPPort {
		return false
	}
	conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}