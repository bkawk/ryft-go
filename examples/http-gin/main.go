package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/bkawk/ryft-go"
	"github.com/gin-gonic/gin"
)

type createCustomerBody struct {
	Email     string            `json:"email" binding:"required,email"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Metadata  map[string]string `json:"metadata"`
}

func main() {
	secretKey := os.Getenv("RYFT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("RYFT_SECRET_KEY is required")
	}

	client, err := ryft.NewClient(ryft.Config{
		SecretKey: secretKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/customers", func(c *gin.Context) {
		var body createCustomerBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_request",
				"message": err.Error(),
			})
			return
		}

		customer, err := client.Customers.Create(context.Background(), ryft.CreateCustomerRequest{
			Email:     body.Email,
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Metadata:  body.Metadata,
		})
		if err != nil {
			var apiErr *ryft.APIError
			if errors.As(err, &apiErr) {
				c.JSON(apiErr.Status, gin.H{
					"error":     apiErr.Code,
					"message":   apiErr.Message,
					"requestId": apiErr.RequestID,
					"details":   apiErr.Errors,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, customer)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
