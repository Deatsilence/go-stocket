package routes

import (
	controller "github.com/Deatsilence/go-stocket/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/api/users/verifyemail", controller.VerifyEmail())
	incomingRoutes.POST("/api/users/signup", controller.SignUp())
	incomingRoutes.POST("/api/users/login", controller.Login())
	incomingRoutes.POST("/api/users/logout", controller.Logout())
}
