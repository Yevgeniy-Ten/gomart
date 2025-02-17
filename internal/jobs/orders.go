package jobs

import (
	"context"
	"go.uber.org/zap"
	"gophermart/internal/domain"
	"sync"
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

func (j *OrdersJob) Run() {
	ctx := context.Background()

	orders, err := j.GetNewOrders(ctx, NumJobs)
	if err != nil {
		j.L.Error("failed to get new orders", zap.Error(err))
		return
	}
	doneCh := make(chan struct{})
	defer close(doneCh)
	statusChs := fanOut(doneCh, orders)
	statusesCh := fanIn(doneCh, statusChs...)
	result := make([]*domain.OrderWithAccrual, 0)
	for status := range statusesCh {
		result = append(result, status)
	}
	if err := j.UpdateOrdersWithAccrual(ctx, result); err != nil {
		j.L.Error("failed to update orders with accrual", zap.Error(err))
	}
}

func fanOut(doneCh chan struct{}, orders []string) []chan *domain.OrderWithAccrual {
	numWorkers := len(orders)
	channels := make([]chan *domain.OrderWithAccrual, numWorkers)

	for i, o := range orders {
		addResultCh := getStatus(doneCh, o)
		channels[i] = addResultCh
	}

	return channels
}

func getStatus(doneCh chan struct{}, order string) chan *domain.OrderWithAccrual {
	resultCh := make(chan *domain.OrderWithAccrual)

	go func() {
		defer close(resultCh)
		//request to accrual
		var accrualResponse domain.OrderWithAccrual
		select {
		case <-doneCh:
			return
		case resultCh <- &accrualResponse:
		}
	}()

	return resultCh
}

func fanIn(doneCh chan struct{}, resultChs ...chan *domain.OrderWithAccrual) chan *domain.OrderWithAccrual {
	finalCh := make(chan *domain.OrderWithAccrual)

	var wg sync.WaitGroup

	for _, ch := range resultChs {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			// получаем данные из канала
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
