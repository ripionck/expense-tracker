package controllers

import (
	"context"
	"expense-tracker/pkg/models"
	"expense-tracker/pkg/utils"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

func AddFunds(c *gin.Context) {
	var fund models.Fund
	if err := c.BindJSON(&fund); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fund.Inserted_At = time.Now()
	fund.ID = primitive.NewObjectID()
	fund.Fund_ID = fund.ID.Hex()

	session, err := utils.Client.StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
		return
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Insert the fund
		_, err := utils.FundCollection.InsertOne(sessCtx, fund)
		if err != nil {
			return nil, err
		}

		// Check if a balance document exists
		var balance models.Balance
		err = utils.BalanceCollection.FindOne(sessCtx, bson.M{}).Decode(&balance)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// Initialize balance when no document exists
				newBalance := models.Balance{
					ID:            primitive.NewObjectID(),
					Total_Balance: roundToTwoDecimals(fund.Amount),
					Updated_At:    time.Now(),
				}
				_, err := utils.BalanceCollection.InsertOne(sessCtx, newBalance)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			// Update the existing balance
			newBalance := roundToTwoDecimals(balance.Total_Balance + fund.Amount)
			update := bson.M{
				"$set": bson.M{
					"total_balance": newBalance,
					"updated_at":    time.Now(),
				},
			}
			_, err := utils.BalanceCollection.UpdateOne(sessCtx, bson.M{"_id": balance.ID}, update)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Funds added successfully"})
}
