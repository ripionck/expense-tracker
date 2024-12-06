package controllers

import (
	"expense-tracker/pkg/models"
	"expense-tracker/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddExpenses(c *gin.Context) {
	var expense models.Expense
	if err := c.BindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	expense.Inserted_At = time.Now()
	expense.ID = primitive.NewObjectID()
	expense.Expense_ID = expense.ID.Hex()

	_, err := utils.ExpenseCollection.InsertOne(c, expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expenses added successfully"})
}
