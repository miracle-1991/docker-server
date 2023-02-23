package driver

import (
	"example.com/db_service/module/DBSampled516"
	"example.com/db_service/pkg/gpsTemplate"
	"example.com/db_service/pkg/timeformat"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
	"time"
)

type DriverInfo struct {
	DriverId   int64  `json:"driverid"`
	Time       string `json:"time"`
	Forward    int    `json:"forward"`
	Backward   int    `json:"backward"`
	Outputpath string `json:"outputpath"`
}

func GetDriverGPS(c *gin.Context) {
	dax := &DriverInfo{}
	err := c.ShouldBind(dax)
	if err != nil {
		errmsg := fmt.Sprintf("failed to parse, error:%v\n", err.Error())
		log.Printf(errmsg)
		c.JSON(400, gin.H{
			"error": errmsg,
		})
		return
	}

	driverid := dax.DriverId
	objTime, err := timeformat.Parsetime(dax.Time)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	outputpath := dax.Outputpath
	if len(outputpath) == 0 {
		errmsg := fmt.Sprintf("you must set output path\n")
		log.Printf(errmsg)
		c.JSON(400, gin.H{
			"error": errmsg,
		})
		return
	}

	var records []DBSampled516.GpsPingSimpled516
	startTime := objTime.Add(-1 * time.Duration(dax.Forward) * time.Minute)
	endTime := objTime.Add(time.Duration(dax.Backward) * time.Minute)
	log.Printf("start pulling data, startTime: %s, endTime: %s ......\n", startTime.String(), endTime.String())
	log.Printf("start pulling data, startTime: %v, endTime: %v ......\n", startTime.Unix(), endTime.Unix())
	records = DBSampled516.GetDasGpsInMinuteInRange(dax.DriverId, startTime, endTime)
	if len(records) == 0 {
		log.Printf("get nothing from db\n")
		return
	}
	if outputpath[len(outputpath)-1] == '/' {
		outputpath = strings.TrimRight(outputpath, "/")
	}
	fileName := strconv.Itoa(int(driverid)) + "-" + strings.Replace(strings.Replace(dax.Time, " ", "-", -1), ":", "-", -1)
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

	snapHtmlFile, err := gpsTemplate.TemplateSnapedGpsToHtml(fileName, records, time.Time{}, time.Time{})
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
		"TrajectoryOfProcessedGps": snapHtmlFile,
		"DriverOriginalGpsInCSV":   gpscsvfile,
	})
	return
}
