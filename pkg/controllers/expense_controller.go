package controllers

import (
	"context"
	"expense-tracker/pkg/models"
	"expense-tracker/pkg/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	session, err := utils.Client.StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
		return
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Insert the expense
		_, err := utils.ExpenseCollection.InsertOne(sessCtx, expense)
		if err != nil {
			return nil, err
		}

		// Decrease balance
		newBalance := roundToTwoDecimals(-expense.Amount)
		update := bson.M{
			"$inc": bson.M{"total_balance": newBalance},
			"$set": bson.M{"updated_at": time.Now()},
		}
		result, err := utils.BalanceCollection.UpdateOne(sessCtx, bson.M{}, update)
		if err != nil {
			return nil, err
		}

		// Ensure a balance document exists
		if result.MatchedCount == 0 {
			return nil, fmt.Errorf("balance document not found")
		}

		return nil, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense added successfully"})
}
