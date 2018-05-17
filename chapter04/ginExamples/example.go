package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/time", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"serverTime": time.Now().UTC(),
			"values": []gin.H{
				gin.H{
					"dayOfWeek": "Monday",
				},
				gin.H{
					"dayOfWeek": "Tuesday",
				},
			},
		})
	})
	router.Run(":8080")
}
