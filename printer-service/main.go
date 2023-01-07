package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PrintJob struct {
	Format    string `json:"format" binding:"required"`
	InvoiceId int    `json:"invoiceId" binding:"required,gte=0"`
	JobId     int    `json:"jobId" binding:"gte=0"`
}

func main() {
	router := gin.Default()
	router.POST("/print-jobs", func(c *gin.Context) {
		var p PrintJob
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input!"})
			return
		}
		log.Printf("PrintService: creating new print job from invoice #%v...", p.InvoiceId)
		rand.Seed(time.Now().UnixNano())
		p.JobId = rand.Intn(1000)
		log.Printf("PrintService: created print job #%v", p.JobId)
		c.JSON(http.StatusOK, p)
	})
	router.Run(":5000")
}
