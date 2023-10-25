package webapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/pkg/httpclient"
)

var (
	errBadStatusCode = errors.New("bad status code")
)

type AccrualError struct {
	StatusCode int
	Err        error
}

func (r *AccrualError) Error() string {
	return r.Err.Error()
}

func (r *AccrualError) ToManyRequests() bool {
	return r.StatusCode == http.StatusTooManyRequests // 503
}

// GetRateLimit returns allowed requests per minute
func (r *AccrualError) GetRateLimit() (int, error) {
	var n []byte
	str := r.Err.Error()

	for i := range str {
		if byte(47) < str[i] && str[i] < byte(59) {
			n = append(n, str[i])
		}
	}

	return strconv.Atoi(string(n))
}

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
		switch resp.StatusCode {
		case http.StatusNoContent:
			return order, nil
		case http.StatusTooManyRequests:
			rl, err := io.ReadAll(resp.Body)
			if err != nil {
				return order, err
			}
			return order, &AccrualError{
				StatusCode: resp.StatusCode,
				Err:        errors.New(string(rl)),
			}
		default:
			return order, &AccrualError{
				StatusCode: resp.StatusCode,
				Err:        errBadStatusCode,
			}
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&rBody)
	if err != nil {
		return order, err
	}

	order.Status = rBody.Status
	order.Accrual = rBody.Accrual

	return order, nil
}
