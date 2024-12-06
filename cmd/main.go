package main

import (
	"expense-tracker/pkg/routes"
	"expense-tracker/pkg/utils"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Initialize the database
	client := utils.ConnectDB()
	if client == nil {
		log.Fatal("Failed to initialize MongoDB connection")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatal("PORT environment variable not set")
	}
	fmt.Printf("Server starting on port: %s\n", port)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	routes.AuthRoutes(router)
	routes.FundRoutes(router)
	routes.ExpensesRoutes(router)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
