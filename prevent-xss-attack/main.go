package main

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

type Person struct {
	XMLName xml.Name `xml:"person"`
	Name    string   `xml:"name"`
	Age     int      `xml:"age"`
}

func main() {
	r := gin.Default()

	r.POST("/person", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindXML(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("before ", person.Name)
		p := bluemonday.NewPolicy()
		p.AllowStandardURLs()

		person.Name = p.Sanitize(person.Name)
		fmt.Println("after ", person.Name)
		c.XML(http.StatusOK, gin.H{"person": person})
	})

	r.Run()
}
