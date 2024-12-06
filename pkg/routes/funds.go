package routes

import (
	"expense-tracker/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func FundRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/funds/add", controllers.AddFunds)
	}
}
