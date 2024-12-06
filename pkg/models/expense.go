package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Type        string             `json:"type" validate:"required,"`
	Amount      float64            `json:"amount" validate:"required,gte=5,lte=100000"`
	Note        string             `json:"note" validate:"required,min=10,max=120"`
	Inserted_At time.Time          `json:"inserted_at"`
	Expense_ID  string             `json:"expense_id"`
}
