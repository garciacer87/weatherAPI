# WeatherAPI

This is an API used to consume basic weather info from a specific city.
The source is provided by an external API (https://openweathermap.org/).
This application is developed in Go. So, in order to run it, first download the project. Then, you will have 3 options available:

**Using go run command:**
If you have Go installed on your machine, open a terminal, locate the root of the project you just downloaded, and run the following command. 
```sh
$ go run main.go
```
That will execute the application using the environment variables defined in .env file located in the root of the project.

**Using go build command:**
If you have Go installed on your machine, open a terminal, locate the root of the project you just downloaded, and run the following command. 
```sh
$ go build .
```
That will create an executable file (.exe, .sh, etc depending of your OS) in the root of the project, available to run it, using the environment variables defined in .env file.

**Using docker image:**
If you have docker installed on your machine, open a terminal, locate the root of the project you just downloaded, and create a docker container using the build command.
```sh
$ docker build -t weatherapi:1.0 .
```
Then, you can run the container, passing the required environment variables in the command. For example:
```sh
$ docker run --env OPENWEATHERMAP_HOST=http://api.openweathermap.org --env OPENWEATHERMAP_APIKEY=1508a9a4840a5574c822d70ca2132032 --rm -it -p 8080:8080/tcp weatherapi:1.0
```

# Environment Variables
This application uses a few environment variables to work.
  - OPENWEATHERMAP_HOST **(required)**: this is used to define the OpenWeather API host. Like: http://api.openweathermap.org
  - OPENWEATHERMAP_APIKEY **(required)**: this is used to define the API KEY needed to consume the OpenWeather API.
  - OPENWEATHERMAP_UNIT **(optional)**: this is used to set the unit measurement. Values permitted: "metric" (Cº and m/s), "imperial" (ºF and miles/hr). Default value is "metric".
  - CACHE_DURATION **(optional)**: this is used to set the expiration of cache. This value is represented in Minutes. Default value is 2.

# Endpoints available
 - /health (GET): used as a health check to get an OK response if the service is up.
 - /weather?city=$CITY&country=$COUNTRY (GET): used to get weather info of a city. Query parameters city and country must fulfill the following:
    - City: is required and must be a string of [a-zA-z]. Otherwise, you will get a bad request response.
    - Country: is required and must be a 2 characters string in lowercase. Otherwise, you will get a bad request response.

# Response
The API will always response a JSON. If the response is not 200, the response will be something like this:
```code
{
    "code": 400,
    "message": [
        "country must be a two characters string in lowercase"
    ]
}
```
If you get a successful response, you wil get something like this:
```code
{
    "cloudiness": "clear sky",
    "forecast": [
        {
            "cloudiness": "clear sky",
            "forecasted_datetime": "25/01/2021 09:00",
            "humidity": "58%",
            "maximum_temperature": "22ºC",
            "minimum_temperature": "21ºC",
            "real_feel_temperature": "20ºC",
            "temperature": "21ºC"
        },
        {
            "cloudiness": "clear sky",
            "forecasted_datetime": "25/01/2021 12:00",
            "humidity": "41%",
            "maximum_temperature": "28ºC",
            "minimum_temperature": "26ºC",
            "real_feel_temperature": "25ºC",
            "temperature": "26ºC"
        },
        {
            "cloudiness": "clear sky",
            "forecasted_datetime": "25/01/2021 15:00",
            "humidity": "28%",
            "maximum_temperature": "32ºC",
            "minimum_temperature": "31ºC",
            "real_feel_temperature": "28ºC",
            "temperature": "31ºC"
        }
    ],
    "geo_coordinates": "[-33.456900, -70.648300]",
    "humidity": "77%",
    "location_name": "Santiago, CL",
    "maximum_temperature": "21ºC",
    "minimum_temperature": "16ºC",
    "pressure": "1011 hpa",
    "real_feel_temperature": "18ºC",
    "requested_time": "07:05",
    "sunrise": "06:58",
    "sunset": "20:51",
    "temperature": "19ºC",
    "wind": "3.09 m/s South-SouthEast"
}
```
**Thanks! Enjoy!!!**
