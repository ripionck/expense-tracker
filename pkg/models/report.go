package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	FundList    []Fund             `json:"fund_list"`
	ExpenseList []Expense          `json:"expense_list"`
	ReportID    string             `json:"report_id"`
	Created_At  time.Time          `json:"created_at"`
}
