package routes

import (
	"expense-tracker/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func ExpensesRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/expenses/add", controllers.AddExpenses)
	}
}
