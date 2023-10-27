package syncer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Chystik/gophermart/internal/infrastructure/webapi"
	"github.com/Chystik/gophermart/internal/models"
	"github.com/Chystik/gophermart/internal/usecase"
	"github.com/Chystik/gophermart/pkg/logger"
)

const (
	defaultRequestInterval time.Duration = 100 * time.Millisecond
)

var once sync.Once

type syncer struct {
	userRepo  usecase.UserRepository
	orderRepo usecase.OrderRepository
	accrual   usecase.AccrualWebAPI
	logger    logger.AppLogger
	i         time.Duration
	tick      chan struct{}
}

func NewSyncer(ur usecase.UserRepository, or usecase.OrderRepository, acc usecase.AccrualWebAPI, l logger.AppLogger, opts ...Options) *syncer {
	s := &syncer{
		userRepo:  ur,
		orderRepo: or,
		accrual:   acc,
		logger:    l,
		i:         defaultRequestInterval,
		tick:      make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *syncer) Run() error {
	s.tick <- struct{}{} // init tick

	for range s.tick {
		err := s.sync()
		if err != nil {
			s.logger.Error(err.Error())
			return err
		}
		time.Sleep(s.i)
		s.tick <- struct{}{}
	}

	return nil
}

func (s *syncer) Shutdown(ctx context.Context) error {
	once.Do(func() {
		close(s.tick)
	})
	return nil
}

func (s *syncer) sync() error {
	orders, err := s.orderRepo.GetUnprocessed(context.TODO())
	if err != nil {
		return err
	}

	for i := range orders {
		ctx := context.TODO()
		oldStatus := orders[i].Status

		// check order status in accrual service
		orders[i], err = s.accrual.GetOrder(ctx, orders[i])
		if err != nil {
			s.logger.Error(err.Error())
			var accrualErr *webapi.AccrualError

			// not accrusl service error
			if !errors.As(err, &accrualErr) {
				return err
			}

			// we make to many requests
			if accrualErr.ToManyRequests() {

				// change requests interval and restart sync task,
				// because we cant make requests for a while
				newRatePerMin, err := accrualErr.GetRateLimit()
				if err != nil {
					s.logger.Error(err.Error())
					return err
				}
				s.i = time.Duration((newRatePerMin / 60) * int(time.Second))
				return nil
			}
		}

		// update order with new status
		newStatus := orders[i].Status
		if newStatus != oldStatus {
			err = s.orderRepo.Update(ctx, orders[i])
			if err != nil {
				s.logger.Error(err.Error())
				return err
			}

			// if the new status is valid for accrual - change user balance
			if newStatus == "PROCESSED" {
				user, err := s.userRepo.Get(ctx, models.User{Login: orders[i].User})
				if err != nil {
					s.logger.Error(err.Error())
					return err
				}

				user.Balance += orders[i].Accrual
				err = s.userRepo.Update(ctx, user)
				if err != nil {
					s.logger.Error(err.Error())
					return err
				}
			}
		}

	}
	return nil
}
