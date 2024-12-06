package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Username      string             `json:"username" validate:"required,min=2,max=100"`
	Name          string             `json:"name" validate:"required,min=2,max=100"`
	Email         string             `json:"email" validate:"required,email"`
	Password      string             `json:"password" validate:"required,min=6"`
	Access_Token  string             `json:"access_token,omitempty"`
	Refresh_Token string             `json:"refresh_token,omitempty"`
	Created_At    time.Time          `json:"created_at"`
	Updated_At    time.Time          `json:"updated_at"`
	User_ID       string             `json:"user_id"`
}
