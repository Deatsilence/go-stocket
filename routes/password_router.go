package routes

import (
	controller "github.com/Deatsilence/go-stocket/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func PasswordRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/api/passwordreset/request", controller.RequestPasswordReset())
	incomingRoutes.POST("/api/passwordreset/confirm", controller.ResetPassword())
	incomingRoutes.POST("/api/passwordreset/changepassword", controller.ChangePassword())
}
