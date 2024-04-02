package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/example/go-rest-api-boilerplate/pkg/responses"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Customer struct {
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	DateOfBirth time.Time `json:"dateOfBirth"`
}

type CustomerAdapter struct {
	url string
}

func NewCustomerAdapter(url string) *CustomerAdapter {
	return &CustomerAdapter{url}
}

func (a *CustomerAdapter) GetCustomerById(ctx context.Context, id int) (*Customer, error) {
	resp, err := otelhttp.Get(ctx, fmt.Sprintf("%s/%s/%d", a.url, "customers", 2))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp responses.ApiResponse[Customer]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &apiResp.Data, nil
}
