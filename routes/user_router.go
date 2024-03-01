package routes

import (
	controller "github.com/Deatsilence/go-stocket/pkg/controllers"
	"github.com/Deatsilence/go-stocket/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.AuthMiddleware())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:userId", controller.GetUser())

}
