package controllers

import (
	"expense-tracker/pkg/models"
	"expense-tracker/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTotalBalance(c *gin.Context) {
	var balance models.Balance
	err := utils.BalanceCollection.FindOne(c, bson.M{}).Decode(&balance)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{"total_balance": 0})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_balance": balance.Total_Balance})
}
