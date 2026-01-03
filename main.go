package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type FibonacciResponse struct {
	OrderNumber  int    `json:"order_number"`
	FibonacciNum uint64 `json:"fibonacci_number"`
	StatusCode   int    `json:"status_code"`
	Message      string `json:"message,omitempty"`
}

// Time complexity: O(2^n) – intentionally slow
func recursiveFibonacci(n int) uint64 {
	if n <= 1 {
		return uint64(n)
	}
	return recursiveFibonacci(n-1) + recursiveFibonacci(n-2)
}

func fibonacciHandler(c *gin.Context) {
	orderStr := c.Param("order_number")

	order, err := strconv.Atoi(orderStr)
	if err != nil || order < 0 || order > 93 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": http.StatusBadRequest,
			"message":     "order_number must be between 0 and 93",
		})
		return
	}

	start := time.Now()
	fibNum := recursiveFibonacci(order)
	elapsed := time.Since(start)

	c.JSON(http.StatusOK, FibonacciResponse{
		OrderNumber:  order,
		FibonacciNum: fibNum,
		StatusCode:   http.StatusOK,
		Message:      "Computed in " + elapsed.String(),
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Register Prometheus metrics
	p := ginprometheus.NewWithConfig(ginprometheus.Config{
		Subsystem: "gin",
	})
	p.Use(router)

	// API routes
	api := router.Group("/api/v1")
	{
		api.GET("/fibonacci/:order_number", fibonacciHandler)
	}

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("run server")
	}
}
