# Airport Server
The project library can update the nearest flight to the database and list the flights in the database on request.

Before start, you must change the [config](./cmd/config.yaml)

## Update Flights:
* Send a GET request to http://localhost:8006/update/:{country}
  * such as requesting to update the flight data of Singapore: http://localhost:8006/update/SG
## Get Flights:
* Send a post request to http://localhost:8006/get, the request body is:
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
# Start At Mac
just open this project in GoLand and run it

# Start At Docker
it will auto update every 6 hours because of this command in Dockerfile:
```
RUN echo "0 */6 * * * /app/my-flight-service > /dev/null 2>&1" >> /etc/crontab
```
build:
```
docker build -t ff/airport:latest -f Dockerfile .
```
run 
```
docker run -p 8006:8006 -v ~/Downloads/positioning-data:/data ff/airport:latest
```
if you meet error:
```
Host '' is not allowed to connect to this MySQL server
```
you can try this solution: [create root](./airport-service/module/airports/sql/create-root.sql)
