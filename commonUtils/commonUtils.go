package commonUtils

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

var (
	min, max int
	minFloat, maxFloat float32
)

func init() {
	rand.Seed(time.Now().UnixNano())
	min = 0
	max = 1000
	minFloat = 0.00
	maxFloat = 1000.00
}

func GetRandomInt() int {
	randomNum := rand.Intn(max-min+1) + min
	return randomNum
}

func GetRandomFloat() float32 {
	randomNum := minFloat + rand.Float32() * (maxFloat - minFloat)
	return randomNum
}

func SendJSONResponse(w http.ResponseWriter, data interface{}) {
	//w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}