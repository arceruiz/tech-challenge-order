package rest

import (
	"tech-challenge-order/internal/canonical"
)

func (p *ProductItem) toCanonical() canonical.Product {
	return canonical.Product{
		ID:       p.ID,
		Name:     p.Name,
		Price:    p.Price,
		Category: p.Category,
	}
}

func productToResponse(p canonical.Product) ProductItem {
	return ProductItem{
		ID:       p.ID,
		Name:     p.Name,
		Price:    p.Price,
		Category: p.Category,
	}
}

func (o *OrderRequest) toCanonical() *canonical.Order {
	var orderItems []canonical.OrderItem

	for _, item := range o.OrderItems {
		orderItems = append(orderItems, item.toCanonical())
	}

	x := canonical.MapOrderStatus["RECEIVED"]
	if o.Status != "" {
		ok := false
		x, ok = canonical.MapOrderStatus[o.Status]
		if !ok {
			return nil
		}
	}
	return &canonical.Order{
		ID:         o.ID,
		CustomerID: o.CustomerID,
		Status:     x,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
		Total:      o.Total,
		OrderItems: orderItems,
	}
}

func orderToResponse(order canonical.Order) OrderResponse {
	var productsList []OrderItem

	for _, item := range order.OrderItems {
		oi := OrderItem{}
		oi.ProductItem = productToResponse(item.Product)
		oi.Quantity = item.Quantity
		productsList = append(productsList, oi)
	}

	return OrderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     keyByValue(canonical.MapOrderStatus, order.Status),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		Products:   productsList,
	}
}

func (items *OrderItem) toCanonical() canonical.OrderItem {
	oi := canonical.OrderItem{}
	oi.Product = items.ProductItem.toCanonical()
	oi.Quantity = items.Quantity
	return oi
}

func keyByValue(myMap map[string]canonical.OrderStatus, value canonical.OrderStatus) string {
	for k, v := range myMap {
		if value == v {
			return k
		}
	}
	return ""
}
