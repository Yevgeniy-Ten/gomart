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

	orders, err := j.GetNewOrders(context.Background(), NumJobs)
	if err != nil {
		j.L.Error("failed to get new orders", zap.Error(err))
		return
	}
	doneCh := make(chan struct{})
	defer close(doneCh)
	statusChs := fanOut(doneCh, orders)
	statusesCh := fanIn(doneCh, statusChs...)
	result := make([]string, 0)
	for status := range statusesCh {
		result = append(result, status)
	}
}

func fanOut(doneCh chan struct{}, orders []string) []chan string {
	numWorkers := len(orders)
	channels := make([]chan string, numWorkers)

	for i, o := range orders {
		addResultCh := getStatus(doneCh, o)
		channels[i] = addResultCh
	}

	return channels
}

func getStatus(
	doneCh chan struct{},
	order string) chan string {
	// канал, в который отправляются результаты
	resultCh := make(chan string)

	go func() {
		defer close(resultCh)
		//request to accrual
		select {
		case <-doneCh:
			return
		case resultCh <- order:
		}
	}()

	return resultCh
}

func fanIn(doneCh chan struct{}, resultChs ...chan string) chan string {
	finalCh := make(chan string)

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
