package main

import (
	"fmt"
	"time"
	"encoding/json"
	
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error


//Generated with "JSON-to-Go" tool. (https://mholt.github.io/json-to-go/) and modified by hand to remove unnecessary elements
type WeatherData struct {
	Message float64 `json:"message"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  float64 `json:"pressure"`
			SeaLevel  float64 `json:"sea_level"`
			GrndLevel float64 `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    int     `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   float64 `json:"deg"`
		} `json:"wind"`
		Rain struct {
			ThreeH float64 `json:"3h"`
		} `json:"rain"`
		Sys struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
		Snow  struct {
			ThreeH float64 `json:"3h"`
		} `json:"snow,omitempty"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country string `json:"country"`
	} `json:"city"`
}

type RainForecast struct{
	WillRain bool `json:"willRain"`
	When int `json:"when"`
}


func main(){
	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.AutoMigrate(&Customer{})

  	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"PUT", "POST", "PATCH", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin","Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge: 12 * time.Hour,
	}))
	r.GET("/customers/", GetCustomers)
	r.GET("/customers/:id", GetCustomer)
	r.POST("/customers", CreateCustomer)
	r.PUT("/customers/:id", UpdateCustomer)
	r.DELETE("/customers/:id", DeleteCustomer)
	r.GET("/weather-for-customer/:id", GetWeatherCustomer)
	r.GET("/forecasts/cities/", GetForecastCities)
	r.GET("/forecasts/customers/", GetForecastsCustomers)
	for _, element := range FindCitiesServed(){
		fmt.Println(element)
	}
	r.Run(":8080")
}



func GetWeatherCustomer(c *gin.Context) { 
	id := c.Params.ByName("id")
	var customer Customer
	if err := db.Where("id = ?", id).First(&customer).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		fmt.Println(customer.Location)
		var forecast RainForecast
		forecast.WillRain, forecast.When = willRainAt(customer.Location)
		if !forecast.WillRain{
			c.JSON(200, gin.H{"Customer":customer, "RainExpected":gin.H{"willRain":forecast.WillRain}})
		} else {
			c.JSON(200, gin.H{"Customer":customer, "RainExpected":forecast})
		}
	}
}

func GetForecastsCustomers(c *gin.Context){
	var customersInCities []CityWithCustomers
	var forecasts []gin.H 
	customersInCities = FindCustomersInCities()
	for _, element := range customersInCities {
		var forecast RainForecast
		forecast.WillRain, forecast.When = willRainAt(element.City)
		forecasts = append(forecasts, gin.H{"City": element.City, "Customers" :  element.Customers, "Forecast" : forecast})
	} 
	c.JSON(200, forecasts)
}

func GetForecastCities(c *gin.Context){
	cities := FindCitiesServed()
	var forecasts []gin.H
	for _, element := range cities {
		var forecast RainForecast
		forecast.WillRain, forecast.When = willRainAt(element)
		forecasts = append(forecasts, gin.H{"City": element, "Forecast":forecast})
	}
	c.JSON(200, forecasts)
}

func willRainAt(location string) (bool, int) {
	weather := RequestWeatherAt(location)
	for _, element := range weather.List {
		if element.Weather[0].Main == "Rain"{
			return true, element.Dt
		}
	}
	return false, 0
}

func willRainAtParsedTime(location string) (bool, time.Time) {
	weather := RequestWeatherAt(location)
	for _, element := range weather.List {
		if element.Weather[0].Main == "Rain"{
			return true, time.Unix(int64(element.Dt), 0)
		}
	}
	return false, time.Unix(0,0)
}


func RequestWeatherAt(location string) (weatherData WeatherData){
	response, err := http.Get("http://api.openweathermap.org/data/2.5/forecast?q="+location+"&APPID=2aa5d8c417225481239400cc3a8a5409")
	if err != nil {
		fmt.Printf("The HTTP request failed with err %s", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data)) //Remove
		json.Unmarshal(data, &weatherData)	
		
	}
	return 
}

func FindCitiesServed() []string{
	var customers []Customer
	db.Find(&customers)
	var cities []string
	for _, element := range customers {
		if (!Contains(cities, element.Location)){
			cities = append(cities, element.Location)
		}
	}

	return cities
}

type CityWithCustomers struct{
	City string `json:"city"`
	Customers []Customer `json:"customers"` 
}

func FindCustomersInCities() []CityWithCustomers {
	cities := FindCitiesServed() 
	var customersInCities []CityWithCustomers
	for _, element := range cities {
		var customers []Customer 
		db.Where(&Customer{Location:element}).Find(&customers)
		customersInCities = append(customersInCities, CityWithCustomers{City: element, Customers: customers})
	}
	return customersInCities
}

func Contains(haystack []string, needle string) bool {
	for _, element := range haystack {
		if needle == element{
			return true
		}
	}
	return false
}