package schedules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ScheduleHttpClient interface {
	GET(ScheduleRequest) ([]ScheduleResponse, error)
}

func NewScheduleHttpClient(apikey string) ScheduleHttpClient {
	return &httpClient{
		ApiKey: apikey,
	}
}

const scheduleUrl = "https://airlabs.co/api/v9/schedules"

type httpClient struct {
	ApiKey string `json:"api_key"` //https://airlabs.co/account
}

type scheduleResponse struct {
	Response []ScheduleResponse `json:"response"`
}

func (h *httpClient) GET(req ScheduleRequest) ([]ScheduleResponse, error) {
	params := url.Values{}
	Url, err := url.Parse(scheduleUrl)
	if err != nil {
		return nil, err
	}
	if req.DepIata != "" {
		params.Set("dep_iata", req.DepIata)
	}
	if req.DepIcao != "" {
		params.Set("dep_icao", req.DepIcao)
	}
	if req.ArrIata != "" {
		params.Set("arr_iata", req.ArrIata)
	}
	if req.ArrIcao != "" {
		params.Set("arr_icao", req.ArrIcao)
	}
	if req.AirlineIcao != "" {
		params.Set("airline_icao", req.AirlineIcao)
	}
	if req.AirlineIata != "" {
		params.Set("airline_iata", req.AirlineIata)
	}
	params.Set("api_key", h.ApiKey)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Printf("request: %v", urlPath)
	resp, err := http.Get(urlPath)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res scheduleResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res.Response, nil
}
