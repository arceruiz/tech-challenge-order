package service_test

import (
	"context"
	"errors"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/integration/payment"
	mock_test "tech-challenge-order/internal/mocks"
	"tech-challenge-order/internal/repository"
	"tech-challenge-order/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undefinedlabs/go-mpatch"
)

func TestOrderService_GetByID(t *testing.T) {

	type Given struct {
		id        string
		orderRepo func() repository.OrderRepository
	}
	type Expected struct {
		err assert.ErrorAssertionFunc
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{

		"given order with main fields filled, must return created paymend with all fields filled": {
			given: Given{
				id: "1234",
				orderRepo: func() repository.OrderRepository {
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, "1234").Return(&canonical.Order{
						ID:         "order_valid_id",
						CustomerID: "order_valid_customer_id",
						Status:     canonical.ORDER_CHECKED_OUT,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						Total:      1000,
						OrderItems: []canonical.OrderItem{
							{
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
								Quantity: 10,
							},
							{
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
								Quantity: 10,
							},
						},
					}, nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.NoError,
			},
		},
	}

	for _, tc := range tests {
		_, err := service.NewOrderService(tc.given.orderRepo(), &mock_test.PaymentServiceMock{}).GetByID(context.Background(), tc.given.id)

		tc.expected.err(t, err)
	}
}

func TestOrderService_GetAll(t *testing.T) {

	type Given struct {
		orderRepo func() repository.OrderRepository
	}
	type Expected struct {
		err assert.ErrorAssertionFunc
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{

		"given order with main fields filled, must return created paymend with all fields filled": {
			given: Given{
				orderRepo: func() repository.OrderRepository {
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("GetAll", mock.Anything).Return([]canonical.Order{
						{
							ID:         "order_valid_id",
							CustomerID: "order_valid_customer_id",
							Status:     canonical.ORDER_CHECKED_OUT,
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
							Total:      1000,
							OrderItems: []canonical.OrderItem{
								{
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
									Quantity: 10,
								},
								{
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
									Quantity: 10,
								},
							},
						},
						{
							ID:         "order_valid_id",
							CustomerID: "order_valid_customer_id",
							Status:     canonical.ORDER_CHECKED_OUT,
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
							Total:      1000,
							OrderItems: []canonical.OrderItem{
								{
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
									Quantity: 10,
								},
								{
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
									Quantity: 10,
								},
							},
						},
					}, nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.NoError,
			},
		},
	}

	for _, tc := range tests {
		_, err := service.NewOrderService(tc.given.orderRepo(), &mock_test.PaymentServiceMock{}).GetAll(context.Background())

		tc.expected.err(t, err)
	}
}

func TestOrderService_GetByCategory(t *testing.T) {

	type Given struct {
		status    canonical.OrderStatus
		orderRepo func() repository.OrderRepository
	}
	type Expected struct {
		err assert.ErrorAssertionFunc
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{

		"given order with main fields filled, must return created paymend with all fields filled": {
			given: Given{
				status: canonical.ORDER_CHECKED_OUT,
				orderRepo: func() repository.OrderRepository {
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("GetByStatus", mock.Anything, 1).Return([]canonical.Order{
						{
							ID:         "order_valid_id",
							CustomerID: "order_valid_customer_id",
							Status:     canonical.ORDER_CHECKED_OUT,
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
							Total:      1000,
							OrderItems: []canonical.OrderItem{
								{
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
									Quantity: 10,
								},
								{
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
									Quantity: 10,
								},
							},
						},
						{
							ID:         "order_valid_id",
							CustomerID: "order_valid_customer_id",
							Status:     canonical.ORDER_CHECKED_OUT,
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
							Total:      1000,
							OrderItems: []canonical.OrderItem{
								{
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
									Quantity: 10,
								},
								{
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
									Quantity: 10,
								},
							},
						},
					}, nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.NoError,
			},
		},
	}

	for _, tc := range tests {
		_, err := service.NewOrderService(tc.given.orderRepo(), &mock_test.PaymentServiceMock{}).GetByStatus(context.Background(), tc.given.status)

		tc.expected.err(t, err)
	}
}

func TestOrderService_Create(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})
	mpatch.PatchMethod(canonical.NewUUID, func() string {
		return "order_valid_id"
	})

	type Given struct {
		order     canonical.Order
		orderRepo func() repository.OrderRepository
	}
	type Expected struct {
		err assert.ErrorAssertionFunc
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given order with main fields filled, must return created paymend with all fields filled": {
			given: Given{
				order: canonical.Order{
					ID:         "order_valid_id",
					CustomerID: "order_valid_customer_id",
					Status:     canonical.ORDER_RECEIVED,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Total:      1000,
					OrderItems: []canonical.OrderItem{
						{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						},
						{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						},
					},
				},
				orderRepo: func() repository.OrderRepository {
					order := canonical.Order{
						ID:         "order_valid_id",
						CustomerID: "order_valid_customer_id",
						Status:     canonical.ORDER_RECEIVED,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						Total:      2000,
						OrderItems: []canonical.OrderItem{
							{
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
								Quantity: 10,
							},
							{
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
								Quantity: 10,
							},
						},
					}
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("Create", mock.Anything, order).Return(&order, nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.NoError,
			},
		},
		"given error creating, must return error": {
			given: Given{
				order: canonical.Order{
					CustomerID: "order_valid_customer_id",
					Status:     canonical.ORDER_RECEIVED,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Total:      2000,
					OrderItems: []canonical.OrderItem{
						{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						},
						{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						},
					},
				},
				orderRepo: func() repository.OrderRepository {
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("Create", mock.Anything, mock.Anything).Return(&canonical.Order{}, errors.New("error creating order"))
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for _, tc := range tests {
		err := service.NewOrderService(tc.given.orderRepo(), &mock_test.PaymentServiceMock{}).Create(context.Background(), tc.given.order)

		tc.expected.err(t, err)
	}
}

func TestOrderService_Update(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})
	mpatch.PatchMethod(canonical.NewUUID, func() string {
		return "order_valid_id"
	})

	type Given struct {
		order     canonical.Order
		orderID   string
		orderRepo func() repository.OrderRepository
	}
	type Expected struct {
		err assert.ErrorAssertionFunc
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given order with main fields filled, must return created paymend with all fields filled": {
			given: Given{
				orderID: "order_valid_id",
				order: canonical.Order{
					ID:         "order_valid_id",
					CustomerID: "order_valid_customer_id",
					Status:     canonical.ORDER_CHECKED_OUT,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Total:      1000,
					OrderItems: []canonical.OrderItem{
						{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						},
						{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						},
					},
				},
				orderRepo: func() repository.OrderRepository {
					order := canonical.Order{
						ID:         "order_valid_id",
						CustomerID: "order_valid_customer_id",
						Status:     canonical.ORDER_CHECKED_OUT,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						Total:      2000,
						OrderItems: []canonical.OrderItem{
							{
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
								Quantity: 10,
							},
							{
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
								Quantity: 10,
							},
						},
					}
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("Update", mock.Anything, "order_valid_id", order).Return(nil)
					return repoMock
				},
			},
			expected: Expected{
				err: assert.NoError,
			},
		},
		"given error creating, must return error": {
			given: Given{
				orderID: "order_valid_id",
				order: canonical.Order{
					ID:         "order_valid_id",
					CustomerID: "order_valid_customer_id",
					Status:     canonical.ORDER_CHECKED_OUT,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Total:      1000,
					OrderItems: []canonical.OrderItem{
						{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						},
						{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						},
					},
				},
				orderRepo: func() repository.OrderRepository {
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error creating order"))
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for _, tc := range tests {
		err := service.NewOrderService(tc.given.orderRepo(), &mock_test.PaymentServiceMock{}).Update(context.Background(), tc.given.orderID, tc.given.order)

		tc.expected.err(t, err)
	}
}

func TestOrderService_Checkout(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})
	mpatch.PatchMethod(canonical.NewUUID, func() string {
		return "order_valid_id"
	})

	type Given struct {
		orderID        string
		orderRepo      func() repository.OrderRepository
		paymentService func() payment.PaymentService
	}
	type Expected struct {
		err assert.ErrorAssertionFunc
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given order with main fields filled, must return created paymend with all fields filled": {
			given: Given{
				orderID: "order_valid_id",
				orderRepo: func() repository.OrderRepository {
					order := canonical.Order{
						ID:         "order_valid_id",
						CustomerID: "order_valid_customer_id",
						Status:     canonical.ORDER_CHECKED_OUT,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						Total:      1000,
						OrderItems: []canonical.OrderItem{
							{
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
								Quantity: 10,
							},
							{
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
								Quantity: 10,
							},
						},
					}
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, "order_valid_id").Return(&order, nil)
					repoMock.On("Update", mock.Anything, "order_valid_id", order).Return(nil)
					return repoMock
				},
				paymentService: func() payment.PaymentService {
					paymentServiceMock := &mock_test.PaymentServiceMock{}
					paymentServiceMock.On("Create", mock.Anything, mock.Anything).Return(nil)
					return paymentServiceMock
				},
			},
			expected: Expected{
				err: assert.NoError,
			},
		},
		"given error creating, must return error": {
			given: Given{
				orderID: "order_valid_id",
				orderRepo: func() repository.OrderRepository {
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, mock.Anything).Return(&canonical.Order{}, nil)
					repoMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error creating order"))
					return repoMock
				},
				paymentService: func() payment.PaymentService {
					paymentServiceMock := &mock_test.PaymentServiceMock{}
					paymentServiceMock.On("Create", mock.Anything, mock.Anything).Return(nil)
					return paymentServiceMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
		"given error creating, must returnasd error": {
			given: Given{
				orderID: "order_valid_id",
				orderRepo: func() repository.OrderRepository {
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, mock.Anything).Return(nil, errors.New("error creating order"))
					return repoMock
				},
				paymentService: func() payment.PaymentService {
					paymentServiceMock := &mock_test.PaymentServiceMock{}
					paymentServiceMock.On("Create", mock.Anything, mock.Anything).Return(nil)
					return paymentServiceMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
		"given error creating, must retasdasdurnasd error": {
			given: Given{
				orderID: "order_valid_id",
				orderRepo: func() repository.OrderRepository {
					repoMock := &mock_test.OrderRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, mock.Anything).Return(&canonical.Order{}, nil)
					repoMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(&canonical.Order{}, nil)
					return repoMock
				},
				paymentService: func() payment.PaymentService {
					paymentServiceMock := &mock_test.PaymentServiceMock{}
					paymentServiceMock.On("Create", mock.Anything, mock.Anything).Return(errors.New("error creating order"))
					return paymentServiceMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for _, tc := range tests {
		order, err := service.NewOrderService(tc.given.orderRepo(), tc.given.paymentService()).CheckoutOrder(context.Background(), tc.given.orderID)

		if err == nil {
			assert.Equal(t, canonical.ORDER_CHECKED_OUT, order.Status)
		}
		tc.expected.err(t, err)
	}
}
