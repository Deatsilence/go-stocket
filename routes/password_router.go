package routes

import (
	controller "github.com/Deatsilence/go-stocket/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func PasswordRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/passwordreset/request", controller.RequestPasswordReset())
	incomingRoutes.POST("/passwordreset/confirm", controller.ResetPassword())
	incomingRoutes.POST("/passwordreset/changepassword", controller.ChangePassword())
}
