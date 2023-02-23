package booking

import (
	"example.com/db_service/module/DBSampled516"
	"example.com/db_service/pkg/gpsTemplate"
	"example.com/db_service/pkg/timeformat"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

type BookingInfo struct {
	BookingCode string `json:"bookingcode"`
	Time        string `json:"time"`
	OutputPath  string `json:"outputpath"`
}

func GetBookingGPS(c *gin.Context) {
	booking := &BookingInfo{}
	err := c.ShouldBind(booking)
	if err != nil {
		errmsg := fmt.Sprintf("failed to parse, error:%v\n", err.Error())
		log.Printf(errmsg)
		c.JSON(400, gin.H{
			"error": errmsg,
		})
		return
	}

	objTime, err := timeformat.ParseDate(booking.Time)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	records := DBSampled516.GetGpsOfOneBook(booking.BookingCode, objTime)
	if len(records) == 0 {
		errmsg := fmt.Sprintf("get nothing for %v at %v\n", booking.BookingCode, booking.Time)
		c.JSON(400, gin.H{
			"error": errmsg,
		})
		return
	}

	outputpath := booking.OutputPath
	if outputpath[len(outputpath)-1] == '/' {
		outputpath = strings.TrimRight(outputpath, "/")
	}

	var fileName string
	fileName = booking.BookingCode + "-" + strings.Replace(strings.Replace(booking.Time, " ", "-", -1), ":", "-", -1)
	err = gpsTemplate.InitTempalte(outputpath)
	if err != nil {
		errmsg := fmt.Sprintf("failed to init template, error:%v\n", err.Error())
		log.Printf(errmsg)
		c.JSON(500, gin.H{
			"error": errmsg,
		})
		return
	}
	rawHtmlFile, err := gpsTemplate.TemplateRawGpsToHtml(fileName, records, time.Time{}, time.Time{})
	if err != nil {
		errmsg := fmt.Sprintf("failed to template raw html for %v, error:%v\n", fileName, err.Error())
		log.Printf(errmsg)
		c.JSON(500, gin.H{
			"error": errmsg,
		})
		return
	}
	snaphtmlfile, err := gpsTemplate.TemplateSnapedGpsToHtml(fileName, records, time.Time{}, time.Time{})
	if err != nil {
		errmsg := fmt.Sprintf("failed to template snap html for %v, error:%v\n", fileName, err.Error())
		log.Printf(errmsg)
		c.JSON(500, gin.H{
			"error": errmsg,
		})
		return
	}
	gpscsvfile, err := gpsTemplate.WriteToCsv(fileName, records)
	if err != nil {
		errmsg := fmt.Sprintf("failed to template gps csv for %v, error:%v\n", fileName, err.Error())
		log.Printf(errmsg)
		c.JSON(500, gin.H{
			"error": errmsg,
		})
		return
	}
	c.JSON(200, gin.H{
		"error":                    "",
		"TrajectoryOfOriginalGps":  rawHtmlFile,
		"TrajectoryOfProcessedGps": snaphtmlfile,
		"OriginalGpsInCSV":         gpscsvfile,
	})
	return
}
