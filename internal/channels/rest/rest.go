package rest

import (
	"tech-challenge-order/internal/config"
	"tech-challenge-order/internal/middlewares"

	"github.com/labstack/echo/v4"
)

type Order interface {
	RegisterGroup(g *echo.Group)
	Create(c echo.Context) error
	Get(c echo.Context) error
	Update(c echo.Context) error
	UpdateStatus(c echo.Context) error
	CheckoutOrder(c echo.Context) error
}
type rest struct {
	order Order
}

func New(channel Order) rest {
	return rest{
		order: channel,
	}
}

func (r rest) Start() error {
	router := echo.New()

	router.Use(middlewares.Logger)

	mainGroup := router.Group("/api")

	orderGroup := mainGroup.Group("/order")
	r.order.RegisterGroup(orderGroup)
	orderGroup.Use(middlewares.Authorization)

	return router.Start(":" + config.Get().Server.Port)
}
