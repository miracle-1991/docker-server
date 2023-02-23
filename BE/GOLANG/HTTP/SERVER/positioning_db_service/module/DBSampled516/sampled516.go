package DBSampled516

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GpsPingSimpled516 struct {
	StreamTime         string  `gorm:"column:streamtime" json:"streamtime"`
	DriverId           int64   `gorm:"column:driverid" json:"driverid"`
	MapId              string  `gorm:"column:mapid" json:"mapid"`
	MapVersion         string  `gorm:"column:mapversion" json:"mapversion"`
	VehicleId          int64   `gorm:"column:vehicleid" json:"vehicleid"`
	Latitude           float64 `gorm:"column:latitude" json:"latitude"`
	Longitude          float64 `gorm:"column:longitude" json:"longitude"`
	Accuracy           float64 `gorm:"column:accuracy" json:"accuracy"`
	Bearing            float64 `gorm:"column:bearing" json:"bearing"`
	Timestamp          int64   `gorm:"column:timestamp" json:"timestamp"`
	Altitude           float64 `gorm:"column:altitude" json:"altitude"`
	Speed              float64 `gorm:"column:speed" json:"speed"`
	SegmentstartNode   int64   `gorm:"column:segmentstartnode" json:"segmentstartnode"`
	SegmentendNode     int64   `gorm:"column:segmentendnode" json:"segmentendnode"`
	SegmentoffsetRatio float64 `gorm:"column:segmentoffsetratio" json:"segmentoffsetratio"`
	WayIds             string  `gorm:"column:wayids" json:"wayids"`
	Availability       string  `gorm:"column:availability" json:"availability"`
	PingStatus         int64   `gorm:"column:pingstatus" json:"pingstatus"`
	Projectedlat       float64 `gorm:"column:projectedlat" json:"projectedlat"`
	Projectedlng       float64 `gorm:"column:projectedlng" json:"projectedlng"`
	SegmentBearing     float64 `gorm:"column:segmentbearing" json:"segmentbearing"`
	Filter             string  `gorm:"column:filter" json:"filter"`
	DriverState        string  `gorm:"column:driverstate" json:"'driverstate'"`
	CityId             string  `gorm:"column:cityid" json:"cityid"`
	ChannelId          string  `gorm:"column:channelid" json:"channelid"`
	Message            string  `gorm:"column:message" json:"message"`
	BookingCode        string  `gorm:"column:bookingcode" json:"bookingcode"`
	Confidence         float64 `gorm:"column:confidence" json:"confidence"`
	Source             string  `gorm:"column:source" json:"source"`
	StaleDuration      int64   `gorm:"column:staleduration" json:"staleduration"'`
	Transformed        bool    `gorm:"column:transformed" json:"transformed"` //true: convert from 516
	IsForward          bool    `gorm:"column:isforward" json:"isforward"`
	DeviceInfo         string  `gorm:"column:deviceinfo" json:"deviceinfo"`
	GeoHash            string  `gorm:"column:geohash" json:"geohash"`
	Altlatitude        float64 `gorm:"column:altlatitude" json:"altlatitude"`
	Altlongitude       float64 `gorm:"column:altlongitude" json:"altlongitude"`
	Alttimestamp       int64   `gorm:"column:alttimestamp" json:"alttimestamp"`
	Altbearing         float64 `gorm:"column:altbearing" json:"altbearing"`
	Altspeed           float64 `gorm:"column:altspeed" json:"altspeed"`
	Altaccuracy        float64 `gorm:"column:altaccuracy" json:"altaccuracy"`
	Altsource          string  `gorm:"column:altsource" json:"altsource"`
	Minute             int     `gorm:"column:minute" json:"minute"`
	Year               string  `gorm:"column:year" json:"year"`
	Month              string  `gorm:"column:month" json:"month"`
	Day                string  `gorm:"column:day" json:"day"`
	Hour               string  `gorm:"column:hour" json:"hour"`
}

func (GpsPingSimpled516) TableName() string {
	return "driver_snapped_location"
}

type ByTime []GpsPingSimpled516

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Less(i, j int) bool { return a[i].Timestamp <= a[j].Timestamp }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

const TimeFormat = "2006-01-02 15:04:05"

func (u GpsPingSimpled516) ToString() []string {
	return []string{
		u.StreamTime,
		strconv.FormatInt(u.DriverId, 10),
		u.BookingCode,
		fmt.Sprintf("%f", u.Latitude),
		fmt.Sprintf("%f", u.Longitude),
		fmt.Sprintf("%f", u.Projectedlat),
		fmt.Sprintf("%f", u.Projectedlng),
		strconv.FormatInt(u.StaleDuration, 10),
		fmt.Sprintf("%f", u.Speed),
		strconv.FormatInt(u.Timestamp, 10),
		u.Filter,
		u.Hour,
		strconv.Itoa(u.Minute),
	}
}

type sqlTime struct {
	year   string
	month  string
	day    string
	hour   string
	minute int
}

func parseFromTime(t time.Time) sqlTime {
	var s sqlTime
	tf := t.Format(TimeFormat)
	tfSplit := strings.Split(tf, " ")
	front := strings.Split(tfSplit[0], "-")
	back := strings.Split(tfSplit[1], ":")
	s.year, s.month, s.day, s.hour = front[0], front[1], front[2], back[0]
	s.minute, _ = strconv.Atoi(back[1])
	return s
}

func getSql(startTime, endTime time.Time, driverId int64) []string {
	var s, e sqlTime = parseFromTime(startTime), parseFromTime(endTime)
	selectSql := "select streamtime, driverid, bookingcode, latitude, longitude, projectedlat, projectedlng, staleduration, speed, timestamp, filter, hour, minute from grab_datastore.driver_snapped_location "
	//都在同一小时内
	whereSQL := []string{}
	if s.year == e.year && s.month == e.month && s.day == e.day && s.hour == e.hour {
		tempSql := selectSql + fmt.Sprintf("where year='%s' and month='%s' and day='%s' and hour='%s' and minute>=%d and minute<=%d and driverid=%d order by timestamp asc, streamtime asc",
			s.year, s.month, s.day, s.hour, s.minute, e.minute, driverId)
		whereSQL = append(whereSQL, tempSql)
		return whereSQL
	}

	//跨小时，但是不跨天
	tempSql1 := selectSql + fmt.Sprintf("where year='%s' and month='%s' and day='%s' and hour='%s' and minute>=%d and minute<=%d and driverid=%d order by timestamp asc, streamtime asc",
		s.year, s.month, s.day, s.hour, s.minute, 59, driverId)
	tempSql2 := selectSql + fmt.Sprintf("where year='%s' and month='%s' and day='%s' and hour='%s' and minute>=%d and minute<=%d and driverid=%d order by timestamp asc, streamtime asc",
		e.year, e.month, e.day, e.hour, 00, e.minute, driverId)
	whereSQL = append(whereSQL, tempSql1)
	whereSQL = append(whereSQL, tempSql2)
	return whereSQL
}

func getDasGpsInMinuteInRangeOnce(sqltext string) (pos []GpsPingSimpled516, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("getOnce error:%v\n", err)
			pos = make([]GpsPingSimpled516, 0)
			err = errors.New("JumpCloud authentication error, please check your JumpCloud username and password")
		}
	}()
	rows, err := db.Raw(sqltext).Rows()
	defer rows.Close()
	if err == gorm.ErrRecordNotFound {
		log.Println("find nothing")
		return pos, err
	}
	if err != nil {
		log.Printf("query error:%v\n", err.Error())
		return pos, err
	}

	for rows.Next() {
		var p GpsPingSimpled516
		rows.Scan(&p.StreamTime, &p.DriverId, &p.BookingCode, &p.Latitude, &p.Longitude, &p.Projectedlat, &p.Projectedlng,
			&p.StaleDuration, &p.Speed, &p.Timestamp, &p.Filter, &p.Hour, &p.Minute)
		pos = append(pos, p)
	}
	return pos, nil
}

func GetDasGpsInMinuteInRange(driverid int64, startTime, endTime time.Time) ([]GpsPingSimpled516, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("GetDasGpsInMinuteInRange error:%v\n", err)
		}
	}()

	var pos []GpsPingSimpled516
	sqlText := getSql(startTime, endTime, driverid)
	for i := 0; i < len(sqlText); i++ {
		tmpPos, err := getDasGpsInMinuteInRangeOnce(sqlText[i])
		if err != nil {
			return nil, err
		} else {
			if len(tmpPos) > 0 {
				pos = append(pos, tmpPos...)
			}
		}
	}

	return pos, nil
}

//GetGpsOfOneDay 获取某个司机在某一天内所有的gps
func GetGpsOfOneDay(driverid int64, daytime time.Time) []GpsPingSimpled516 {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("error:%v\n", err)
		}
	}()
	var pos []GpsPingSimpled516
	sqlTime := parseFromTime(daytime)
	sqlText := "select streamtime, driverid, bookingcode, latitude, longitude, timestamp, year, month, day, hour from grab_datastore.driver_snapped_location where year='%s' and month='%s' and day='%s' and driverid=%d"
	sqlQueryText := fmt.Sprintf(sqlText, sqlTime.year, sqlTime.month, sqlTime.day, driverid)
	rows, err := db.Raw(sqlQueryText).Rows()
	defer rows.Close()
	if err == gorm.ErrRecordNotFound {
		log.Println("find nothing")
		return pos
	}
	if err != nil {
		log.Printf("query error:%v\n", err)
		return pos
	}
	for rows.Next() {
		var p GpsPingSimpled516
		rows.Scan(&p.StreamTime, &p.DriverId, &p.BookingCode, &p.Latitude, &p.Longitude, &p.Timestamp, &p.Year, &p.Month, &p.Day, &p.Hour)
		pos = append(pos, p)
	}
	sort.Sort(ByTime(pos))
	return pos
}

func GetGpsOfOneBook(bookingcode string, daytime time.Time) []GpsPingSimpled516 {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("error:%v\n", err)
		}
	}()
	var pos []GpsPingSimpled516
	sqlTime := parseFromTime(daytime)
	sqlText := "select streamtime, driverid, bookingcode, latitude, longitude, timestamp, year, month, day, hour from grab_datastore.driver_snapped_location where year='%s' and month='%s' and day='%s' and bookingcode='%s'"
	sqlQueryText := fmt.Sprintf(sqlText, sqlTime.year, sqlTime.month, sqlTime.day, bookingcode)
	rows, err := db.Raw(sqlQueryText).Rows()
	if err == gorm.ErrRecordNotFound {
		log.Println("find nothing")
		return pos
	}
	if err != nil {
		log.Printf("query error:%v\n", err)
		return pos
	}
	for rows.Next() {
		var p GpsPingSimpled516
		rows.Scan(&p.StreamTime, &p.DriverId, &p.BookingCode, &p.Latitude, &p.Longitude, &p.Timestamp, &p.Year, &p.Month, &p.Day, &p.Hour)
		pos = append(pos, p)
	}
	sort.Sort(ByTime(pos))
	return pos
}
