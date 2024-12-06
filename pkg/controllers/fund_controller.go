package controllers

import (
	"expense-tracker/pkg/models"
	"expense-tracker/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddFunds(c *gin.Context) {
	var fund models.Fund
	if err := c.BindJSON(&fund); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err := utils.FundCollection.InsertOne(c, fund)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Funds added successfully"})
}
