package today

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"xiaolong.ji.com/airport/airport-service/module/airports"
	"xiaolong.ji.com/airport/airport-service/module/schedules"
	"xiaolong.ji.com/airport/pkg/config"
)

type TodayRequest struct {
	CountryCode string `json:"country_code"`
	OutputPath  string `json:"output_path"`
}

func writeScheduleRecordsToCSV(csvfile string, airportMap map[string]airports.AirportResponse, records map[string][]schedules.ScheduleResponse, isArrive bool) error {
	file, err := os.OpenFile(csvfile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("open file failed", err)
		return err
	}
	defer file.Close()

	var getCountryCode = func(airportMap map[string]airports.AirportResponse, airportIcalCode string) string {
		if r, ok := airportMap[airportIcalCode]; ok {
			return r.CountryCode
		}
		return ""
	}
	var getCityName = func(airportMap map[string]airports.AirportResponse, airportIcalCode string) string {
		if r, ok := airportMap[airportIcalCode]; ok {
			return r.City
		}
		return ""
	}
	var getAirportName = func(airportMap map[string]airports.AirportResponse, airportIcalCode string) string {
		if r, ok := airportMap[airportIcalCode]; ok {
			return r.Name
		}
		return ""
	}

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

func GetSchedulesForToday(c *gin.Context) {
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
	outputpath := r.OutputPath
	if len(outputpath) == 0 {
		errmsg := fmt.Sprintf("you must set output path\n")
		log.Printf(errmsg)
		c.JSON(400, gin.H{
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
	airportMap := make(map[string]airports.AirportResponse)
	arrMap := make(map[string][]schedules.ScheduleResponse)
	depMap := make(map[string][]schedules.ScheduleResponse)
	for _, airport := range resp {
		icao := airport.IcaoCode
		airportMap[icao] = airport
		scheduleHttpCLient := schedules.NewScheduleHttpClient(config.GetAirLabConfig().ApiKey)
		//arrive
		arrResp, err := scheduleHttpCLient.GET(schedules.ScheduleRequest{
			ArrIcao: icao,
		})
		if err == nil {
			arrMap[icao] = arrResp
		}

		//departure
		depResp, err := scheduleHttpCLient.GET(schedules.ScheduleRequest{
			DepIcao: icao,
		})
		if err == nil {
			depMap[icao] = depResp
		}
	}
	//write to csv
	csvfile1 := outputpath + "/" + r.CountryCode + "-arrive-schedules.csv"
	_ = writeScheduleRecordsToCSV(csvfile1, airportMap, arrMap, true)
	csvfile2 := outputpath + "/" + r.CountryCode + "-departure-schedules.csv"
	_ = writeScheduleRecordsToCSV(csvfile2, airportMap, depMap, false)
	c.JSON(200, gin.H{
		"code": 0,
		"output": gin.H{
			"file1": csvfile1,
			"file2": csvfile2,
		},
	})
	return
}
