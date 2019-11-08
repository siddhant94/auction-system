package utils

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

var (
	min, max int
)

func init() {
	rand.Seed(time.Now().UnixNano())
	min = 0
	max = 1000
}

func GetRandomID() int {
	randomNum := rand.Intn(max-min+1) + min
	return randomNum
}

func SendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}
