package main

import (
  "fmt"

  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

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
  r.POST("/customers", CreateCustomer)
  r.PUT("/customers/:id", UpdateCustomer)
  r.DELETE("/customers/:id", DeleteCustomer)
  r.Run(":8080")
}

func DeleteCustomer(c *gin.Context) {
	id := c.Params.ByName("id")
	var Customer Customer
	d := db.Where("id = ?", id).Delete(&Customer)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdateCustomer(c *gin.Context) {
	var Customer Customer
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&Customer).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&Customer)

	db.Save(&Customer)
	c.JSON(200, Customer)
}

func CreateCustomer(c *gin.Context) {
	var Customer Customer
	c.BindJSON(&Customer)

	db.Create(&Customer)
	c.JSON(200, Customer)
}

func GetCustomer(c *gin.Context) {
	id := c.Params.ByName("id")
	var Customer Customer
	if err := db.Where("id = ?", id).First(&Customer).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, Customer)
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