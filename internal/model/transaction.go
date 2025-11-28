package model

import (
	"time"
)

type Status string

const (
	Pending   Status = "PENDING"
	Completed Status = "COMPLETED"
	Failed    Status = "FAILED"
	Refunded  Status = "REFUNDED"
)

type PaymentMethod string

const (
	Card       PaymentMethod = "PAYMENT_METHOD_CARD"
	OnDelivery PaymentMethod = "PAYMENT_METHOD_ON_DELIVERY"
)

type Currency string

const (
	EUR Currency = "EUR"
	USD Currency = "USD"
)

type Transaction struct {
	ID                   int64         `json:"id"`
	UserID               int64         `json:"user_id"`
	OrderID              *int64        `json:"order_id"`
	Amount               float64       `json:"amount"`
	Currency             Currency      `json:"currency"`
	Status               Status        `json:"Status"`
	GatewayTransactionID *string       `json:"gateway_transaction_id"`
	PaymentMethod        PaymentMethod `json:"payment_method"`
	CreatedAt            time.Time     `json:"created_at"`
}
