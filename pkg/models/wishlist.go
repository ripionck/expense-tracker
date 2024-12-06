package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wishlist struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	ItemName   string             `json:"item_name"`
	Price      float64            `json:"price,omitempty"`
	Priority   int                `json:"priority,omitempty"`
	Note       string             `json:"note" validate:"required,min=10,max=120"`
	Created_At time.Time          `json:"created_at"`
	Updated_At time.Time          `json:"updated_at"`
}
