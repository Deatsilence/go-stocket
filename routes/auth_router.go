package routes

import (
	controller "github.com/Deatsilence/go-stocket/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/verifyemail", controller.VerifyEmail())
	incomingRoutes.POST("users/signup", controller.SignUp())
	incomingRoutes.POST("users/login", controller.Login())
	incomingRoutes.POST("users/logout", controller.Logout())
}
