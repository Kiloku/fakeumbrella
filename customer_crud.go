package main

import (

	"fmt"
	"github.com/gin-gonic/gin"
)

type Customer struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Location string `json:"location"`
	Contact string `json:"contact"`
	Telephone string `json:"telephone"`
	Employees uint `json:"employees"`
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