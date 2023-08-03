package repository

import (
	"context"
	"math/rand"
	"reflect"
	"testing"
)

func TestSavePayment(t *testing.T) {
	repository, err := NewDatabaseRepository("root:root@tcp(localhost:3306)/mekano_test")
	if err != nil {
		t.Fatalf("Error initializing the database")
	}

	payment := Payment{
		Consecutive: rand.Intn(100),
		CreateAt:    "2023-01-13",
		FileName:    "payment_test.xlsx",
	}

	err = repository.SavePayment(context.Background(), payment)
	if err != nil {
		t.Fatalf("Error saving payment: %v", err)
	}
}
func TestSaveBilling(t *testing.T) {
	repository, err := NewDatabaseRepository("root:root@tcp(localhost:3306)/mekano_test")
	if err != nil {
		t.Fatalf("Error initializing the database")
	}

	billing := Billing{
		Debit:    100,
		Credit:   200,
		Base:     300,
		FileName: "billing_test.xlsx",
		CreateAt: "2023-01-13",
	}

	err = repository.SaveBilling(context.Background(), billing)
	if err != nil {
		t.Fatalf("Error saving billing: %v", err)
	}
}

func TestGetPayment(t *testing.T) {
	repository, err := NewDatabaseRepository("root:root@tcp(localhost:3306)/mekano_test")
	if err != nil {
		t.Fatalf("Error initializing the database")
	}

	payments, err := repository.GetPayment(context.Background())
	if err != nil {
		t.Fatalf("Error getting payment: %v", err)
	}

	expectedData := Payment{
		Consecutive: payments.Consecutive,
		CreateAt:    "2023-01-13",
		FileName:    "payment_test.xlsx",
	}

	if !reflect.DeepEqual(payments, expectedData) {
		t.Errorf("Expected data: %+v, got: %+v", expectedData, payments)
	}

}
