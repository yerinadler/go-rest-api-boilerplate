package services

import (
	"context"
	"errors"

	"github.com/example/go-rest-api-boilerplate/internal/customer-service/db/gorm/models"
	"github.com/example/go-rest-api-boilerplate/internal/customer-service/dtos"
	exception "github.com/example/go-rest-api-boilerplate/pkg/exceptions"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type CustomerService struct {
	tracer trace.Tracer
	logger *logrus.Logger
	db     *gorm.DB
}

func NewCustomerService(tracer trace.Tracer, logger *logrus.Logger, db *gorm.DB) *CustomerService {
	return &CustomerService{
		tracer,
		logger,
		db,
	}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, dto *dtos.CustomerDto) error {
	ctx, span := s.tracer.Start(ctx, "creating a new product")
	defer span.End()

	customer := &models.Customer{
		Firstname:   dto.Firstname,
		Middlename:  dto.Middlename,
		Lastname:    dto.Lastname,
		DateOfBirth: dto.DateOfBirth,
	}

	result := s.db.WithContext(ctx).Create(customer)

	s.logger.WithContext(ctx).Infof("successfully created the customer with the ID of %d", customer.ID)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *CustomerService) GetCustomerById(ctx context.Context, id string) (*dtos.CustomerDto, error) {
	var customer models.Customer
	result := s.db.WithContext(ctx).First(&customer, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &exception.NotFoundException
		}
		return nil, result.Error
	}

	return &dtos.CustomerDto{
		Id:          customer.ID,
		Firstname:   customer.Firstname,
		Middlename:  customer.Middlename,
		Lastname:    customer.Lastname,
		DateOfBirth: customer.DateOfBirth,
	}, nil
}
