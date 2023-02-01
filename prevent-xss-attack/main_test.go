package main

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXSSSanitization(t *testing.T) {
	payload := `<person>
		<name>rohit</name>
	</person>`

	req, err := http.NewRequest("POST", "http://localhost:8080/person", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/xml")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode, "Expected status code to be 200 OK")

	var response struct {
		Person struct {
			Name    string `xml:"name"`
			Address string `xml:"address"`
		} `xml:"person"`
	}

	err = xml.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotContains(t, response.Person.Address, "<script>", "Expected address to be sanitized")
}
