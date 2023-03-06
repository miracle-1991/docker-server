package schedules

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"xiaolong.ji.com/airport/airport-service/module/common"
)

type ScheduleHttpClient interface {
	GET(ScheduleRequest) ([]Schedule, error)
	UPDATE(ScheduleRequest) ([]Schedule, error)
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
	Response []Schedule `json:"response"`
}

func (h *httpClient) GET(req ScheduleRequest) ([]Schedule, error) {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	//check if exist
	if req.ArrDate != "" {
		exist, err := ExistScheduleByArrDate(req.ArrDate)
		if err == nil && exist == false {
			return nil, errors.New("no data")
		}
	}

	if req.DepDate != "" {
		exist, err := ExistScheduleByDepDate(req.DepDate)
		if err == nil && exist == false {
			return nil, errors.New("no data")
		}
	}

	if req.DepIata != "" {
		maps["dep_iata"] = req.DepIata
		maps["dep_date"] = req.DepDate
	} else if req.DepIcao != "" {
		maps["dep_icao"] = req.DepIcao
		maps["dep_date"] = req.DepDate
	} else if req.ArrIata != "" {
		maps["arr_iata"] = req.ArrIata
		maps["arr_date"] = req.ArrDate
	} else if req.ArrIcao != "" {
		maps["arr_icao"] = req.ArrIcao
		maps["arr_date"] = req.ArrDate
	} else if req.AirlineIcao != "" {
		maps["airline_icao"] = req.AirlineIcao
	} else if req.AirlineIata != "" {
		maps["airline_iata"] = req.AirlineIata
	}
	schedules, err := GetSchedules(req.PageNum, req.PageSize, maps)
	if err != nil {
		return nil, err
	}
	if req.ArrIata != "" || req.ArrIcao != "" {
		//sort and unique
		return UniqueByArrTime(schedules), nil
	} else if req.DepIata != "" || req.DepIcao != "" {
		return UniqueByDepTime(schedules), nil
	}
	return schedules, nil
}

func (h *httpClient) UPDATE(req ScheduleRequest) ([]Schedule, error) {
	params := url.Values{}
	Url, err := url.Parse(scheduleUrl)
	if err != nil {
		return nil, err
	}
	params.Set("api_key", h.ApiKey)
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
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Printf("request: %v\n", urlPath)
	resp, err := http.Get(urlPath)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res scheduleResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	schedules := res.Response
	newSchedules := make([]Schedule, 0)
	for _, s := range schedules {
		var exist bool
		if s.CSFlightIata != "" {
			exist, err = ExistSchedule(s.CSFlightIata, s.ArrTimeTs, s.DepTimeTs)
		} else {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !exist {
			s.ArrDate, _ = common.Date(s.ArrTime)
			s.DepDate, _ = common.Date(s.DepTime)
			err = AddSchedule(s)
			if err != nil {
				return nil, err
			}
			newSchedules = append(newSchedules, s)
		}
	}
	return newSchedules, nil
}
