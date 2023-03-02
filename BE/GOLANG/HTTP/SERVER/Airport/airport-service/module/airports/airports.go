package airports

type AirportRequest struct {
	IataCode    string `json:"iata_code"`    //Filtering by Airport IATA code.
	IcaoCode    string `json:"icao_code"`    //Filtering by Airport ICAO code.
	CityCode    string `json:"city_code"`    //Filtering by IATA City code.
	CountryCode string `json:"country_code"` //Filtering by Country ISO 2 code
}

type AirportResponse struct {
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
