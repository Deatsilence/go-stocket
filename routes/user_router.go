package routes

import (
	controller "github.com/Deatsilence/go-stocket/pkg/controllers"
	"github.com/Deatsilence/go-stocket/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/api/users", controller.GetUsers())
	incomingRoutes.GET("/api/users/:userid", controller.GetUser())

}
