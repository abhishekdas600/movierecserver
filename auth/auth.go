package auth

import (
    "log"
    "os"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/markbates/goth"
    "github.com/markbates/goth/gothic"
    "github.com/markbates/goth/providers/google"
)

var (
    MaxAge = 86400 * 30 
    IsProd = false      
)

func NewAuth(r *gin.Engine) {
    
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

   
    key := os.Getenv("AUTH_KEY")
    if key == "" {
        log.Fatal("AUTH_KEY environment variable not set")
    }

    googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
    googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

    store := cookie.NewStore([]byte(key))
    store.Options(sessions.Options{
        Path:     "/",
        MaxAge:   MaxAge,
        HttpOnly: true,
        Secure:   IsProd,
    })

    r.Use(sessions.Sessions("mysession", store))

    gothic.Store = store

    goth.UseProviders(
        google.New(googleClientId, googleClientSecret, "http://localhost:8080/auth/google/callback"),
    )
}

