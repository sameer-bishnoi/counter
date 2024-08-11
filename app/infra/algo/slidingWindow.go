package algo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type SlidingWindower interface {
	Enqueue(timestamp int64)
	GetSize() int64
	RemoveExpired(beforeTimestamp int64)
	IsEmpty() bool
	Load(file *os.File) error
	Store(file *os.File) error
}

type Queue struct {
	mutex      sync.Mutex
	timestamps []int64
}

func NewQueue() SlidingWindower {
	return &Queue{
		timestamps: make([]int64, 0),
	}
}

func (q *Queue) Enqueue(timestamp int64) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.timestamps = append(q.timestamps, timestamp)
}

func (q *Queue) GetSize() int64 {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return int64(len(q.timestamps))
}

func (q *Queue) RemoveExpired(beforeTimestamp int64) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	newStartIndex := -1
	for idx := range q.timestamps {
		if q.timestamps[idx] < beforeTimestamp {
			newStartIndex = idx
			continue
		}
		break
	}
	if newStartIndex < len(q.timestamps) {
		q.timestamps = q.timestamps[newStartIndex+1:]
	}
}

func (q *Queue) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.timestamps) == 0
}

func (q *Queue) Load(file *os.File) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		timestamp, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		q.timestamps = append(q.timestamps, timestamp)
	}
	return scanner.Err()
}

func (q *Queue) Store(file *os.File) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	writer := bufio.NewWriter(file)
	for _, timestamp := range q.timestamps {
		_, err := writer.WriteString(strconv.FormatInt(timestamp, 10) + "\n")
		if err != nil {
			return fmt.Errorf("unable to write: %v", err)
		}
	}
	err := writer.Flush()
	if err != nil {
		return fmt.Errorf("unable to flush the writer: %v", err)
	}
	return nil
}
