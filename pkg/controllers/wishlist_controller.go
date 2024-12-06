package controllers

import (
	"context"
	"expense-tracker/pkg/models"
	"expense-tracker/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateWishlist adds a new wishlist item to the database.
func CreateWishlist(c *gin.Context) {
	var wishlist models.Wishlist

	// Bind JSON data to the wishlist struct
	if err := c.ShouldBindJSON(&wishlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Set the Created_At and Updated_At timestamps
	wishlist.Created_At = time.Now()
	wishlist.Updated_At = time.Now()

	// Generate a new ID for the wishlist item
	wishlist.ID = primitive.NewObjectID()

	// Insert the new wishlist item into the collection
	_, err := utils.WishlistCollection.InsertOne(context.TODO(), wishlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wishlist item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Wishlist item created successfully",
		"wishlist": wishlist})
}

// UpdateWishlist updates an existing wishlist item.
func UpdateWishlist(c *gin.Context) {
	// Get the wishlist item ID from the URL parameter
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist item ID"})
		return
	}

	var wishlist models.Wishlist
	// Bind JSON data to the wishlist struct
	if err := c.ShouldBindJSON(&wishlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Set the Updated_At timestamp
	wishlist.Updated_At = time.Now()

	// Prepare the update query
	update := bson.M{
		"$set": bson.M{
			"item_name":  wishlist.Item_Name,
			"price":      wishlist.Price,
			"priority":   wishlist.Priority,
			"note":       wishlist.Note,
			"updated_at": wishlist.Updated_At,
		},
	}

	// Update the wishlist item in the collection
	_, err = utils.WishlistCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wishlist item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wishlist item updated successfully"})
}

// DeleteWishlist deletes a wishlist item.
func DeleteWishlist(c *gin.Context) {
	// Get the wishlist item ID from the URL parameter
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist item ID"})
		return
	}

	// Delete the wishlist item from the collection
	_, err = utils.WishlistCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete wishlist item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wishlist item deleted successfully"})
}

// GetWishlist retrieves all wishlist items or a specific item.
func GetWishlist(c *gin.Context) {
	id := c.Param("id")

	// If the ID is provided, fetch the specific item
	if id != "" {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist item ID"})
			return
		}

		var wishlist models.Wishlist
		err = utils.WishlistCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&wishlist)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"message": "Wishlist item not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist item"})
			}
			return
		}

		c.JSON(http.StatusOK, wishlist)
		return
	}

	// If no ID is provided, fetch all wishlist items
	cursor, err := utils.WishlistCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist items"})
		return
	}
	defer cursor.Close(context.TODO())

	var wishlists []models.Wishlist
	for cursor.Next(context.TODO()) {
		var wishlist models.Wishlist
		if err := cursor.Decode(&wishlist); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode wishlist item"})
			return
		}
		wishlists = append(wishlists, wishlist)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while iterating through wishlist items"})
		return
	}

	c.JSON(http.StatusOK, wishlists)
}
