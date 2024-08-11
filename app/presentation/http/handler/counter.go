package handler

import (
	"log"
	"net/http"

	"github.com/sameer-bishnoi/counter/app/domain/entities/dto"
	"github.com/sameer-bishnoi/counter/app/infra/shttp"
)

type RequestCounterServicer interface {
	GetRequestCount() (*dto.RequestCounterDTO, error)
}

type RequestCounter struct {
	requestCounterService RequestCounterServicer
}

func NewRequestCounter() *RequestCounter {
	return &RequestCounter{}
}

func (c *RequestCounter) WithRequestCounterService(rcs RequestCounterServicer) {
	c.requestCounterService = rcs
}

func (c *RequestCounter) GetCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("unsupported http method"))
	}

	data, err := c.requestCounterService.GetRequestCount()
	if err != nil {
		log.Printf("error while getting the counters: %v", err)
		shttp.FailedResponseJSON(w, http.StatusInternalServerError, "1001", "internal error")
		return
	}
	log.Printf("Response: %+v", data)
	shttp.ResponseJSON(w, http.StatusOK, data)
}
