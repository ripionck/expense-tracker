package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client             *mongo.Client
	UserCollection     *mongo.Collection
	BalanceCollection  *mongo.Collection
	FundCollection     *mongo.Collection
	ExpenseCollection  *mongo.Collection
	ReportCollection   *mongo.Collection
	WishlistCollection *mongo.Collection
)

func ConnectDB() *mongo.Client {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not loaded. Ensure environment variables are set.")
	}

	// Get DB URI
	uri := os.Getenv("DB_URI")
	if uri == "" {
		log.Fatal("Error: DB_URI is not set in environment variables or .env file")
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Ping MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")

	// Set global variables
	Client = client
	UserCollection = Client.Database("expense_tracker").Collection("users")
	BalanceCollection = Client.Database("expense_tracker").Collection("balance")
	FundCollection = Client.Database("expense_tracker").Collection("funds")
	ExpenseCollection = Client.Database("expense_tracker").Collection("expenses")
	ReportCollection = Client.Database("expense_tracker").Collection("reports")
	WishlistCollection = Client.Database("expense_tracker").Collection("wishlists")

	return Client
}
