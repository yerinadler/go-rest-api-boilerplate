package services

import (
	"context"
	"errors"

	"github.com/example/go-rest-api-revision/internal/db/gorm/models"
	"github.com/example/go-rest-api-revision/internal/dtos"
	exception "github.com/example/go-rest-api-revision/pkg/exceptions"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type ProductService struct {
	tracer trace.Tracer
	db     *gorm.DB
}

func NewProductService(tracer trace.Tracer, db *gorm.DB) *ProductService {
	return &ProductService{
		tracer,
		db,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, dto *dtos.ProductDto) error {
	ctx, span := s.tracer.Start(ctx, "creating a new product")
	defer span.End()

	result := s.db.WithContext(ctx).Create(&models.Product{
		Name:        dto.Name,
		Description: dto.Description,
		UnitPrice:   dto.UnitPrice,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]dtos.ProductDto, error) {
	var products []models.Product
	result := s.db.WithContext(ctx).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(products) == 0 {
		return []dtos.ProductDto{}, nil
	}

	var productDtos []dtos.ProductDto

	for _, product := range products {
		productDtos = append(productDtos, dtos.ProductDto{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			UnitPrice:   product.UnitPrice,
		})
	}

	return productDtos, nil
}

func (s *ProductService) GetProductById(ctx context.Context, id string) (*dtos.ProductDto, error) {
	var product models.Product
	result := s.db.WithContext(ctx).First(&product, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &exception.NotFoundException
		}
		return nil, result.Error
	}

	return &dtos.ProductDto{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		UnitPrice:   product.UnitPrice,
	}, nil
}
