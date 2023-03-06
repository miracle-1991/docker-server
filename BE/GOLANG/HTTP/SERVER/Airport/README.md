# Airport Server
The project library can update the nearest flight to the database and list the flights in the database on request.

If you run it at goland, before start, you must change the [config](./cmd/config.yaml)

## Update Flights:
* Send a GET request to http://localhost:8006/update/:{country}
  * such as requesting to update the flight data of Singapore: http://localhost:8006/update/SG
  ```
  curl -v -X GET http://localhost:8006/update/SG
  ```
## Get Flights:
* Send a post request to http://localhost:8006/get, the request body is:
```
curl -v -X POST -H "Content-Type: application/json" -d '{"country_code": "SG", "arr_date": "2023-03-03", "dep_date": "2023-03-03", "output_path": "/data"}' http://localhost:8006/get 
```
post body:
```
{
  "country_code": "SG",
  //At least one of arr_date and dep_date is not empty
  "arr_date": "2023-03-03", 
  "dep_date": "2023-03-03",
  //If this parameter is empty, the result will be returned in the form of json, otherwise the result will be placed in the specified directory
  "output_path": "/data" 
}
```
# Run At Mac
just open this project in GoLand and run it

# Run At Docker
start
```
docker-compose up -d
```
stop
```
docker-compose stop
```


