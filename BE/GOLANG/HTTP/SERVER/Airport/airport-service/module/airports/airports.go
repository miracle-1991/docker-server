package airports

import (
	"github.com/jinzhu/gorm"
	"xiaolong.ji.com/airport/airport-service/module/common"
)

type AirportRequest struct {
	PageNum     int    `json:"page_num"`
	PageSize    int    `json:"page_size"`
	IataCode    string `json:"iata_code"`    //Filtering by Airport IATA code.
	IcaoCode    string `json:"icao_code"`    //Filtering by Airport ICAO code.
	CityCode    string `json:"city_code"`    //Filtering by IATA City code.
	CountryCode string `json:"country_code"` //Filtering by Country ISO 2 code
}

type Airport struct {
	common.Model
	Name        string  `json:"name"`         //Public name.
	IataCode    string  `json:"iata_code"`    //Official IATA code.
	IcaoCode    string  `json:"icao_code"`    //Official ICAO code.
	Lat         float64 `json:"lat"`          //Geo Latitude
	Lng         float64 `json:"lng"`          //Geo Longitude.
	City        string  `json:"city"`         //Airport metropolitan city name
	CityCode    string  `json:"city_code"`    //Airport metropolitan 3 letter city code
	UNLocode    string  `json:"un_locode"`    //United Nations location code.
	Timezone    string  `json:"timezone"`     //Airport location timezone.
	CountryCode string  `json:"country_code"` //ISO 2 country code
	Departures  int     `json:"departures"`   //Total departures from airport per year.
}

func (a Airport) TableName() string {
	return "airports"
}

func ExistAirportByIataCode(iatacode string) (bool, error) {
	var airport Airport
	err := common.DB.Select("id").Where("iata_code = ? AND deleted_on = ? ", iatacode, 0).First(&airport).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if airport.ID > 0 {
		return true, nil
	}
	return false, nil
}

func ExistAirportByIcaoCode(icao_code string) (bool, error) {
	var airport Airport
	err := common.DB.Select("id").Where("icao_code = ? AND deleted_on = ? ", icao_code, 0).First(&airport).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if airport.ID > 0 {
		return true, nil
	}
	return false, nil
}

func AddAirport(airport Airport) error {
	if err := common.DB.Create(&airport).Error; err != nil {
		return err
	}
	return nil
}

func GetAirports(pagenum int, pagesize int, maps interface{}) ([]Airport, error) {
	var (
		airports []Airport
		err      error
	)
	if pagesize > 0 && pagenum > 0 {
		err = common.DB.Where(maps).Find(&airports).Offset(pagenum).Limit(pagesize).Error
	} else {
		err = common.DB.Where(maps).Find(&airports).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return airports, nil
}

func GetAirportsTotal(maps interface{}) (int, error) {
	var count int
	if err := common.DB.Model(&Airport{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func ExistAirportByID(id int) (bool, error) {
	var airport Airport
	err := common.DB.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&airport).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if airport.ID > 0 {
		return true, nil
	}
	return false, nil
}

func DelteAirport(id int) error {
	if err := common.DB.Where("id = ?", id).Delete(&Airport{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateAirport(id int, data interface{}) error {
	if err := common.DB.Model(&Airport{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
