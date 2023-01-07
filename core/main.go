package main

import (
	"fmt"
	"log"
	"runtime"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func v1EndpointHandler(c *gin.Context) {
	c.String(200, "v1: %s %s", c.Request.Method, c.Request.URL.Path)
}

func v2EndpointHandler(c *gin.Context) {
	c.String(200, "v2: %s %s", c.Request.Method, c.Request.URL.Path)
}

func add(c *gin.Context) {
	x, _ := strconv.ParseFloat(c.Param("x"), 64)
	y, _ := strconv.ParseFloat(c.Param("y"), 64)
	c.String(200, fmt.Sprintf("%f", x+y))
}

type MultiplexParams struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func multiplex(c *gin.Context) {
	var ap MultiplexParams
	if err := c.ShouldBindJSON(&ap); err != nil {
		c.JSON(400, gin.H{"error": "Calculator error"})
		return
	}

	c.JSON(200, gin.H{"answer": ap.X * ap.Y})
}

type PrintJob struct {
	JobId int `json:"jobId" binding:"required,gte=10000"`
	Pages int `json:"pages" binding:"required,gte=1,lte=100"`
}

func FindUserAgent() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println(ctx.GetHeader("User-Agent"))
		// before calling handler
		ctx.Next()
		// after
	}
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	// middleware
	router.Use(FindUserAgent())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "CORS and middleware work as well!"})
	})

	// basic
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	router.GET("/os", func(ctx *gin.Context) {
		ctx.String(200, runtime.GOOS)
	})

	// versioning
	v1 := router.Group("/v1")
	v1.GET("/products", v1EndpointHandler)
	v1.GET("/products/:productId", v1EndpointHandler)
	v1.POST("/products", v1EndpointHandler)
	v1.PUT("/products/:productId", v1EndpointHandler)
	v1.DELETE("/products/:productId", v1EndpointHandler)

	v2 := router.Group("/v2")
	v2.GET("/products", v2EndpointHandler)
	v2.GET("/products/:productId", v2EndpointHandler)
	v2.POST("/products", v2EndpointHandler)
	v2.PUT("/products/:productId", v2EndpointHandler)
	v2.DELETE("/products/:productId", v2EndpointHandler)

	// url param
	router.GET("/add/:x/:y", add)

	// http message body
	router.POST("/multiple", multiplex)

	// validation
	router.POST("/print", func(ctx *gin.Context) {
		var p PrintJob
		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input!"})
			log.Println(err)
			return
		}
	})

	router.Run(":5000")
}
