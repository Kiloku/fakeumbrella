package main

import (
	"fmt"
	
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

type Customer struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Location string `json:"location"`
	Contact string `json:"contact"`
	Telephone string `json:"telephone"`
	Employees uint `json:"employees"`
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

func DeleteCustomer(c *gin.Context) {
	id := c.Params.ByName("id")
	var customer Customer
	d := db.Where("id = ?", id).Delete(&customer)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdateCustomer(c *gin.Context) {
	var customer Customer
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&customer).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&customer)

	db.Save(&customer)
	c.JSON(200, customer)
}

func CreateCustomer(c *gin.Context) {
	var customer Customer
	c.BindJSON(&customer)

	db.Create(&customer)
	c.JSON(200, customer)
}

func GetCustomer(c *gin.Context) {
	id := c.Params.ByName("id")
	var customer Customer
	if err := db.Where("id = ?", id).First(&customer).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, customer)
	}
}

func GetCustomers(c *gin.Context) {
	var customers []Customer
	err := db.Find(&customers).Error
	if err != nil {
    	c.AbortWithStatus(404)
    	fmt.Println(err)
 	} else {
    	c.JSON(200, customers)
	}
}

func GetWeatherCustomer(c *gin.Context) { 
	id := c.Params.ByName("id")
	var customer Customer
	if err := db.Where("id = ?", id).First(&customer).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		fmt.Println(customer.Location)
		c.JSON(200, gin.H{"Customer":customer, "RainExpected":willRainAt(customer.Location)})
		//c.JSON(200, RequestWeatherAt(customer.Location))
		//fmt.Println(willRainAt(customer.Location))
	}
}

func willRainAt(location string) bool{
	weather := RequestWeatherAt(location)
	for _, element := range weather.List {
		if element.Weather[0].Main == "Rain"{
			return true
		}
	}
	return false
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