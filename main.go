package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abhishekdas600/movierecserver/auth"
	"github.com/abhishekdas600/movierecserver/db" 
	"github.com/abhishekdas600/movierecserver/router"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv" 
)

func main() {
	
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Error loading .env.local file: %v", err)
	}

	
	db.Init() 

	
	r := router.SetupRouter()

	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	
	auth.NewAuth(r)
	router.SetupRoutes(r)

	
	port := os.Getenv("HTTP_PLATFORM_PORT")
	if port == "" {
		port = "8080" 
	}

	fmt.Printf("Starting server on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

