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

func GetDriverGPS(c *gin.Context) {
	_ = c.Request.ParseForm()
	log.Println(c.Request.PostForm)

	driverid := c.PostForm("driverid")
	id, err := strconv.ParseInt(driverid, 10, 64)
	if err != nil {
		errmsg := fmt.Sprintf("failed to parse the driverid:%v, error:%v\n", driverid, err.Error())
		log.Printf(errmsg)
		c.JSON(400, gin.H{
			"error": errmsg,
		})
		return
	}

	timestr := c.PostForm("time")
	objTime, err := timeformat.Parsetime(timestr)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	forward := c.PostForm("forward")
	iBefore, err := strconv.Atoi(forward)
	if err != nil {
		errmsg := fmt.Sprintf("failed to parse the forward:%v, error:%v\n", forward, err.Error())
		log.Printf(errmsg)
		c.JSON(400, gin.H{
			"error": errmsg,
		})
		return
	}

	backward := c.PostForm("backward")
	iAfter, err := strconv.Atoi(backward)
	if err != nil {
		errmsg := fmt.Sprintf("failed to parse the backward:%v, error:%v\n", backward, err.Error())
		log.Printf(errmsg)
		c.JSON(400, gin.H{
			"error": errmsg,
		})
		return
	}

	outputpath := c.PostForm("outputpath")
	if len(outputpath) == 0 {
		errmsg := fmt.Sprintf("you must set output path\n")
		log.Printf(errmsg)
		c.JSON(400, gin.H{
			"error": errmsg,
		})
		return
	}

	var records []DBSampled516.GpsPingSimpled516
	startTime := objTime.Add(-1 * time.Duration(iBefore) * time.Minute)
	endTime := objTime.Add(time.Duration(iAfter) * time.Minute)
	log.Printf("start pulling data, startTime: %s, endTime: %s ......\n", startTime.String(), endTime.String())
	log.Printf("start pulling data, startTime: %v, endTime: %v ......\n", startTime.Unix(), endTime.Unix())
	records = DBSampled516.GetDasGpsInMinuteInRange(id, startTime, endTime)
	if len(records) == 0 {
		log.Printf("get nothing from db\n")
		return
	}
	if outputpath[len(outputpath)-1] == '/' {
		outputpath = strings.TrimRight(outputpath, "/")
	}
	fileName := driverid + "-" + strings.Replace(strings.Replace(timestr, " ", "-", -1), ":", "-", -1)
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
