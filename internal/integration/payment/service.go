package payment

type PaymentService interface {
	Create(Payment) error
}
