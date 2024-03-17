package routes

import (
	controller "github.com/Deatsilence/go-stocket/pkg/controllers"
	"github.com/Deatsilence/go-stocket/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("products/add", controller.AddAProduct())
	// incomingRoutes.DELETE("products/delete/:productid", controller.DeleteAProduct())
	// incomingRoutes.GET("products", controller.GetProducts())
	// incomingRoutes.GET("products/:productid", controller.GetProduct())
	// incomingRoutes.PUT("products/update/:productid", controller.UpdateAProduct())
	// incomingRoutes.PATCH("products/update/:productid", controller.UpdateAProduct())
}
