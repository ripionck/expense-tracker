package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fund struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Type        string             `json:"type" validate:"required"`
	Amount      float64            `json:"amount" validate:"required"`
	Note        string             `json:"note" validate:"required,min=10,max=120"`
	Inserted_At time.Time          `json:"inserted_at"`
	Fund_ID     string             `json:"fund_id"`
}
