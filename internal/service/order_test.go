package service

import (
	"context"
	"errors"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/integration/product"
	"tech-challenge-order/internal/integration/sqs_publisher"
	"tech-challenge-order/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
					repoMock := &OrderRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, "1234").Return(&canonical.Order{
						ID:         "order_valid_id",
						CustomerID: "order_valid_customer_id",
						Status:     canonical.ORDER_RECEIVED,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						Total:      1000,
						OrderItems: map[string]*canonical.OrderItem{
							"product_valid_id": {
								Quantity: 10,
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
							},
							"product_valid_id1": {
								Quantity: 10,
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
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
		order := orderService{
			repo: tc.given.orderRepo(),
		}
		_, err := order.GetByID(context.Background(), tc.given.id)

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
					repoMock := &OrderRepositoryMock{}
					repoMock.On("GetAll", mock.Anything).Return([]canonical.Order{
						{
							ID:         "order_valid_id",
							CustomerID: "order_valid_customer_id",
							Status:     canonical.ORDER_RECEIVED,
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
							Total:      1000,
							OrderItems: map[string]*canonical.OrderItem{
								"product_valid_id": {
									Quantity: 10,
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
								},
								"product_valid_id1": {
									Quantity: 10,
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
								},
							},
						},
						{
							ID:         "order_valid_id",
							CustomerID: "order_valid_customer_id",
							Status:     canonical.ORDER_RECEIVED,
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
							Total:      1000,
							OrderItems: map[string]*canonical.OrderItem{
								"product_valid_id": {
									Quantity: 10,
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
								},
								"product_valid_id1": {
									Quantity: 10,
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
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
		order := orderService{
			repo: tc.given.orderRepo(),
		}
		_, err := order.GetAll(context.Background())

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
				status: canonical.ORDER_RECEIVED,
				orderRepo: func() repository.OrderRepository {
					repoMock := &OrderRepositoryMock{}
					repoMock.On("GetByStatus", mock.Anything, 0).Return([]canonical.Order{
						{
							ID:         "order_valid_id",
							CustomerID: "order_valid_customer_id",
							Status:     canonical.ORDER_RECEIVED,
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
							Total:      1000,
							OrderItems: map[string]*canonical.OrderItem{
								"product_valid_id": {
									Quantity: 10,
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
								},
								"product_valid_id1": {
									Quantity: 10,
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
								},
							},
						},
						{
							ID:         "order_valid_id",
							CustomerID: "order_valid_customer_id",
							Status:     canonical.ORDER_RECEIVED,
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
							Total:      1000,
							OrderItems: map[string]*canonical.OrderItem{
								"product_valid_id": {
									Quantity: 10,
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
								},
								"product_valid_id1": {
									Quantity: 10,
									Product: canonical.Product{
										ID:       "product_valid_id",
										Name:     "product_valid_name",
										Price:    50,
										Category: "product_valid_category",
									},
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
		order := orderService{
			repo: tc.given.orderRepo(),
		}
		_, err := order.GetByStatus(context.Background(), tc.given.status)

		tc.expected.err(t, err)
	}
}

func TestOrderService_Create(t *testing.T) {
	type Given struct {
		order          canonical.Order
		orderRepo      func() repository.OrderRepository
		productService func() product.ProductService
		publisher      func() sqs_publisher.Publisher
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
					OrderItems: map[string]*canonical.OrderItem{
						"product_valid_id": {
							Quantity: 10,
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
						},
						"product_valid_id1": {
							Quantity: 10,
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
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
						OrderItems: map[string]*canonical.OrderItem{
							"product_valid_id": {
								Quantity: 10,
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
							},
							"product_valid_id1": {
								Quantity: 10,
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
							},
						},
					}
					repoMock := &OrderRepositoryMock{}
					repoMock.On("Create", mock.MatchedBy(func(o canonical.Order) bool {
						return o.CustomerID == order.CustomerID
					})).Return(&canonical.Order{}, nil)

					return repoMock
				},
				productService: func() product.ProductService {
					pMock := &ProductMock{}
					pMock.MockGetProducts("product_valid_id", func(args mock.Arguments) {
						products := args.Get(0).(map[string]*canonical.OrderItem)

						products["product_valid_id"] = &canonical.OrderItem{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						}
						products["product_valid_id1"] = &canonical.OrderItem{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						}
					})

					return pMock
				},
				publisher: func() sqs_publisher.Publisher {
					publisherMock := new(PublisherMock)

					publisherMock.On("SendMessage").Return(nil)

					return publisherMock
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
					OrderItems: map[string]*canonical.OrderItem{
						"product_valid_id": {
							Quantity: 10,
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
						},
						"product_valid_id1": {
							Quantity: 10,
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
						},
					},
				},
				orderRepo: func() repository.OrderRepository {
					repoMock := &OrderRepositoryMock{}

					repoMock.On("Create", mock.MatchedBy(func(o canonical.Order) bool {
						return o.CustomerID == "order_valid_customer_id"
					})).Return(&canonical.Order{}, errors.New("generic error"))

					return repoMock
				},
				productService: func() product.ProductService {
					pMock := &ProductMock{}
					pMock.MockGetProducts("product_valid_id", func(args mock.Arguments) {
						products := args.Get(0).(map[string]*canonical.OrderItem)

						products["product_valid_id"] = &canonical.OrderItem{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						}
						products["product_valid_id1"] = &canonical.OrderItem{
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
							Quantity: 10,
						}
					})

					return pMock
				},
				publisher: func() sqs_publisher.Publisher {
					return new(PublisherMock)
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for _, tc := range tests {
		order := orderService{
			repo:           tc.given.orderRepo(),
			productService: tc.given.productService(),
			publisher:      tc.given.publisher(),
		}

		err := order.Create(context.Background(), tc.given.order)

		tc.expected.err(t, err)
	}
}

func TestOrderService_Update(t *testing.T) {
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
					Status:     canonical.ORDER_RECEIVED,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Total:      1000,
					OrderItems: map[string]*canonical.OrderItem{
						"product_valid_id": {
							Quantity: 10,
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
						},
						"product_valid_id1": {
							Quantity: 10,
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
						},
					},
				},
				orderRepo: func() repository.OrderRepository {
					order := canonical.Order{
						ID: "order_valid_id",
					}
					repoMock := &OrderRepositoryMock{}
					repoMock.On("Update", order.ID).Return(nil)
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
					Status:     canonical.ORDER_RECEIVED,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Total:      1000,
					OrderItems: map[string]*canonical.OrderItem{
						"product_valid_id": {
							Quantity: 10,
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
						},
						"product_valid_id1": {
							Quantity: 10,
							Product: canonical.Product{
								ID:       "product_valid_id",
								Name:     "product_valid_name",
								Price:    50,
								Category: "product_valid_category",
							},
						},
					},
				},
				orderRepo: func() repository.OrderRepository {
					repoMock := &OrderRepositoryMock{}
					repoMock.On("Update", mock.Anything).Return(errors.New("error creating order"))
					return repoMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for _, tc := range tests {
		order := orderService{
			repo: tc.given.orderRepo(),
		}
		err := order.Update(context.Background(), tc.given.orderID, tc.given.order)

		tc.expected.err(t, err)
	}
}

func TestOrderService_Checkout(t *testing.T) {
	type Given struct {
		orderID   string
		orderRepo func() repository.OrderRepository
		publisher func() sqs_publisher.Publisher
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
						Status:     canonical.ORDER_RECEIVED,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						Total:      1000,
						OrderItems: map[string]*canonical.OrderItem{
							"product_valid_id": {
								Quantity: 10,
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
							},
							"product_valid_id1": {
								Quantity: 10,
								Product: canonical.Product{
									ID:       "product_valid_id",
									Name:     "product_valid_name",
									Price:    50,
									Category: "product_valid_category",
								},
							},
						},
					}
					repoMock := &OrderRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, "order_valid_id").Return(&order, nil)
					repoMock.On("UpdateStatus", "order_valid_id").Return(nil)
					return repoMock
				},
				publisher: func() sqs_publisher.Publisher {
					publisherMock := &PublisherMock{}
					publisherMock.On("SendMessage").Return(nil)
					return publisherMock
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
					repoMock := &OrderRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, mock.Anything).Return(&canonical.Order{}, nil)
					repoMock.On("UpdateStatus", "order_valid_id").Return(errors.New("error creating order"))
					return repoMock
				},
				publisher: func() sqs_publisher.Publisher {
					publisherMock := &PublisherMock{}
					publisherMock.On("SendMessage").Return(nil)
					return publisherMock
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
					repoMock := &OrderRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, mock.Anything).Return(nil, errors.New("error creating order"))
					return repoMock
				},
				publisher: func() sqs_publisher.Publisher {
					publisherMock := &PublisherMock{}
					publisherMock.On("SendMessage").Return(nil)
					return publisherMock
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
					repoMock := &OrderRepositoryMock{}
					repoMock.On("GetByID", mock.Anything, mock.Anything).Return(&canonical.Order{}, nil)
					repoMock.On("UpdateStatus", "order_valid_id").Return(errors.New("error creating order"))
					return repoMock
				},
				publisher: func() sqs_publisher.Publisher {
					publisherMock := &PublisherMock{}
					publisherMock.On("SendMessage").Return(nil)
					return publisherMock
				},
			},
			expected: Expected{
				err: assert.Error,
			},
		},
	}

	for _, tc := range tests {
		ordersvc := orderService{
			repo:      tc.given.orderRepo(),
			publisher: tc.given.publisher(),
		}

		order, err := ordersvc.CheckoutOrder(context.Background(), tc.given.orderID)

		if err == nil {
			assert.Equal(t, canonical.ORDER_RECEIVED, int(order.Status))
		}
		tc.expected.err(t, err)
	}
}

func TestUpdateStatus(t *testing.T) {
	mockRepo := new(OrderRepositoryMock)

	order := canonical.Order{
		ID: "fakeId",
	}

	mockRepo.On("GetByID", mock.Anything, order.ID).Return(&canonical.Order{}, nil)
	mockRepo.On("UpdateStatus", order.ID).Return(nil)

	svc := orderService{
		repo: mockRepo,
	}

	err := svc.UpdateStatus(context.Background(), order.ID, canonical.ORDER_COMPLETED)

	assert.Nil(t, err)
	mockRepo.AssertExpectations(t)
}
