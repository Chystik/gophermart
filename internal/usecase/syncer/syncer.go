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

	"github.com/avito-tech/go-transaction-manager/trm"
)

const (
	defaultRequestInterval time.Duration = 100 * time.Millisecond
)

type syncer struct {
	userRepo  usecase.UserRepository
	orderRepo usecase.OrderRepository
	accrual   usecase.AccrualWebAPI
	logger    logger.AppLogger
	i         time.Duration
	tick      time.Ticker
	quit      chan struct{}
	once      sync.Once
	trm       trm.Manager
}

func NewSyncer(
	ur usecase.UserRepository,
	or usecase.OrderRepository,
	trm trm.Manager,
	acc usecase.AccrualWebAPI,
	l logger.AppLogger,
	opts ...Options) *syncer {

	s := &syncer{
		userRepo:  ur,
		orderRepo: or,
		trm:       trm,
		accrual:   acc,
		logger:    l,
		i:         defaultRequestInterval,
		quit:      make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.tick = *time.NewTicker(s.i)

	return s
}

func (s *syncer) Run() error {

	for {
		select {
		case <-s.quit:
			s.tick.Stop()
			return nil
		case <-s.tick.C:
			webAPIReqCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			err := s.sync(webAPIReqCtx)
			if err != nil {
				s.logger.Error(err.Error())
				return err
			}
		}
	}
}

func (s *syncer) Shutdown(ctx context.Context) error {
	s.once.Do(func() {
		s.quit <- struct{}{}
	})
	return nil
}

func (s *syncer) sync(ctx context.Context) error {
	orders, err := s.orderRepo.GetUnprocessed(ctx)
	if err != nil {
		return err
	}

	for i := range orders {
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
			if newStatus == models.Processed {
				err = s.trm.Do(ctx, func(ctx context.Context) error {
					user, err := s.userRepo.Get(ctx, models.User{Login: orders[i].User})
					if err != nil {
						s.logger.Error(err.Error())
						return err
					}

					user.Balance.Add(orders[i].Accrual)
					err = s.userRepo.Update(ctx, user)
					if err != nil {
						s.logger.Error(err.Error())
						return err
					}

					return s.orderRepo.Update(ctx, orders[i])
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
