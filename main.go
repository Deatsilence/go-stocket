package main

import (
	"log"
	"os"

	routes "github.com/Deatsilence/go-stocket/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file `%v`", err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.PasswordRoutes(router)
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.ProductRoutes(router)

	router.Run(":" + port)
}
