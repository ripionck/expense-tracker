package controllers

import (
	"context"
	"expense-tracker/pkg/helpers.go"
	"expense-tracker/pkg/models"
	"expense-tracker/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := utils.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error detected while fetching the email."})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists."})
			return
		}

		hashedPassword, err := utils.HashedPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password."})
			return
		}
		user.Password = hashedPassword

		user.Created_At = time.Now()
		user.Updated_At = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		access_token, refresh_token, err := helpers.GenerateTokens(user.Email, user.Username, user.Name, user.User_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating tokens."})
			return
		}
		user.Access_Token = access_token
		user.Refresh_Token = refresh_token

		_, err = utils.UserCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user details were not saved"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user registration successful"})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := utils.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := utils.VerifyPassword(foundUser.Password, user.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		responseUser := gin.H{
			"email":      foundUser.Email,
			"username":   foundUser.Username,
			"name":       foundUser.Name,
			"user_id":    foundUser.User_ID,
			"created_at": foundUser.Created_At,
			"updated_at": foundUser.Updated_At,
		}

		// Send response
		c.JSON(http.StatusOK, gin.H{
			"access_token":  foundUser.Access_Token,
			"refresh_token": foundUser.Refresh_Token,
			"user":          responseUser,
		})
	}
}
