package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter(t *testing.T) {
	router := gin.Default()
	router.Use(rateLimit(1000))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Send multiple requests within a minute to see if rate limiting is enforced
	for i := 0; i < 2000; i++ {
		router.ServeHTTP(w, req)
		if w.Code != 200 && w.Code != 429 {
			t.Fatalf("Unexpected response code: %d", w.Code)
		}
		if w.Code == 429 {
			break
		}
	}

	// Check that the rate limit was enforced by verifying that at least one 429 response was received
	if w.Code != 429 {
		t.Fatalf("Expected at least one 429 response, but got %d", w.Code)
	}

	// Wait for a few minutes to allow the rate limiter to reset
	time.Sleep(1 * time.Minute)

	// Send a new request to see if it is accepted
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check that the new request was accepted
	if w.Code != 200 {
		t.Fatalf("Unexpected response code: %d", w.Code)
	}
}
