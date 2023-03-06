package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"xiaolong.ji.com/airport/airport-service/module/airports"
	"xiaolong.ji.com/airport/airport-service/module/common"
	"xiaolong.ji.com/airport/airport-service/module/schedules"
	"xiaolong.ji.com/airport/pkg/config"
)

type TodayRequest struct {
	CountryCode string `json:"country_code"`
	ArrDate     string `json:"arr_date"`
	DepDate     string `json:"dep_date"`
	PageNum     int    `json:"page_num"`
	PageSize    int    `json:"page_size"`
	OutputPath  string `json:"output_path"`
}

var getCountryCode = func(airportMap map[string]airports.Airport, airportIcalCode string) string {
	if r, ok := airportMap[airportIcalCode]; ok {
		return r.CountryCode
	}
	return ""
}
var getCityName = func(airportMap map[string]airports.Airport, airportIcalCode string) string {
	if r, ok := airportMap[airportIcalCode]; ok {
		return r.City
	}
	return ""
}
var getAirportName = func(airportMap map[string]airports.Airport, airportIcalCode string) string {
	if r, ok := airportMap[airportIcalCode]; ok {
		return r.Name
	}
	return ""
}

func writeScheduleRecordsSummaryToCSV(csvfile string, airportMap map[string]airports.Airport, records map[string][]schedules.Schedule, isArrive bool) error {
	file, err := os.OpenFile(csvfile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("open file %v failed", csvfile)
		return err
	}
	defer file.Close()

	var getHourRow = func() []string {
		h := []string{""}
		for i := 0; i <= 23; i++ {
			t := strconv.Itoa(i)
			if len(t) == 1 {
				t = "0" + t
			}
			h = append(h, t)
		}
		return h
	}

	w := csv.NewWriter(file)
	header := getHourRow()
	w.Write(header)

	summary := make(map[string]map[string]int) //[date][hour]count
	for _, recordsForOneAirport := range records {
		for _, record := range recordsForOneAirport {
			date, hour := "", ""
			if isArrive {
				date, _ = common.Date(record.ArrTime)
				hour, _ = common.Hour(record.ArrTime)
			} else {
				date, _ = common.Date(record.DepTime)
				hour, _ = common.Hour(record.DepTime)
			}
			if _, ok := summary[date]; !ok {
				summary[date] = make(map[string]int)
			}
			if _, ok := summary[date][hour]; !ok {
				summary[date][hour] = 0
			}
			summary[date][hour] = summary[date][hour] + 1
		}
	}

	for date, _ := range summary {
		row := make([]string, 0)
		row = append(row, date)
		for _, hkey := range header {
			if hkey == "" {
				continue
			}
			countStr := strconv.Itoa(summary[date][hkey])
			row = append(row, countStr)
		}
		w.Write(row)
	}
	w.Flush()
	return nil
}

func writeScheduleRecordsToCSV(csvfile string, airportMap map[string]airports.Airport, records map[string][]schedules.Schedule, isArrive bool) error {
	file, err := os.OpenFile(csvfile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("open file failed", err)
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	if isArrive {
		header := []string{"CountryCode", "City", "Airport", "AirportIcaoCode", "Flight", "timestamp", "ArriveTimeLocal", "ArriveTimeUTC", "FlightStatus"}
		w.Write(header)
		for key, val := range records {
			for _, record := range val {
				row := []string{
					getCountryCode(airportMap, key),
					getCityName(airportMap, key),
					getAirportName(airportMap, key),
					key,
					record.FlightIata,
					strconv.Itoa(int(record.ArrTimeTs)),
					record.ArrTime,
					record.ArrTimeUtc,
					record.Status,
				}
				w.Write(row)
			}
		}
		w.Flush()
		log.Printf("output file: %v", csvfile)
	} else {
		header := []string{"CountryCode", "City", "Airport", "AirportIcaoCode", "Flight", "timestamp", "DepartureTimeLocal", "DepartureTimeUTC", "FlightStatus"}
		w.Write(header)
		for key, val := range records {
			for _, record := range val {
				row := []string{
					getCountryCode(airportMap, key),
					getCityName(airportMap, key),
					getAirportName(airportMap, key),
					key,
					record.FlightIata,
					strconv.Itoa(int(record.DepTimeTs)),
					record.DepTime,
					record.DepTimeUtc,
					record.Status,
				}
				w.Write(row)
			}
		}
		w.Flush()
		log.Printf("output file: %v", csvfile)
	}
	return nil
}

func GetSchedules(c *gin.Context) {
	r := &TodayRequest{}
	err := c.ShouldBind(r)
	log.Printf("%v,%v", r, err)
	if err != nil {
		errmsg := fmt.Sprintf("failed to parse request, error:%v", err.Error())
		log.Printf(errmsg)
		c.JSON(400, gin.H{
			"code":  -1,
			"error": errmsg,
		})
		return
	}

	//get all airports for the country
	airportHttpClient := airports.NewAirportHttpClient(config.GetAirLabConfig().ApiKey)
	resp, err := airportHttpClient.GET(airports.AirportRequest{
		CountryCode: r.CountryCode,
	})
	if err != nil {
		errmsg := fmt.Sprintf("failed to request airlabs, error:%v", err.Error())
		log.Printf(errmsg)
		c.JSON(500, gin.H{
			"code":  -1,
			"error": errmsg,
		})
		return
	}

	//get all schedules for the country
	airportMap := make(map[string]airports.Airport)
	arrMap := make(map[string][]schedules.Schedule)
	depMap := make(map[string][]schedules.Schedule)
	for _, airport := range resp {
		icao := airport.IcaoCode
		airportMap[icao] = airport
		scheduleHttpCLient := schedules.NewScheduleHttpClient(config.GetAirLabConfig().ApiKey)
		//arrive
		arrResp, err := scheduleHttpCLient.GET(schedules.ScheduleRequest{
			ArrIcao:  icao,
			PageNum:  r.PageNum,
			PageSize: r.PageSize,
			ArrDate:  r.ArrDate,
		})
		if err == nil {
			arrMap[icao] = arrResp
		} else {
			log.Printf("%v", err)
		}

		//departure
		depResp, err := scheduleHttpCLient.GET(schedules.ScheduleRequest{
			DepIcao:  icao,
			PageNum:  r.PageNum,
			PageSize: r.PageSize,
			DepDate:  r.DepDate,
		})
		if err == nil {
			depMap[icao] = depResp
		} else {
			log.Printf("%v", err)
		}
	}

	outputpath := r.OutputPath
	if len(outputpath) == 0 {
		airportJson, _ := json.Marshal(airportMap)
		arriveJson, _ := json.Marshal(arrMap)
		departureJson, _ := json.Marshal(depMap)
		c.JSON(200, gin.H{
			"code":      0,
			"airport":   string(airportJson),
			"arrive":    string(arriveJson),
			"departure": string(departureJson),
		})
	} else {
		//write to csv
		csvfile1 := outputpath + "/" + r.CountryCode + "-arrive-schedules.csv"
		_ = writeScheduleRecordsToCSV(csvfile1, airportMap, arrMap, true)
		csvfile2 := outputpath + "/" + r.CountryCode + "-departure-schedules.csv"
		_ = writeScheduleRecordsToCSV(csvfile2, airportMap, depMap, false)
		csvfile3 := outputpath + "/" + r.CountryCode + "-arrive-schedules-summary.csv"
		_ = writeScheduleRecordsSummaryToCSV(csvfile3, airportMap, arrMap, true)
		csvfile4 := outputpath + "/" + r.CountryCode + "-departure-schedules-summary.csv"
		_ = writeScheduleRecordsSummaryToCSV(csvfile4, airportMap, depMap, false)
		c.JSON(200, gin.H{
			"code": 0,
			"output": gin.H{
				"arrive_file":            csvfile1,
				"departure_file":         csvfile2,
				"arrive_summary_file":    csvfile3,
				"departure_summary_file": csvfile4,
			},
		})
	}

	return
}
