package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Queue struct {
	Items        []interface{}
	MaxSize      int
	Done         chan bool
	Mu           sync.Mutex
	LastItemTime time.Time
	ClearedTimes int
}

type QueuePayload struct {
	Content string
	Type    string
}

func newQueue(maxSize int) *Queue {
	return &Queue{
		Items:   make([]interface{}, 0, maxSize),
		MaxSize: maxSize,
		Done:    make(chan bool, 1),
	}
}

func (q *Queue) add(item string) error {
	q.Mu.Lock()
	defer q.Mu.Unlock()

	if len(q.Items) >= q.MaxSize {
		return fmt.Errorf("queue is full")
	}

	q.Items = append(q.Items, item)
	q.LastItemTime = time.Now()

	//fmt.Printf("Queue type %s has %d items\n", item, len(q.Items))

	if len(q.Items) >= q.MaxSize {
		q.Done <- true
	}
	return nil
}

func purgeFullQueues(mux *sync.Mutex, queues map[string]*Queue) {
	for {
		mux.Lock()
		for _, q := range queues {
			select {
			case <-q.Done:
				//fmt.Printf("Purging queue %v \n", t)
				q.Items = nil
				q.ClearedTimes++
			default:
				// do nothing
			}
		}
		mux.Unlock()
	}
}

func purgeTimedOutQueues(mux *sync.Mutex, queues map[string]*Queue) {
	for {
		time.Sleep(1 * time.Second)

		mux.Lock()
		for _, q := range queues {
			if !q.LastItemTime.IsZero() && time.Since(q.LastItemTime) > (30*time.Second) {
				//fmt.Printf("Purging old queue %v\n", t)
				q.Items = nil
				q.LastItemTime = time.Time{}
				q.ClearedTimes++
			}
		}
		mux.Unlock()
	}
}

func printCurrentStatus(queues map[string]*Queue) {
	for {
		time.Sleep(1 * time.Second)

		fmt.Print("\r")
		for queueName, queue := range queues {
			fmt.Printf("Queue %s: %d items cleared %d times | ", queueName, len(queue.Items), queue.ClearedTimes)
		}
		fmt.Printf("Total Queues: %d", len(queues))
	}
}

func main() {
	queues := make(map[string]*Queue)
	var mux sync.Mutex

	go purgeFullQueues(&mux, queues)
	go purgeTimedOutQueues(&mux, queues)
	go printCurrentStatus(queues)

	http.HandleFunc("/queue", func(w http.ResponseWriter, r *http.Request) {
		var payload QueuePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mux.Lock()

		q, ok := queues[payload.Type]
		if !ok {
			q = newQueue(10)
		}

		q.add(payload.Content)

		queues[payload.Type] = q

		w.WriteHeader(http.StatusCreated)

		mux.Unlock()
	})

	http.ListenAndServe(":8080", nil)
}
