package schedules

import (
	"github.com/jinzhu/gorm"
	"sort"
	"xiaolong.ji.com/airport/airport-service/module/common"
)

type ScheduleRequest struct {
	PageNum     int    `json:"page_num"`
	PageSize    int    `json:"page_size"`
	DepIata     string `json:"dep_iata"`     //Filtering by departure Airport IATA code.
	DepIcao     string `json:"dep_icao"`     //Filtering by departure Airport ICAO code.
	DepDate     string `json:"dep_date"`     //departure date
	ArrIata     string `json:"arr_iata"`     //Filtering by arrival Airport IATA code.
	ArrIcao     string `json:"arr_icao"`     //Filtering by arrival Airport ICAO code.
	ArrDate     string `json:"arr_date"`     //arrival date
	AirlineIcao string `json:"airline_icao"` //Filtering by Airline ICAO code.
	AirlineIata string `json:"airline_iata"` //Filtering by Airline IATA code.
}

type Schedule struct {
	common.Model
	AirlineIata     string `json:"airline_iata"`      //Airline IATA code.
	AirlineIcao     string `json:"airline_icao"`      //Airline ICAO code.
	FlightIata      string `json:"flight_iata"`       //Flight IATA code-number
	FlightIcao      string `json:"flight_icao"`       //Flight ICAO code-number.
	FlightNumber    string `json:"flight_number"`     //Flight number only.
	CSAirlineIata   string `json:"cs_airline_iata"`   //Codeshared airline IATA code.
	CSFlightIata    string `json:"cs_flight_iata"`    //Codeshared flight IATA code-number.
	CSFlightNumber  string `json:"cs_flight_number"`  //Codeshared flight number.
	DepIata         string `json:"dep_iata"`          //Departure airport IATA code.
	DepIcao         string `json:"dep_icao"`          //Departure airport ICAO code.
	DepTerminal     string `json:"dep_terminal"`      //Estimated departure terminal.
	DepGate         string `json:"dep_gate"`          //Estimated departure gate.
	DepTime         string `json:"dep_time"`          //Departure time in the airport time zone.
	DepDate         string `json:"dep_date"`          //Departure date in the airport time zone.
	DepTimeTs       uint32 `json:"dep_time_ts"`       //Departure UNIX timestamp.
	DepTimeUtc      string `json:"dep_time_utc"`      //Departure time in UTC time zone.
	DepEstimated    string `json:"dep_estimated"`     //Updated departure time in the airport time zone.
	DepEstimatedTs  uint32 `json:"dep_estimated_ts"`  //Updated departure UNIX timestamp.
	DepEstimatedUtc string `json:"dep_estimated_utc"` //Updated departure time in UTC time zone.
	DepActual       string `json:"dep_actual"`        //Actual departure time in the airport time zone.
	DepActualTs     uint32 `json:"dep_actual_ts"`     //Actual departure UNIX timestamp.
	DepActualUtc    string `json:"dep_actual_utc"`    //Actual departure time in UTC time zone.
	ArrIata         string `json:"arr_iata"`          //Arrival airport IATA code.
	ArrIcao         string `json:"arr_icao"`          //Arrival airport ICAO code.
	ArrTerminal     string `json:"arr_terminal"`      //Estimated arrival terminal.
	ArrGate         string `json:"arr_gate"`          //Estimated arrival gate.
	ArrBaggage      string `json:"arr_baggage"`       //Arrival baggage claim carousel number.
	ArrTime         string `json:"arr_time"`          //Arrival time in the airport time zone.
	ArrDate         string `json:"arr_date"`          //Arrival date in the airport time zone.
	ArrTimeTs       uint32 `json:"arr_time_ts"`       //Arrival UNIX timestamp.
	ArrTimeUtc      string `json:"arr_time_utc"`      //Arrival time in UTC time zone.
	ArrEstimated    string `json:"arr_estimated"`     //Updated arrival time in the airport time zone.
	ArrEstimatedTs  uint32 `json:"arr_estimated_ts"`  //Updated arrival UNIX timestamp.
	ArrEstimatedUtc string `json:"arr_estimated_utc"` //Updated arrival time in UTC time zone.
	ArrActual       string `json:"arr_actual"`        //Actual arrival time in the airport time zone.
	ArrActualTs     uint32 `json:"arr_actual_ts"`     //Actual arrival UNIX timestamp.
	ArrActualUtc    string `json:"arr_actual_utc"`    //Actual arrival time in UTC time zone.
	Duration        int    `json:"duration"`          //Estimated flight time (in minutes).
	Delayed         int    `json:"delayed"`           //(deprecated) Estimated flight delay time (in minutes).
	DepDelayed      int    `json:"dep_delayed"`       //Estimated time of flight departure delay (in minutes).
	ArrDelayed      int    `json:"arr_delayed"`       //Estimated time of flight arrival delay (in minutes).
	Status          string `json:"status"`            //Flight status - scheduled, cancelled, active, landed.
}

func (a Schedule) TableName() string {
	return "schedules"
}

//self def sort func: by arr time
type ByArrTime []Schedule

func (a ByArrTime) Len() int           { return len(a) }
func (a ByArrTime) Less(i, j int) bool { return a[i].ArrTimeTs < a[j].ArrTimeTs }
func (a ByArrTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func UniqueByArrTime(slist []Schedule) []Schedule {
	m := make(map[string]bool)
	tempSlist := make([]Schedule, 0)
	for _, s := range slist {
		key := s.CSFlightIata + s.ArrTime
		if m[key] == true {
			continue
		} else {
			m[key] = true
			tempSlist = append(tempSlist, s)
		}
	}
	sort.Sort(ByArrTime(tempSlist))
	return tempSlist
}

//self def sort func: by dep time
type ByDepTime []Schedule

func (a ByDepTime) Len() int           { return len(a) }
func (a ByDepTime) Less(i, j int) bool { return a[i].DepTimeTs < a[j].DepTimeTs }
func (a ByDepTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func UniqueByDepTime(slist []Schedule) []Schedule {
	m := make(map[string]bool)
	tempSlit := make([]Schedule, 0)
	for _, s := range slist {
		key := s.CSFlightIata + s.DepTime
		if m[key] == true {
			continue
		} else {
			m[key] = true
			tempSlit = append(tempSlit, s)
		}
	}
	sort.Sort(ByDepTime(tempSlit))
	return tempSlit
}

func ExistScheduleByCSFlightIata(cs_flight_iata string) (bool, error) {
	var item Schedule
	err := common.DB.Select("id").Where("cs_flight_iata = ? AND deleted_on = ? ", cs_flight_iata, 0).First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if item.ID > 0 {
		return true, nil
	}
	return false, nil
}

func ExistScheduleByDepDate(dep_date string) (bool, error) {
	var item Schedule
	err := common.DB.Select("id").Where("dep_date = ? AND deleted_on = ? ", dep_date, 0).First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if item.ID > 0 {
		return true, nil
	}
	return false, nil
}

func ExistScheduleByArrDate(arr_date string) (bool, error) {
	var item Schedule
	err := common.DB.Select("id").Where("arr_date = ? AND deleted_on = ? ", arr_date, 0).First(&item).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if item.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddSchedule(schedule Schedule) error {
	if err := common.DB.Create(&schedule).Error; err != nil {
		return err
	}
	return nil
}

func GetSchedules(pagenum int, pagesize int, maps interface{}) ([]Schedule, error) {
	var (
		schedules []Schedule
		err       error
	)
	if pagesize > 0 && pagenum > 0 {
		err = common.DB.Where(maps).Find(&schedules).Offset(pagenum).Limit(pagesize).Error
	} else {
		err = common.DB.Where(maps).Find(&schedules).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return schedules, nil
}

func GetSchedulesTotal(maps interface{}) (int, error) {
	var count int
	if err := common.DB.Model(&Schedule{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func ExistScheduleByID(id int) (bool, error) {
	var schedule Schedule
	err := common.DB.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&schedule).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if schedule.ID > 0 {
		return true, nil
	}
	return false, nil
}

func DelteSchedule(id int) error {
	if err := common.DB.Where("id = ?", id).Delete(&Schedule{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateSchedule(id int, data interface{}) error {
	if err := common.DB.Model(&Schedule{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
