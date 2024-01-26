package repository_test

import (
	"context"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/undefinedlabs/go-mpatch"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestOrderRepository_GetByID(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})

	type Given struct {
		mtestFunc func(mt *mtest.T)
	}
	type Expected struct {
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given valid search result, must return valid order": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := repository.NewOrderRepo(mt.DB)
					mt.AddMockResponses(
						mtest.CreateCursorResponse(1,
							"order.order",
							mtest.FirstBatch,
							bson.D{
								{Key: "_id", Value: "order_valid_id"},
								{Key: "customer_id", Value: "order_valid_customer_id"},
								{Key: "status", Value: canonical.ORDER_RECEIVED},
								{Key: "created_at", Value: time.Now()},
								{Key: "updated_at", Value: time.Now()},
								{Key: "total", Value: 0.0},
								{Key: "order_items", Value: []bson.D{
									{
										{Key: "quantity", Value: 5},
										{Key: "product_id", Value: "product_id"},
										{Key: "name", Value: "product_name"},
										{Key: "price", Value: 10.0},
										{Key: "category", Value: "product_category"},
									},
									{
										{Key: "quantity", Value: 10},
										{Key: "product_id", Value: "product_id"},
										{Key: "name", Value: "product_name"},
										{Key: "price", Value: 10.0},
										{Key: "category", Value: "product_category"},
									},
								},
								},
							},
						),
					)
					order, err := repo.GetByID(context.Background(), "order_valid_id")
					assert.Nil(t, err)
					assert.Equal(t, order.ID, "order_valid_id")
					assert.Equal(t, order.Status, canonical.ORDER_RECEIVED)
				},
			},
		},
		"given entity not found must return error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := repository.NewOrderRepo(mt.DB)
					mt.AddMockResponses(mtest.CreateCursorResponse(0, "order.order", mtest.FirstBatch))
					order, err := repo.GetByID(context.Background(), "asd")
					assert.NotNil(t, err)
					assert.Equal(t, err.Error(), "mongo: no documents in result")
					assert.Nil(t, order)
				},
			},
		},
	}

	for _, tc := range tests {
		db := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		db.Run("", tc.given.mtestFunc)
	}
}

func TestOrderRepository_GetAll(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})

	type Given struct {
		mtestFunc func(mt *mtest.T)
	}
	type Expected struct {
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given valid search result, must return valid order": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := repository.NewOrderRepo(mt.DB)

					first := mtest.CreateCursorResponse(1, "order.order", mtest.FirstBatch, bson.D{
						{Key: "_id", Value: "order_valid_id"},
						{Key: "customer_id", Value: "order_valid_customer_id"},
						{Key: "status", Value: canonical.ORDER_RECEIVED},
						{Key: "created_at", Value: time.Now()},
						{Key: "updated_at", Value: time.Now()},
						{Key: "total", Value: 0.0},
						{Key: "order_items", Value: []bson.D{
							{
								{Key: "quantity", Value: 5},
								{Key: "product_id", Value: "product_id"},
								{Key: "name", Value: "product_name"},
								{Key: "price", Value: 10.0},
								{Key: "category", Value: "product_category"},
							},
							{
								{Key: "quantity", Value: 10},
								{Key: "product_id", Value: "product_id"},
								{Key: "name", Value: "product_name"},
								{Key: "price", Value: 10.0},
								{Key: "category", Value: "product_category"},
							},
						},
						},
					})
					getMore := mtest.CreateCursorResponse(1, "order.order", mtest.NextBatch, bson.D{
						{Key: "_id", Value: "order_valid_id"},
						{Key: "customer_id", Value: "order_valid_customer_id"},
						{Key: "status", Value: canonical.ORDER_RECEIVED},
						{Key: "created_at", Value: time.Now()},
						{Key: "updated_at", Value: time.Now()},
						{Key: "total", Value: 0.0},
						{Key: "order_items", Value: []bson.D{
							{
								{Key: "quantity", Value: 5},
								{Key: "product_id", Value: "product_id"},
								{Key: "name", Value: "product_name"},
								{Key: "price", Value: 10.0},
								{Key: "category", Value: "product_category"},
							},
							{
								{Key: "quantity", Value: 10},
								{Key: "product_id", Value: "product_id"},
								{Key: "name", Value: "product_name"},
								{Key: "price", Value: 10.0},
								{Key: "category", Value: "product_category"},
							},
						},
						},
					})
					lastCursor := mtest.CreateCursorResponse(0, "order.order", mtest.NextBatch)
					mt.AddMockResponses(first, getMore, lastCursor)

					orders, err := repo.GetAll(context.Background())
					assert.Nil(t, err)
					for _, order := range orders {
						assert.Equal(t, order.ID, "order_valid_id")
						assert.Equal(t, order.Status, canonical.ORDER_RECEIVED)
					}
				},
			},
		},
		"given entity not found must return error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := repository.NewOrderRepo(mt.DB)
					mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Message: "mongo: no documents in result"}))
					order, err := repo.GetAll(context.Background())
					assert.NotNil(t, err)
					assert.Equal(t, err.Error(), "write command error: [{write errors: [{mongo: no documents in result}]}, {<nil>}]")
					assert.Nil(t, order)
				},
			},
		},
	}

	for _, tc := range tests {
		db := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		db.Run("", tc.given.mtestFunc)
	}
}

func TestOrderRepository_GetByCategory(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})

	type Given struct {
		mtestFunc func(mt *mtest.T)
	}
	type Expected struct {
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given valid search result, must return valid order": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := repository.NewOrderRepo(mt.DB)

					first := mtest.CreateCursorResponse(1, "order.order", mtest.FirstBatch, bson.D{
						{Key: "_id", Value: "order_valid_id"},
						{Key: "customer_id", Value: "order_valid_customer_id"},
						{Key: "status", Value: canonical.ORDER_PREPARING},
						{Key: "created_at", Value: time.Now()},
						{Key: "updated_at", Value: time.Now()},
						{Key: "total", Value: 0.0},
						{Key: "order_items", Value: []bson.D{
							{
								{Key: "quantity", Value: 5},
								{Key: "product_id", Value: "product_id"},
								{Key: "name", Value: "product_name"},
								{Key: "price", Value: 10.0},
								{Key: "category", Value: "product_category"},
							},
							{
								{Key: "quantity", Value: 10},
								{Key: "product_id", Value: "product_id"},
								{Key: "name", Value: "product_name"},
								{Key: "price", Value: 10.0},
								{Key: "category", Value: "product_category"},
							},
						},
						},
					})
					getMore := mtest.CreateCursorResponse(1, "order.order", mtest.NextBatch, bson.D{
						{Key: "_id", Value: "order_valid_id"},
						{Key: "customer_id", Value: "order_valid_customer_id"},
						{Key: "status", Value: canonical.ORDER_PREPARING},
						{Key: "created_at", Value: time.Now()},
						{Key: "updated_at", Value: time.Now()},
						{Key: "total", Value: 0.0},
						{Key: "order_items", Value: []bson.D{
							{
								{Key: "quantity", Value: 5},
								{Key: "product_id", Value: "product_id"},
								{Key: "name", Value: "product_name"},
								{Key: "price", Value: 10.0},
								{Key: "category", Value: "product_category"},
							},
							{
								{Key: "quantity", Value: 10},
								{Key: "product_id", Value: "product_id"},
								{Key: "name", Value: "product_name"},
								{Key: "price", Value: 10.0},
								{Key: "category", Value: "product_category"},
							},
						},
						},
					})
					lastCursor := mtest.CreateCursorResponse(0, "order.order", mtest.NextBatch)
					mt.AddMockResponses(first, getMore, lastCursor)

					orders, err := repo.GetByStatus(context.Background(), 0)
					for _, order := range orders {
						assert.Nil(t, err)
						assert.Equal(t, order.Status, canonical.ORDER_PREPARING)
					}
				},
			},
		},
		"given entity not found must return error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := repository.NewOrderRepo(mt.DB)
					mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Message: "mongo: no documents in result"}))
					order, err := repo.GetByStatus(context.Background(), 2)
					assert.NotNil(t, err)
					assert.Equal(t, err.Error(), "write command error: [{write errors: [{mongo: no documents in result}]}, {<nil>}]")
					assert.Nil(t, order)
				},
			},
		},
	}

	for _, tc := range tests {
		db := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		db.Run("", tc.given.mtestFunc)
	}
}

func TestCreate(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})

	type Given struct {
		mtestFunc func(mt *mtest.T)
	}
	type Expected struct {
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given given no error saving must return correct entity": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := repository.NewOrderRepo(mt.DB)
					mt.AddMockResponses(mtest.CreateSuccessResponse())

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

					createdProduct, err := repo.Create(context.Background(), order)

					assert.Nil(t, err)
					assert.Equal(t, *createdProduct, order)

				},
			},
		},
		"given given error saving must return error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := repository.NewOrderRepo(mt.DB)
					mt.AddMockResponses(
						bson.D{
							{Key: "ok", Value: -1},
						},
					)

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

					createdProduct, err := repo.Create(context.Background(), order)

					assert.NotNil(t, err)
					assert.Nil(t, createdProduct)

				},
			},
		},
	}

	for _, tc := range tests {
		db := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		db.Run("", tc.given.mtestFunc)
	}
}

func TestUpdate(t *testing.T) {

	mpatch.PatchMethod(time.Now, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})

	type Given struct {
		mtestFunc func(mt *mtest.T)
	}
	type Expected struct {
	}
	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"given given no error updating must return no error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := repository.NewOrderRepo(mt.DB)
					mt.AddMockResponses(bson.D{
						{Key: "ok", Value: 1},
						{Key: "value", Value: bson.D{
							{Key: "_id", Value: "order_valid_id"},
							{Key: "customer_id", Value: "order_valid_customer_id"},
							{Key: "status", Value: canonical.ORDER_PREPARING},
							{Key: "created_at", Value: time.Now()},
							{Key: "updated_at", Value: time.Now()},
							{Key: "total", Value: 0.0},
							{Key: "order_items", Value: []bson.D{
								{
									{Key: "quantity", Value: 5},
									{Key: "product_id", Value: "product_id"},
									{Key: "name", Value: "product_name"},
									{Key: "price", Value: 10.0},
									{Key: "category", Value: "product_category"},
								},
								{
									{Key: "quantity", Value: 10},
									{Key: "product_id", Value: "product_id"},
									{Key: "name", Value: "product_name"},
									{Key: "price", Value: 10.0},
									{Key: "category", Value: "product_category"},
								},
							},
							},
						}},
					})

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

					err := repo.Update(context.Background(), "product_valid", order)

					assert.Nil(t, err)

				},
			},
		},
		"given error saving must return error": {
			given: Given{
				mtestFunc: func(mt *mtest.T) {
					repo := repository.NewOrderRepo(mt.DB)
					mt.AddMockResponses(
						bson.D{
							{Key: "ok", Value: -1},
						},
					)
					product := canonical.Order{
						ID:         "",
						CustomerID: "",
						Status:     0,
						CreatedAt:  time.Time{},
						UpdatedAt:  time.Time{},
						Total:      0,
						OrderItems: []canonical.OrderItem{},
					}

					err := repo.Update(context.Background(), "product_valid", product)

					assert.NotNil(t, err)

				},
			},
		},
	}

	for _, tc := range tests {
		db := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		db.Run("", tc.given.mtestFunc)
	}
}
