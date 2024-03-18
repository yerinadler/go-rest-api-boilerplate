package services

import (
	"context"

	"github.com/example/go-rest-api-revision/internal/db/gorm/models"
	"github.com/example/go-rest-api-revision/internal/dtos"
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
			Name:        product.Name,
			Description: product.Description,
			UnitPrice:   product.UnitPrice,
		})
	}

	return productDtos, nil
}
