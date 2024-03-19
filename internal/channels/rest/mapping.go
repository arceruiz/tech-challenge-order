package rest

import (
	"tech-challenge-order/internal/canonical"
)

func productToResponse(p canonical.Product) ProductItem {
	return ProductItem{
		ID:       p.ID,
		Name:     p.Name,
		Price:    p.Price,
		Category: p.Category,
	}
}

func (o *OrderRequest) toCanonical(customerId string) *canonical.Order {
	return &canonical.Order{
		OrderItems: toCanonical(o.OrderItems),
		CustomerID: customerId,
	}
}

func orderToResponse(order canonical.Order) OrderResponse {
	var productsList []OrderItemResponse

	for _, item := range order.OrderItems {
		oi := OrderItemResponse{}
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
		Total:      order.Total,
	}
}

func toCanonical(ordersItems []OrderItem) map[string]*canonical.OrderItem {
	v := map[string]*canonical.OrderItem{}

	for _, item := range ordersItems {
		v[item.ProductId] = &canonical.OrderItem{
			Quantity: item.Quantity,
		}
	}

	return v
}

func keyByValue(myMap map[string]canonical.OrderStatus, value canonical.OrderStatus) string {
	for k, v := range myMap {
		if value == v {
			return k
		}
	}
	return ""
}
