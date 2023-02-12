package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func rateLimit(maxRate int) gin.HandlerFunc {
	limiter := sync.Map{}
	// go func() {
	// 	for {
	// 		time.Sleep(time.Minute)
	// 		limiter.Range(func(key, value interface{}) bool {
	// 			limiter := value.(*rate.Limiter)
	// 			if time.Since(limiter.Allow().Time) > 2*time.Minute {
	// 				limiter.Delete(key)

	// 			}
	// 			return true
	// 		})
	// 	}
	// }()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l, _ := limiter.LoadOrStore(ip, rate.NewLimiter(rate.Limit(maxRate), maxRate))
		if !l.(*rate.Limiter).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// https://www.makeuseof.com/rate-limiting-go-applications/
// https://www.google.com/search?q=how+to+cleanup+sync+map+in+rate+limiter+middleware&oq=how+to+cleanup+sync+map+in+rate+limiter+middleware&aqs=edge..69i57j0i546l3.26146j0j1&sourceid=chrome&ie=UTF-8
// cleanup deletes old entries
// func  cleanup(frequency time.Duration) {
// 	for {
// 		time.Sleep(frequency)

// 		l.Lock()
// 		for ip, v := range l.visitors {
// 			if time.Since(v.lastSeen) > frequency {
// 				delete(l.visitors, ip)

// 				m.rl.Delete(ip)
// 			}
// 		}
// 		l.Unlock()
// 	}
// }

func main() {
	router := gin.Default()
	router.Use(rateLimit(10))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hi this is from Rate Limiter"})
	})

	router.Run(":8080")
}
