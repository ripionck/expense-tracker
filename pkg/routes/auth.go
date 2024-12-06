package routes

import (
	"expense-tracker/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/users/register", controllers.Register())
		v1.POST("/users/login", controllers.Login())
	}
}
