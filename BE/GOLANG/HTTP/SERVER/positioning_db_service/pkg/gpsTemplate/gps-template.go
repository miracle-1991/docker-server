package gpsTemplate

import (
	"bufio"
	"encoding/csv"
	"example.com/db_service/module/DBSampled516"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"
)

var tplEngine *template.Template

const tplName = "gps-template.html"

var outPutPath string

func InitTempalte(outputpath string) error {
	if tplEngine != nil && outPutPath != "" {
		return nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := pwd + "/pkg/gpsTemplate/" + tplName
	tpl, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("parse error:%v\n", err)
		return err
	}
	tplEngine = template.Must(tpl, err)
	outPutPath = outputpath
	return nil
}

func TemplateRawGpsToHtml(outputFileName string, pings []DBSampled516.GpsPingSimpled516, removeTime, addTime time.Time) (string, error) {
	filename := outPutPath + "/" + outputFileName + "-raw.html"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("open file failed", err)
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	gpsPingData := ""
	gpsTimeData := ""
	gpsWhenQueRm := ""
	gpsWhenQueAdd := ""
	for i := 0; i < len(pings); i++ {
		tmpPingStr := fmt.Sprintf("%f", pings[i].Latitude) + "," + fmt.Sprintf("%f", pings[i].Longitude)
		gpsPingData = gpsPingData + tmpPingStr + " "
		tmpGpsTime := time.Unix(pings[i].Timestamp, 0).Add(-8 * time.Hour)
		tmpTimeStr := fmt.Sprintf("%s", tmpGpsTime.Format(DBSampled516.TimeFormat)) + "#"
		gpsTimeData = gpsTimeData + tmpTimeStr

		if tmpGpsTime.Hour() == removeTime.Hour() && tmpGpsTime.Minute() == removeTime.Minute() && gpsWhenQueRm == "" {
			gpsWhenQueRm = tmpPingStr
		}
		if tmpGpsTime.Hour() == addTime.Hour() && tmpGpsTime.Minute() == addTime.Minute() {
			gpsWhenQueAdd = tmpPingStr
		}
	}
	//log.Printf("removePoint:%v, addPoint:%v", gpsWhenQueRm, gpsWhenQueAdd)
	gpsPingData = strings.TrimRight(gpsPingData, " ")
	gpsTimeData = strings.TrimRight(gpsTimeData, "#")
	gpsWhenQueRm = strings.TrimRight(gpsWhenQueRm, " ")
	gpsWhenQueAdd = strings.TrimRight(gpsWhenQueAdd, " ")

	err = tplEngine.ExecuteTemplate(writer, tplName, map[string]interface{}{
		"GpsPings":        gpsPingData,
		"GpsTimes":        gpsTimeData,
		"GpsWhenQueueRm":  gpsWhenQueRm,
		"GpsWhenQueueAdd": gpsWhenQueAdd,
	})
	if err != nil {
		log.Printf("tempalte error:%v\n", err)
		return "", err
	}

	writer.Flush()
	log.Printf("output file: %v\n", filename)
	return filename, nil
}

func TemplateSnapedGpsToHtml(outputFileName string, pings []DBSampled516.GpsPingSimpled516, removeTime, addTime time.Time) (string, error) {
	filename := outPutPath + "/" + outputFileName + "-snaped.html"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("open file failed", err)
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	gpsPingData := ""
	gpsTimeData := ""
	gpsWhenQueRm := ""
	gpsWhenQueAdd := ""
	for i := 0; i < len(pings); i++ {
		tmpPingStr := fmt.Sprintf("%f", pings[i].Projectedlat) + "," + fmt.Sprintf("%f", pings[i].Projectedlng)
		gpsPingData = gpsPingData + tmpPingStr + " "
		tmpGpsTime := time.Unix(pings[i].Timestamp, 0).Add(-8 * time.Hour)
		tmpTimeStr := fmt.Sprintf("%s", tmpGpsTime.Format(DBSampled516.TimeFormat)) + "#"
		gpsTimeData = gpsTimeData + tmpTimeStr

		if tmpGpsTime.Hour() == removeTime.Hour() && tmpGpsTime.Minute() == removeTime.Minute() && gpsWhenQueRm == "" {
			gpsWhenQueRm = tmpPingStr
		}
		if tmpGpsTime.Hour() == addTime.Hour() && tmpGpsTime.Minute() == addTime.Minute() {
			gpsWhenQueAdd = tmpPingStr
		}
	}
	//log.Printf("removePoint:%v, addPoint:%v", gpsWhenQueRm, gpsWhenQueAdd)
	gpsPingData = strings.TrimRight(gpsPingData, " ")
	gpsTimeData = strings.TrimRight(gpsTimeData, "#")
	gpsWhenQueRm = strings.TrimRight(gpsWhenQueRm, " ")
	gpsWhenQueAdd = strings.TrimRight(gpsWhenQueAdd, " ")

	err = tplEngine.ExecuteTemplate(writer, tplName, map[string]interface{}{
		"GpsPings":        gpsPingData,
		"GpsTimes":        gpsTimeData,
		"GpsWhenQueueRm":  gpsWhenQueRm,
		"GpsWhenQueueAdd": gpsWhenQueAdd,
	})
	if err != nil {
		log.Printf("tempalte error:%v\n", err)
		return "", err
	}

	writer.Flush()
	log.Printf("output file: %v\n", filename)
	return filename, nil
}

func WriteToCsv(outputFileName string, pings []DBSampled516.GpsPingSimpled516) (string, error) {
	filename := outPutPath + "/" + outputFileName + ".csv"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("open file failed", err)
		return "", err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	header := []string{"streamtime", "driverid", "bookingcode", "latitude", "longitude", "projectedlat", "projectedlng", "staleduration", "speed", "timestamp", "filter", "hour", "minute"}
	w.Write(header)
	for i := 0; i < len(pings); i++ {
		err = w.Write(pings[i].ToString())
		if err != nil {
			return "", err
		}
	}
	w.Flush()
	log.Printf("output file: %v\n", filename)
	return filename, nil
}
