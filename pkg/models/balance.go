package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Balance struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Total_Balance float64            `json:"total_balance"`
	Updated_At    time.Time          `json:"updated_at"`
}
