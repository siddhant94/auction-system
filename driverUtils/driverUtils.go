package driverUtils

import (
	"auction-system/auctioneer/auctionModels"
	"auction-system/bidder/bidderModels"
	"auction-system/commonUtils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"time"
)

var (
	httpClient *http.Client
)

func init() {
	httpClient = &http.Client{
		Timeout: time.Second * 2,
	}
}

type auctionReq struct {
	Name string `json:"name"`
}

type bidderReq struct {
	Name string `json:"name"`
	Port int `json:"port"`
	delay time.Duration `json:"delay"`
}

type auctionResponse struct {
	Id string `json:"auction_id"`
	Success	  string `json:success`
}

type bidderResponse struct {
	Id string `json:"bidder_id"`
	Success	  string `json:success`
}

type StringBidResponse struct{
	BidderId string `json:"bidder_id"`
	Price    string `price:"price"`
}

func CreateAuctions(num int, url string) {
	fmt.Println("Creating " + strconv.Itoa(num) + " Auctions")
	var reqData auctionReq
	var respData []byte
	i := 1
	for i <= num {
		reqData.Name = "Auction Item " + strconv.Itoa(i)
		i++
		data, err := json.Marshal(reqData)
		if err != nil{
			fmt.Println("JSON Marshall Error " + err.Error())
			continue
		}
		respData, err = doHTTPRequest("POST", url, data)
		if err != nil {
			fmt.Println("Could not create Auction with name - ", reqData.Name, err)
			continue
		}
		var apiResp auctionResponse
		err = json.Unmarshal(respData, &apiResp)
		if err != nil {
			fmt.Println("JSON Ummarshall failed " + err.Error())
			continue
		}
		if apiResp.Success == "true" {
			fmt.Println("Created an Auction with ID - " + apiResp.Id)
		}
	}
	return
}

func GetAuctionsList(url string) {
	fmt.Println("Fetching Auctions List")
	var auctionsList auctionModels.AppState
	respData, err := doHTTPRequest("GET", url, nil)
	if (err != nil) {
		fmt.Println("Error fetching Auctions Data "+ err.Error())
		return
	}
	err = json.Unmarshal(respData, &auctionsList)
	if err != nil {
		fmt.Println("JSON Ummarshall failed " + err.Error())
	}
	fmt.Printf("%+v\n", auctionsList)
	return
}

func CreateBidders(num int, url string) {
	i := 0
	var bidderRequest bidderReq
	var respData []byte
	for i < num {
		bidderRequest.Name = "Bidder " + strconv.Itoa(i)
		bidderRequest.Port = 2000 + i
		bidderRequest.delay = time.Duration(commonUtils.GetRandomInt(10, 500))
		i++
		data, err := json.Marshal(bidderRequest)
		if err != nil{
			fmt.Println("JSON Marshall Error " + err.Error())
			continue
		}
		respData, err = doHTTPRequest("POST", url, data)
		if err != nil {
			fmt.Println("Could not create Bidder with name - ", bidderRequest.Name, err)
			continue
		}
		var apiResp bidderResponse
		err = json.Unmarshal(respData, &apiResp)
		if err != nil {
			fmt.Println("JSON Ummarshall failed " + err.Error())
			continue
		}
		if apiResp.Success == "true" {
			fmt.Println("Created a Bidder live with ID - " + apiResp.Id)
		}
	}
	return
}

func GetBiddersList(url string) {
	fmt.Println("Fetching Bidders List")
	var biddersList bidderModels.AppState
	respData, err := doHTTPRequest("GET", url, nil)
	if (err != nil) {
		fmt.Println("Error fetching Bidders Data "+ err.Error())
		return
	}
	err = json.Unmarshal(respData, &biddersList)
	if err != nil {
		fmt.Println("JSON Ummarshall failed " + err.Error())
	}
	fmt.Printf("%+v\n", biddersList)
	return
}

func StartBidRound(url string) {
	fmt.Println("Staring Bid Round")
	var resp StringBidResponse
	jsonBody, err := doHTTPRequest("GET", url, nil)
	if (err != nil) {
		fmt.Println("Error in getting Start Bid Round Response "+ err.Error())
		return
	}
	err = json.Unmarshal(jsonBody, &resp)
	if err != nil {
		fmt.Println("JSON Ummarshall failed " + err.Error())
	}
	fmt.Printf("%+v\n", resp)
	return
}

func doHTTPRequest(method string, url string, data []byte) ([]byte, error) {
	var req *http.Request
	var err error
	if method == "POST" {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			fmt.Println("Error creating POST REQ " + err.Error())
		}
	} else if method == "GET" {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating GET REQ " + err.Error())
		}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Request Failed for -  ", url, err)
		return nil, err
	}
	//if resp.StatusCode <= 200 || resp.StatusCode >= 209 {
	//	return nil, errors.New("Api status is not 200")
	//}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Could not read response for ", url, err)
		return nil, err
	}
	return body, nil
}