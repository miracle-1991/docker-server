package schedules

type ScheduleRequest struct {
	DepIata     string `json:"dep_iata"`     //Filtering by departure Airport IATA code.
	DepIcao     string `json:"dep_icao"`     //Filtering by departure Airport ICAO code.
	ArrIata     string `json:"arr_iata"`     //Filtering by arrival Airport IATA code.
	ArrIcao     string `json:"arr_icao"`     //Filtering by arrival Airport ICAO code.
	AirlineIcao string `json:"airline_icao"` //Filtering by Airline ICAO code.
	AirlineIata string `json:"airline_iata"` //Filtering by Airline IATA code.
}

type ScheduleResponse struct {
	AirlineIata     string `json:"airline_iata"`      //Airline IATA code.
	AirlineIcao     string `json:"airline_icao"`      //Airline ICAO code.
	FlightIata      string `json:"flight_iata"`       //Flight IATA code-number
	FlightIcao      string `json:"flight_icao"`       //Flight ICAO code-number.
	FlightNumber    string `json:"flight_number"`     //Flight number only.
	CSAirlineIata   string `json:"cs_airline_iata"`   //Codeshared airline IATA code.
	CSFlightIata    string `json:"cs_flight_iata"`    //Codeshared flight IATA code-number.
	CSFlightNumber  string `json:"cs_flight_number"`  //Codeshared flight number.
	DepIata         string `json:"dep_iata"`          //Departure airport IATA code.
	DepIcao         string `json:"dep_icao"`          //Departure airport ICAO code.
	DepTerminal     string `json:"dep_terminal"`      //Estimated departure terminal.
	DepGate         string `json:"dep_gate"`          //Estimated departure gate.
	DepTime         string `json:"dep_time"`          //Departure time in the airport time zone.
	DepTimeTs       uint32 `json:"dep_time_ts"`       //Departure UNIX timestamp.
	DepTimeUtc      string `json:"dep_time_utc"`      //Departure time in UTC time zone.
	DepEstimated    string `json:"dep_estimated"`     //Updated departure time in the airport time zone.
	DepEstimatedTs  uint32 `json:"dep_estimated_ts"`  //Updated departure UNIX timestamp.
	DepEstimatedUtc string `json:"dep_estimated_utc"` //Updated departure time in UTC time zone.
	DepActual       string `json:"dep_actual"`        //Actual departure time in the airport time zone.
	DepActualTs     uint32 `json:"dep_actual_ts"`     //Actual departure UNIX timestamp.
	DepActualUtc    string `json:"dep_actual_utc"`    //Actual departure time in UTC time zone.
	ArrIata         string `json:"arr_iata"`          //Arrival airport IATA code.
	ArrIcao         string `json:"arr_icao"`          //Arrival airport ICAO code.
	ArrTerminal     string `json:"arr_terminal"`      //Estimated arrival terminal.
	ArrGate         string `json:"arr_gate"`          //Estimated arrival gate.
	ArrBaggage      string `json:"arr_baggage"`       //Arrival baggage claim carousel number.
	ArrTime         string `json:"arr_time"`          //Arrival time in the airport time zone.
	ArrTimeTs       uint32 `json:"arr_time_ts"`       //Arrival UNIX timestamp.
	ArrTimeUtc      string `json:"arr_time_utc"`      //Arrival time in UTC time zone.
	ArrEstimated    string `json:"arr_estimated"`     //Updated arrival time in the airport time zone.
	ArrEstimatedTs  uint32 `json:"arr_estimated_ts"`  //Updated arrival UNIX timestamp.
	ArrEstimatedUtc string `json:"arr_estimated_utc"` //Updated arrival time in UTC time zone.
	ArrActual       string `json:"arr_actual"`        //Actual arrival time in the airport time zone.
	ArrActualTs     uint32 `json:"arr_actual_ts"`     //Actual arrival UNIX timestamp.
	ArrActualUtc    string `json:"arr_actual_utc"`    //Actual arrival time in UTC time zone.
	Duration        int    `json:"duration"`          //Estimated flight time (in minutes).
	Delayed         int    `json:"delayed"`           //(deprecated) Estimated flight delay time (in minutes).
	DepDelayed      int    `json:"dep_delayed"`       //Estimated time of flight departure delay (in minutes).
	ArrDelayed      int    `json:"arr_delayed"`       //Estimated time of flight arrival delay (in minutes).
	Status          string `json:"status"`            //Flight status - scheduled, cancelled, active, landed.
}
