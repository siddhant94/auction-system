package utils

import (
	"auction-system/auctioneer/auctionModels"
)

// Remove() : Takes slice and index as input and removes it. (Does not maintains order and does not performs bounds-checking).
func Remove(s []auctionModels.AuctionStruct, i int) ([]auctionModels.AuctionStruct, auctionModels.AuctionStruct) {
	removedItem := s[i]
	s[i] = s[len(s)-1]
	return s[:len(s)-1], removedItem
}
