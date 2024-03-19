package rest

import (
	"context"
	"fmt"
	"net/http"
	"tech-challenge-order/internal/auth/token"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/service"

	"github.com/labstack/echo/v4"
)

type order struct {
	service service.OrderService
}

func NewOrderChannel() Order {
	return &order{
		service: service.NewOrderService(),
	}
}

func (p *order) RegisterGroup(g *echo.Group) {
	g.GET("/", p.Get)
	g.POST("/", p.Create)
	g.PUT("/:id", p.Update)
	g.PATCH("/", p.UpdateStatus)
	g.POST("/checkout", p.CheckoutOrder)
}

func (p *order) Get(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	status := ctx.QueryParam("status")

	response, err := p.get(ctx.Request().Context(), id, status)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{
			Message: "searching orders: " + err.Error(),
		})
	}

	if len(response) == 0 {
		return ctx.NoContent(http.StatusNotFound)
	} else if len(response) == 1 {
		return ctx.JSON(http.StatusOK, response[0])
	} else {
		return ctx.JSON(http.StatusOK, response)
	}
}

func (p *order) get(ctx context.Context, orderID string, status string) ([]OrderResponse, error) {
	if orderID != "" {
		order, err := p.service.GetByID(ctx, orderID)
		if err != nil {
			return nil, err
		}

		if order == nil {
			return nil, nil
		}

		return []OrderResponse{orderToResponse(*order)}, nil
	}

	var response []OrderResponse
	if status != "" {

		status, ok := canonical.MapOrderStatus[status]
		if !ok {
			return nil, fmt.Errorf("invalid status")
		}

		orders, err := p.service.GetByStatus(ctx, status)
		if err != nil {
			return nil, err
		}

		for _, order := range orders {
			response = append(response, orderToResponse(order))
		}

		return response, nil
	}

	orders, err := p.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		response = append(response, orderToResponse(order))
	}

	return response, nil
}

func (p *order) Create(c echo.Context) error {
	var orderRequest OrderRequest

	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid data").Error(),
		})
	}

	if orderRequest.OrderItems == nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid data").Error(),
		})
	}

	customerId, err := token.ExtractCustomerId(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid customer").Error(),
		})
	}

	orderCan := orderRequest.toCanonical(customerId)

	err = p.service.Create(context.Background(), *orderCan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

func (p *order) Update(c echo.Context) error {
	orderID := c.Param("id")
	if len(orderID) == 0 {
		return c.JSON(http.StatusBadRequest, Response{
			Message: "missing id query param",
		})
	}

	var orderRequest OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid data").Error(),
		})
	}

	if orderRequest.OrderItems == nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid data").Error(),
		})
	}

	customerId, err := token.ExtractCustomerId(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: fmt.Errorf("invalid customer").Error(),
		})
	}

	orderCan := orderRequest.toCanonical(customerId)

	err = p.service.Update(c.Request().Context(), orderID, *orderCan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

func (p *order) UpdateStatus(c echo.Context) error {
	orderID := c.QueryParam("id")
	paramStatus := c.QueryParam("status")

	if len(orderID) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "missing id query param"})
	}

	status, ok := canonical.MapOrderStatus[paramStatus]
	if !ok {
		return c.JSON(http.StatusBadRequest, Response{Message: "invalid status"})
	}

	err := p.service.UpdateStatus(c.Request().Context(), orderID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: "error updating order: " + err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

func (p *order) CheckoutOrder(c echo.Context) error {
	orderID := c.QueryParam("id")
	if len(orderID) == 0 {
		return c.JSON(http.StatusBadRequest, Response{Message: "missing id query param"})
	}

	order, err := p.service.CheckoutOrder(c.Request().Context(), orderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Message: "error checking out order: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, orderToResponse(*order))
}
