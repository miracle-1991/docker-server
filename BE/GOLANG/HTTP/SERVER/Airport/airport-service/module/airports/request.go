package airports

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AirportHttpClient interface {
	GET(AirportRequest) ([]Airport, error)
	UPDATE(AirportRequest) ([]Airport, error)
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
	Response []Airport `json:"response"`
}

func (h *httpClient) GET(req AirportRequest) ([]Airport, error) {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if req.IataCode != "" {
		maps["iata_code"] = req.IataCode
	} else if req.IcaoCode != "" {
		maps["icao_code"] = req.IcaoCode
	} else if req.CityCode != "" {
		maps["city_code"] = req.CityCode
	} else if req.CountryCode != "" {
		maps["country_code"] = req.CountryCode
	}
	airports, err := GetAirports(req.PageNum, req.PageSize, maps)
	if err != nil {
		return nil, err
	}
	return airports, nil
}

func (h *httpClient) UPDATE(req AirportRequest) ([]Airport, error) {
	params := url.Values{}
	Url, err := url.Parse(airportUrl)
	if err != nil {
		return nil, err
	}
	params.Set("api_key", h.ApiKey)
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
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Printf("request: %v\n", urlPath)
	resp, err := http.Get(urlPath)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res httpResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	airports := res.Response
	for _, a := range airports {
		var exist bool
		if a.IcaoCode != "" {
			exist, err = ExistAirportByIcaoCode(a.IcaoCode)
		} else if a.IataCode != "" {
			exist, err = ExistAirportByIataCode(a.IataCode)
		} else {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !exist {
			err = AddAirport(a)
			if err != nil {
				return nil, err
			}
		}
	}
	return airports, nil
}
