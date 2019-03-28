package solver

import (
	"giveaway/instagram/commands"
)

type Solver struct {
	ch 			chan commands.Command
	started 	bool
}

func (s *Solver) Enqueue(command commands.Command) chan interface{} {
	ch := make(chan interface{})
	command.SetChannel(ch)
	s.ch <- command
	return ch
}

func (s *Solver) Run() {
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
	close(s.ch)
	s.ch = make(chan commands.Command)
	s.started = false
}

var singleToneSolverInstance *Solver = nil

func New() *Solver {
	if singleToneSolverInstance == nil {
		singleToneSolverInstance = &Solver{}
		singleToneSolverInstance.ch = make(chan commands.Command)
	}
	return singleToneSolverInstance
}

func NewAndRun() *Solver {
	s := New()
	s.Run()
	return s
}
