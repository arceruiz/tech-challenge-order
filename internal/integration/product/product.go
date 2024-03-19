package product

import (
	"context"
	"strconv"
	"tech-challenge-order/internal/canonical"
	"tech-challenge-order/internal/config"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductService interface {
	GetProducts(ctx context.Context, orderItems map[string]*canonical.OrderItem) error
}

type productService struct {
	productService ProductServiceClient
}

func NewProduct() ProductService {
	client, err := grpc.Dial(":"+config.Get().Server.ProductPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatal(err)
	}

	grpcClient := NewProductServiceClient(client)

	return &productService{
		productService: grpcClient,
	}
}

func (p *productService) GetProducts(ctx context.Context, orderItems map[string]*canonical.OrderItem) error {
	var idList []string

	for id := range orderItems {
		idList = append(idList, id)
	}

	products, err := p.productService.GetProduct(ctx, &Ids{
		Ids: idList,
	})
	if err != nil {
		return err
	}

	for _, product := range products.Products {
		price, err := strconv.ParseFloat(product.Price, 64)
		if err != nil {
			return err
		}

		p := orderItems[product.Id]

		p.ID = product.Id
		p.Name = product.Name
		p.Price = price
		p.Category = product.Category
	}

	return nil
}
