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
	FundCollection     *mongo.Collection
	ExpenseCollection  *mongo.Collection
	ReportCollection   *mongo.Collection
	WishlistCollection *mongo.Collection
)

func ConnectDB() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	uri := os.Getenv("DB_URI")
	if uri == "" {
		log.Fatal("DB_URI environment variable not set")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("error connecting to MongoDB:", err)
	}

	// Set a timeout for the ping
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("failed to connect to MongoDB:", err)
	}

	fmt.Println("connected to MongoDB")

	// Set the global Client reference
	Client = client

	UserCollection = Client.Database("expense_tracker").Collection("users")
	FundCollection = Client.Database("expense_tracker").Collection("funds")
	ExpenseCollection = Client.Database("expense_tracker").Collection("expenses")
	ReportCollection = Client.Database("expense_tracker").Collection("reports")
	WishlistCollection = Client.Database("expense_tracker").Collection("wishlists")

	return Client
}
