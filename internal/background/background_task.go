package background

import (
	"fmt"
	"time"
)

type TaskHandler struct {
	ticker *time.Ticker
	done   chan bool
	tasks  []func(t time.Time)
}

func NewHandler(interval time.Duration) TaskHandler {
	return TaskHandler{
		ticker: time.NewTicker(time.Second),
		done:   make(chan bool),
		tasks:  make([]func(t time.Time), 0),
	}
}

func (th *TaskHandler) AddTask(fn func(t time.Time)) {
	th.tasks = append(th.tasks, fn)
}

func (th *TaskHandler) Run() {
	go func() {
		defer th.ticker.Stop()
		for {
			select {
			case <-th.done:
				fmt.Println("Done")
				return
			case t := <-th.ticker.C:
				for _, v := range th.tasks {
					go v(t)
				}
			}
		}

	}()
}

func (th *TaskHandler) Stop() {
	th.done <- true
}
