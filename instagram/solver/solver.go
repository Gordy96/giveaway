package solver

import (
	"giveaway/instagram/commands"
	"sync"
)

type Solver struct {
	ch      chan commands.Command
	started bool
	mux     sync.Mutex
}

func (s *Solver) Enqueue(command commands.Command) chan interface{} {
	ch := make(chan interface{}, 1)
	command.SetChannel(ch)
	s.ch <- command
	return ch
}

func (s *Solver) Run() {
	s.mux.Lock()
	defer s.mux.Unlock()
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
	s.mux.Lock()
	defer s.mux.Unlock()
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

var solverSingletonMux = sync.Mutex{}

func GetInstance() *Solver {
	solverSingletonMux.Lock()
	if singleToneSolverInstance == nil {
		singleToneSolverInstance = New()
	}
	solverSingletonMux.Unlock()
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
