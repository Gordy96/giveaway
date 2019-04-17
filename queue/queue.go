package queue

import "sync"

type Command interface {
	Handle()
}

func NewWorker(pool chan chan Command) *Worker {
	return &Worker{
		pool,
		make(chan Command),
		make(chan bool),
	}
}

type Worker struct {
	pool  chan chan Command
	input chan Command
	quit  chan bool
}

func (w *Worker) Start() {
	go func() {
		for {
			w.pool <- w.input
			select {
			case job := <-w.input:
				job.Handle()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Queue struct {
	workers     []*Worker
	pool        chan chan Command
	WorkerCount int
	jobQueue    chan Command
}

func (q *Queue) Run() {
	for i := 0; i < q.WorkerCount; i++ {
		q.workers[i] = NewWorker(q.pool)
		q.workers[i].Start()
	}
	go func() {
		for {
			select {
			case job := <-q.jobQueue:
				go func(c Command) {
					workerChan := <-q.pool
					workerChan <- c
				}(job)
			}
		}
	}()
}

func (q *Queue) Enqueue(c Command) {
	q.jobQueue <- c
}

func NewQueue(numWorkers int) *Queue {
	q := &Queue{
		make([]*Worker, numWorkers),
		make(chan chan Command),
		numWorkers,
		make(chan Command),
	}
	return q
}

var singletonMux = sync.Mutex{}

var globalInstance *Queue = nil

func GetGlobalInstance(numWorkers ...int) *Queue {
	var i int = 10
	if len(numWorkers) > 0 {
		i = numWorkers[0]
	}
	singletonMux.Lock()
	if globalInstance == nil {
		globalInstance = NewQueue(i)
	} else if len(numWorkers) > 0 && globalInstance.WorkerCount != numWorkers[0] {
		for _, w := range globalInstance.workers {
			w.Stop()
		}
		globalInstance = NewQueue(i)
	}
	singletonMux.Unlock()
	return globalInstance
}
