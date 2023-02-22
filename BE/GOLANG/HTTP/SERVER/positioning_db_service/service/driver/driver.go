package driver

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

func getYear(timestr string) (string, error) {
	t, e := time.Parse(timeFormat, timestr)
	if e != nil {
		gin.Logger()
	}
}

func GetDriverGPS(c *gin.Context) {
	driverid := c.PostForm("id")
	starttime := c.PostForm("starttime")
	endtime := c.PostForm("endtime")

	id, err := strconv.ParseInt(driverID, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	if len(year) < 4 {
		log.Printf("year should be like 2021、2022....")
		return
	}
	if len(month) < 2 {
		log.Printf("year should be like 01、02....")
		return
	}
	if len(day) < 2 {
		log.Printf("day should be like 01、02....")
		return
	}
	if len(hour) < 2 {
		log.Printf("day should be like 09、15....")
		return
	}
	if minute != "" && len(minute) < 2 {
		log.Printf("minute must be like 00、01、60")
		return
	}
	var records []DBSampled516.GpsPingSimpled516

	objTimeStr := year + "-" + month + "-" + day + " " + hour + ":" + minute + ":00"
	objTime, err := time.Parse(timeFormat, objTimeStr)
	if err != nil {
		log.Printf("failed to parse time:%v", objTimeStr)
		return
	}
	iBefore, err := strconv.Atoi(intervalBefore)
	if err != nil {
		log.Printf("failed to parse the interval:%v", intervalBefore)
		return
	}
	iAfter, err := strconv.Atoi(intervalAfter)
	if err != nil {
		log.Printf("failed to parse the interval:%v", intervalAfter)
		return
	}
	startTime := objTime.Add(-1 * time.Duration(iBefore) * time.Minute)
	endTime := objTime.Add(time.Duration(iAfter) * time.Minute)
	log.Printf("start pulling data, startTime: %s, endTime: %s ......\n", startTime.String(), endTime.String())
	log.Printf("start pulling data, startTime: %v, endTime: %v ......\n", startTime.Unix(), endTime.Unix())
	records = DBSampled516.GetDasGpsInMinuteInRange(id, startTime, endTime)
	if len(records) == 0 {
		log.Printf("get nothing from db\n")
		return
	}
	var fileName string
	if minute == "" {
		fileName = driverID + "-" + year + month + day + hour
	} else {
		fileName = driverID + "-" + year + month + day + hour + minute
	}
	err = gpsTemplate.InitTempalte()
	if err != nil {
		log.Printf("failed to init template")
		return
	}
	err = gpsTemplate.TemplateRawGpsToHtml(fileName, records, time.Time{}, time.Time{})
	if err != nil {
		log.Fatal(err)
	}
	err = gpsTemplate.TemplateSnapedGpsToHtml(fileName, records, time.Time{}, time.Time{})
	if err != nil {
		log.Fatal(err)
	}
	err = gpsTemplate.WriteToCsv(fileName, records)
	if err != nil {
		log.Fatal(err)
	}
}
