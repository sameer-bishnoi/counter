package usecase

import (
	"time"

	"github.com/sameer-bishnoi/counter/app/domain/entities/dto"
	"github.com/sameer-bishnoi/counter/app/infra/algo"
)

type RequestCounterService struct {
	queue algo.SlidingWindower
}

func NewRequestCounterService() *RequestCounterService {
	return &RequestCounterService{}
}

func (rcs *RequestCounterService) WithQueue(queue algo.SlidingWindower) {
	rcs.queue = queue
}

func (rcs *RequestCounterService) GetRequestCount() (*dto.RequestCounterDTO, error) {
	requestCount := rcs.queue.GetSize()
	rcs.queue.Enqueue(time.Now().UTC().Unix())

	return &dto.RequestCounterDTO{
		Count: requestCount,
	}, nil
}
