package routes

import (
	controller "github.com/Deatsilence/go-stocket/pkg/controllers"
	"github.com/Deatsilence/go-stocket/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/api/products/add", controller.AddAProduct())
	incomingRoutes.DELETE("/api/products/delete/:productid", controller.DeleteAProduct())
	incomingRoutes.GET("/api/products", controller.GetProducts())
	incomingRoutes.GET("/api/products/:productid", controller.GetProduct())
	incomingRoutes.GET("/api/products/search", controller.SearchByBarcodePrefix())
	incomingRoutes.PUT("/api/products/update/:productid", controller.UpdateAProduct())
	incomingRoutes.PATCH("/api/products/updatepartially/:productid", controller.UpdateSomePropertiesOfProduct())
}
