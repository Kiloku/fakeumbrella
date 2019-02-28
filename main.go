package main

import (
	"fmt"
	"time"
	"encoding/json"
	
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
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
	When time.Time `json:"when"`
}


func main(){
  db, err = gorm.Open("sqlite3", "./gorm.db")
  if err != nil {
  	fmt.Println(err)
  }
  defer db.Close()

  db.AutoMigrate(&Customer{})

  r := gin.Default()
  r.GET("/customers/", GetCustomers)
  r.GET("/customers/:id", GetCustomer)
  r.GET("/weather/:id", GetWeatherCustomer)
  r.POST("/customers", CreateCustomer)
  r.PUT("/customers/:id", UpdateCustomer)
  r.DELETE("/customers/:id", DeleteCustomer)

  for _, element := range FindCustomerCities(){
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

func willRainAt(location string) (bool, time.Time) {
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
	//weatherData := WeatherData{}
	return 
	
	//jsonData := map[string]string
}

func FindCustomerCities() []string{
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

func Contains(haystack []string, needle string) bool {
	for _, element := range haystack {
		if needle == element{
			return true
		}
	}
	return false
}