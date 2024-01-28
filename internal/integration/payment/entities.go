package payment

type Payment struct {
	PaymentType int    `json:"payment_type"`
	OrderID     string `json:"order_id"`
}
