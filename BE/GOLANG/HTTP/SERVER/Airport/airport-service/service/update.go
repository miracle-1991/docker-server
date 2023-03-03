package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"xiaolong.ji.com/airport/airport-service/module/airports"
	"xiaolong.ji.com/airport/airport-service/module/schedules"
	"xiaolong.ji.com/airport/pkg/config"
)

func UpdateSchedule(c *gin.Context) {
	countryCode := c.Param("country")
	if len(countryCode) == 0 {
		errmsg := fmt.Sprintf("failed to parse country code %v", countryCode)
		log.Printf(errmsg)
		c.JSON(400, gin.H{
			"code":  -1,
			"error": errmsg,
		})
		return
	}
	airportHttpClient := airports.NewAirportHttpClient(config.GetAirLabConfig().ApiKey)
	airportlist, err := airportHttpClient.UPDATE(airports.AirportRequest{
		CountryCode: countryCode,
	})
	if err != nil {
		errmsg := fmt.Sprintf("failed to update from airlabs, error:%v", err.Error())
		log.Printf(errmsg)
		c.JSON(500, gin.H{
			"code":  -1,
			"error": errmsg,
		})
		return
	}
	for _, air := range airportlist {
		log.Printf("update airport %v\n", air.Name)
		scheduleHttpCLient := schedules.NewScheduleHttpClient(config.GetAirLabConfig().ApiKey)
		// dep from this airport
		_, err = scheduleHttpCLient.UPDATE(schedules.ScheduleRequest{
			DepIata: air.IataCode,
			DepIcao: air.IcaoCode,
		})
		if err != nil {
			errmsg := fmt.Sprintf("failed to update departure from airport %v, error:%v", air.Name, err.Error())
			log.Printf(errmsg)
			c.JSON(500, gin.H{
				"code":  -1,
				"error": errmsg,
			})
			return
		}
		// arr to this airport
		_, err = scheduleHttpCLient.UPDATE(schedules.ScheduleRequest{
			ArrIata: air.IataCode,
			ArrIcao: air.IcaoCode,
		})
		if err != nil {
			errmsg := fmt.Sprintf("failed to update arrive to airport %v, error:%v", air.Name, err.Error())
			log.Printf(errmsg)
			c.JSON(500, gin.H{
				"code":  -1,
				"error": errmsg,
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"code":  0,
		"error": "success",
	})
	return
}
