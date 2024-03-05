package repository

import (
	"context"
	"tech-challenge-order/internal/canonical"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collection = "order"
)

type OrderRepository interface {
	GetAll(context.Context) ([]canonical.Order, error)
	Create(context.Context, canonical.Order) (*canonical.Order, error)
	Update(context.Context, string, canonical.Order) error
	GetByID(context.Context, string) (*canonical.Order, error)
	GetByStatus(context.Context, int) ([]canonical.Order, error)
}

type orderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepo() OrderRepository {
	return &orderRepository{collection: NewMongo().Collection(collection)}
}

func (r *orderRepository) GetAll(ctx context.Context) ([]canonical.Order, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	var results []canonical.Order
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *orderRepository) Create(ctx context.Context, order canonical.Order) (*canonical.Order, error) {
	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Update(ctx context.Context, id string, updatedOrder canonical.Order) error {
	filter := bson.M{"_id": id}
	fields := bson.M{"$set": updatedOrder}

	_, err := r.collection.UpdateOne(ctx, filter, fields)
	if err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*canonical.Order, error) {
	var order canonical.Order

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetByStatus(ctx context.Context, status int) ([]canonical.Order, error) {
	filter := bson.D{{Key: "status", Value: status}}
	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []canonical.Order
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
