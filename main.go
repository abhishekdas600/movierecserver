package main

import (
	"fmt"
	"github.com/abhishekdas600/movierecserver/auth"
	"github.com/abhishekdas600/movierecserver/router"
	"os"

	
)

func main() {
    
	r := router.SetupRouter()
	auth.NewAuth(r)
	router.SetupRoutes(r)
	port := os.Getenv("HTTP_PLATFORM_PORT")
	if port == "" {
		port = "8080"
	}

	err := r.Run(":" + port)
	if err != nil {
		fmt.Println("ERROR starting server:", err)
	}
}
