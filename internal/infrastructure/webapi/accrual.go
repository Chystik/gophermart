package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/pkg/httpclient"
)

var (
	errBadStatusCode = "resp status code: %s"
)

type respBody struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type accrual struct {
	client httpclient.HTTPClient
	url    string
}

func NewAccrualWebAPI(c httpclient.HTTPClient, opts ...Options) *accrual {
	a := &accrual{
		client: c,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *accrual) GetOrder(ctx context.Context, order models.Order) (models.Order, error) {
	var rBody respBody
	url := fmt.Sprintf("%s/api/orders/%s", a.url, order.Number)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return order, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return order, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			order.Status = "NEW"
			order.UploadedAt = models.RFC3339Time{Time: time.Now()}
			return order, nil
		}
		return order, fmt.Errorf(errBadStatusCode, resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&rBody)
	if err != nil {
		return order, err
	}

	order.Status = rBody.Status
	order.Accrual = rBody.Accrual

	return order, nil
}
