package airports

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AirportHttpClient interface {
	GET(AirportRequest) ([]AirportResponse, error)
}

func NewAirportHttpClient(apikey string) AirportHttpClient {
	return &httpClient{
		ApiKey: apikey,
	}
}

const airportUrl = "https://airlabs.co/api/v9/airports"

type httpClient struct {
	ApiKey string `json:"api_key"` //https://airlabs.co/account
}

type httpResponse struct {
	Response []AirportResponse `json:"response"`
}

func (h *httpClient) GET(req AirportRequest) ([]AirportResponse, error) {
	params := url.Values{}
	Url, err := url.Parse(airportUrl)
	if err != nil {
		return nil, err
	}
	if req.IataCode != "" {
		params.Set("iata_code", req.IataCode)
	}
	if req.IcaoCode != "" {
		params.Set("icao_code", req.IcaoCode)
	}
	if req.CityCode != "" {
		params.Set("city_code", req.CityCode)
	}
	if req.CountryCode != "" {
		params.Set("country_code", req.CountryCode)
	}
	params.Set("api_key", h.ApiKey)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Printf("request: %v", urlPath)
	resp, err := http.Get(urlPath)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res httpResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res.Response, nil
}
