package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

type Person struct {
	Name    string `json:"name" binding:"required"`
	Age     int    `json:"age" binding:"required"`
	Address string `json:"address" binding:"required,addressLength=30"`
}

func addressLength(v interface{}, param string) (bool, error) {
	value, ok := v.(string)
	if !ok {
		return false, nil
	}
	maxLength, err := strconv.Atoi(param)
	if err != nil {
		return false, err
	}
	return len(value) <= maxLength, nil
}

func main() {
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// v.RegisterValidation("addressLength", addressLength)
		v.RegisterValidation("addressLength", func(fl validator.FieldLevel) bool {
			fmt.Println(fl.Field())
			return true
		})
	}

	r.POST("/person", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"person": person})
	})

	r.Run()
}
