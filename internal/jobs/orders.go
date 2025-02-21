package jobs

import (
	"context"
	"errors"
	"gophermart/internal/domain"
	"sync"
	"time"

	"go.uber.org/zap"
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
			<-timer.C
			timer.Reset(interval)
			continue
		}
		if len(orders) == 0 {
			j.L.Info("no new orders")
			<-timer.C
			timer.Reset(interval)
			continue
		}
		doneCh := make(chan struct{})
		requestsToAccrual := j.fanOut(doneCh, orders)
		responsesFromAccrual := j.fanIn(doneCh, requestsToAccrual...)
		result := make([]*domain.OrderWithAccrual, 0)
		for order := range responsesFromAccrual {
			if order.Error != nil {
				var tooManyRequests *domain.TooManyRequestsError
				if errors.As(order.Error, &tooManyRequests) {
					interval = time.Duration(tooManyRequests.RetryAfter) * time.Second
				}
				continue
			}
			result = append(result, &domain.OrderWithAccrual{
				Number: order.Number,
				AccrualResponse: domain.AccrualResponse{
					Accrual: order.AccrualResponse.Accrual,
					Status:  order.AccrualResponse.Status,
				},
			})
		}
		j.L.Info("ORDERS", zap.Any("orders", result))
		if err := j.UpdateOrdersWithAccrual(ctx, result); err != nil {
			j.L.Error("failed to update orders with accrual", zap.Error(err))
		}
		close(doneCh)
		<-timer.C
		timer.Reset(interval)
	}
}
func (j *OrdersJob) getStatus(doneCh chan struct{}, order string, ctx context.Context, cancel context.CancelFunc) chan *domain.OrderInJobs {
	resultCh := make(chan *domain.OrderInJobs)

	go func() {
		select {
		case <-doneCh:
			cancel()
			return
		case <-ctx.Done():
			return
		default:
			fullOrder := domain.OrderInJobs{
				Number: order,
			}
			accrualResp, err := j.GetOrderStatus(order)
			if err != nil {
				j.L.Error("failed to get order status", zap.Error(err))
				fullOrder.Error = err
				cancel()
			} else {
				fullOrder.AccrualResponse = *accrualResp
				resultCh <- &fullOrder
			}
		}
		defer close(resultCh)
	}()
	return resultCh
}
func (j *OrdersJob) fanOut(doneCh chan struct{}, orders []string) []chan *domain.OrderInJobs {
	numWorkers := len(orders)
	channels := make([]chan *domain.OrderInJobs, numWorkers)
	ctx, cancel := context.WithCancel(context.Background())
	for i, o := range orders {
		response := j.getStatus(doneCh, o, ctx, cancel)
		channels[i] = response
	}
	return channels
}

func (j *OrdersJob) fanIn(doneCh chan struct{}, resultChs ...chan *domain.OrderInJobs) chan *domain.OrderInJobs {
	finalCh := make(chan *domain.OrderInJobs)

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
