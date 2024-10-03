package router

import (
	"github.com/abhishekdas600/movierecserver/services/user"
	"github.com/gin-gonic/gin"
)

func SetupUser(router *gin.Engine) {
    
    router.GET("/user", user.GetUserDetails)
}