package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/sameer-bishnoi/counter/app/infra/algo"
)

type Worker interface {
	BackgroundJob(ctx context.Context)
}

type Job struct {
	queue    algo.SlidingWindower
	interval time.Duration
}

func NewJob() *Job {
	return &Job{}
}

func (j *Job) WithQueue(q algo.SlidingWindower) {
	j.queue = q
}

func (j *Job) WithInterval(interval time.Duration) {
	j.interval = interval
}

func (j *Job) BackgroundJob(done <-chan struct{}) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			log.Println("shutting down the background job")
			return
		case <-ticker.C:
			beforeTimestamp := time.Now().UTC().Add(-1 * time.Minute).Unix()
			j.queue.RemoveExpired(beforeTimestamp)
		}
	}
}
