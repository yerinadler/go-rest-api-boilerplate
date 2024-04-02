package services

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/example/go-rest-api-boilerplate/internal/product-service/adapters"
	"github.com/example/go-rest-api-boilerplate/internal/product-service/db/gorm/models"
	"github.com/example/go-rest-api-boilerplate/internal/product-service/dtos"
	"github.com/example/go-rest-api-boilerplate/internal/product-service/events"
	exception "github.com/example/go-rest-api-boilerplate/pkg/exceptions"
	"github.com/example/go-rest-api-boilerplate/pkg/messaging/kafka"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type ProductService struct {
	tracer          trace.Tracer
	logger          *logrus.Logger
	db              *gorm.DB
	producer        *kafka.KafkaProducer
	customerAdapter *adapters.CustomerAdapter
}

func NewProductService(tracer trace.Tracer, logger *logrus.Logger, db *gorm.DB, producer *kafka.KafkaProducer, customerAdapter *adapters.CustomerAdapter) *ProductService {
	return &ProductService{
		tracer,
		logger,
		db,
		producer,
		customerAdapter,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, dto *dtos.ProductDto) error {
	ctx, span := s.tracer.Start(ctx, "creating a new product")
	defer span.End()

	data, _ := s.customerAdapter.GetCustomerById(ctx, 2)

	s.logger.WithContext(ctx).Infof("successfullly retrieved the customer data %v", data)

	product := &models.Product{
		Name:        dto.Name,
		Description: dto.Description,
		UnitPrice:   dto.UnitPrice,
	}
	result := s.db.WithContext(ctx).Create(product)

	s.logger.WithContext(ctx).Infof("successfully created the product %d : %s", product.ID, product.Name)

	event := &events.ProductCreated{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		UnitPrice:   product.UnitPrice,
	}

	jsonEvent, _ := json.Marshal(event)

	if err := s.producer.Publish(ctx, "test", string(jsonEvent), strconv.FormatUint(uint64(product.ID), 10)); err != nil {
		return err
	}

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
