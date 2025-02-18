package jobs

import (
	"context"
	"go.uber.org/zap"
	"gophermart/internal/domain"
	"sync"
	"time"
)

const NumJobs = 5

type Repo interface {
	GetNewOrders(ctx context.Context, limit int) ([]string, error)
	UpdateOrdersWithAccrual(ctx context.Context, orders []*domain.OrderWithAccrual) error
}

type OrdersJob struct {
	Repo
	*domain.Utils
}

func NewOrdersJob(repo Repo, utils *domain.Utils) *OrdersJob {
	return &OrdersJob{
		Repo:  repo,
		Utils: utils,
	}
}

func (j *OrdersJob) Run(initialInterval time.Duration) {
	interval := initialInterval
	timer := time.NewTimer(interval)
	defer timer.Stop()
	for {
		ctx := context.Background()
		orders, err := j.GetNewOrders(ctx, NumJobs)
		if err != nil {
			j.L.Error("failed to get new orders", zap.Error(err))
			return
		}
		doneCh := make(chan struct{})
		statusChs := j.fanOut(doneCh, orders)
		statusesCh := j.fanIn(doneCh, statusChs...)
		result := make([]*domain.OrderWithAccrual, 0)
		for status := range statusesCh {
			result = append(result, status)
		}
		if err := j.UpdateOrdersWithAccrual(ctx, result); err != nil {
			j.L.Error("failed to update orders with accrual", zap.Error(err))
		}
		close(doneCh)
		<-timer.C
		timer.Reset(interval)
	}
}
func (j *OrdersJob) getStatus(doneCh chan struct{}, order string) chan *domain.OrderWithAccrual {
	resultCh := make(chan *domain.OrderWithAccrual)

	go func() {
		defer close(resultCh)
		fullOrder := domain.OrderWithAccrual{
			Number: order,
		}
		accrualResp, err := j.GetOrderStatus(order)
		if err != nil {
			j.L.Error("failed to get order status", zap.Error(err))
			return
		}
		fullOrder.AccrualResponse = *accrualResp
		select {
		case <-doneCh:
			return
		case resultCh <- &fullOrder:
		}
	}()

	return resultCh
}
func (j *OrdersJob) fanOut(doneCh chan struct{}, orders []string) []chan *domain.OrderWithAccrual {
	numWorkers := len(orders)
	channels := make([]chan *domain.OrderWithAccrual, numWorkers)

	for i, o := range orders {
		addResultCh := j.getStatus(doneCh, o)
		channels[i] = addResultCh
	}

	return channels
}

func (j *OrdersJob) fanIn(doneCh chan struct{}, resultChs ...chan *domain.OrderWithAccrual) chan *domain.OrderWithAccrual {
	finalCh := make(chan *domain.OrderWithAccrual)

	var wg sync.WaitGroup

	for _, ch := range resultChs {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-doneCh:
					return
				case finalCh <- data:
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(finalCh)
	}()
	return finalCh
}
