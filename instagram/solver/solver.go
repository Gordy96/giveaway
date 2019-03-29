package solver

import (
	"giveaway/instagram/commands"
	"sync"
)

type Solver struct {
	ch      chan commands.Command
	started bool
}

func (s *Solver) Enqueue(command commands.Command) chan interface{} {
	ch := make(chan interface{}, 1)
	command.SetChannel(ch)
	s.ch <- command
	return ch
}

func (s *Solver) Run() {
	mux := sync.Mutex{}
	mux.Lock()
	defer mux.Unlock()
	if s.started {
		return
	}
	go func() {
		for {
			command, ok := <-s.ch
			if !ok {
				break
			}
			command.Handle()
		}
	}()
	s.started = true
}

func (s *Solver) Close() {
	mux := sync.Mutex{}
	mux.Lock()
	defer mux.Unlock()
	close(s.ch)
	s.ch = make(chan commands.Command)
	s.started = false
}

var singleToneSolverInstance *Solver = nil

func New() *Solver {
	temp := &Solver{}
	temp.ch = make(chan commands.Command)
	return temp
}

func GetInstance() *Solver {
	mux := sync.Mutex{}
	mux.Lock()
	if singleToneSolverInstance == nil {
		singleToneSolverInstance = New()
	}
	mux.Unlock()
	return singleToneSolverInstance
}

func NewAndRun() *Solver {
	s := New()
	s.Run()
	return s
}

func GetRunningInstance() *Solver {
	s := GetInstance()
	s.Run()
	return s
}
