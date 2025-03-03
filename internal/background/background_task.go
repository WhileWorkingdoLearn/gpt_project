package background

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type backgroundTask struct {
	id         string
	name       string
	started_At time.Time
	last_run   time.Time
	task       func()
}

type TaskTicket struct {
	id   string
	name string
}

type TaskInfo struct {
	id         string
	name       string
	started_At time.Time
	last_run   time.Time
}

func NewTask(name string, Task func()) backgroundTask {
	task := backgroundTask{
		id:         "",
		started_At: time.Now(),
		task:       Task,
		name:       name,
	}
	return task
}

func (bt backgroundTask) StartedAt() time.Time {
	return bt.started_At
}

func (bt backgroundTask) LastRun() time.Time {
	return bt.last_run
}

type TaskHandler struct {
	ticker   *time.Ticker
	interval time.Duration
	done     chan bool
	tasks    map[string]backgroundTask
	runUntil *time.Time
	isRuning bool
	cycles   atomic.Uint64
	mutex    sync.Mutex
}

func NewBackgroundWorker() *TaskHandler {
	return &TaskHandler{
		done:     make(chan bool),
		tasks:    make(map[string]backgroundTask, 0),
		isRuning: false,
		interval: 1 * time.Second,
		cycles:   atomic.Uint64{},
	}
}

func (th *TaskHandler) PastInterval() int {
	val := &th.cycles
	return int(val.Load())
}

func (th *TaskHandler) ResetCounter() {
	val := &th.cycles
	val.Store(0)
}

func (th *TaskHandler) AddTask(bt backgroundTask) TaskTicket {
	th.mutex.Lock()
	defer th.mutex.Unlock()
	th.tasks[bt.id] = bt
	return TaskTicket{id: bt.id, name: bt.name}
}

func (th *TaskHandler) RemoveTask(id string) {
	th.mutex.Lock()
	defer th.mutex.Unlock()
	delete(th.tasks, id)
}

func (th *TaskHandler) GetTaskCount() int {
	return len(th.tasks)
}

func (th *TaskHandler) GetTaskInfo(id string) (TaskInfo, error) {
	th.mutex.Lock()
	defer th.mutex.Unlock()
	t, found := th.tasks[id]
	if !found {
		return TaskInfo{}, fmt.Errorf("id was not found")
	}
	return TaskInfo{id: t.id, name: t.name, started_At: t.started_At, last_run: t.last_run}, nil
}

func (th *TaskHandler) SetInterval(value time.Duration) {
	th.interval = value
	if th.ticker != nil {
		th.ticker.Reset(th.interval)
	}
}

func (th *TaskHandler) GetSetInterval() time.Duration {
	return th.interval
}

func (th *TaskHandler) internalRun() {
	defer th.ticker.Stop()

	for {
		select {
		case t := <-th.ticker.C:
			if th.runUntil != nil {
				if t.After(*th.runUntil) {
					fmt.Println("Done")
					th.Stop()
					return
				}
			}

			fmt.Println("Stil running")
			th.cycles.Add(1)
			for _, v := range th.tasks {
				v.last_run = time.Now()
				go v.task()
			}
		case done := <-th.done:
			fmt.Println("Done th.Done ", done)
			th.ticker.Stop()
			th.isRuning = false
			return
		}
	}

}

func (th *TaskHandler) Run() error {

	if th.isRuning {
		return fmt.Errorf("task is already Running")
	}

	th.ticker = time.NewTicker(th.interval)

	go th.internalRun()
	return nil
}

func (th *TaskHandler) RunUntil(until time.Duration) error {
	if th.isRuning {
		return fmt.Errorf("task is already Running")
	}
	timeUntil := time.Now().Add(until)
	th.runUntil = &timeUntil
	th.ticker = time.NewTicker(th.interval)

	go th.internalRun()

	return nil
}

func (th *TaskHandler) Stop() {
	if th.isRuning {
		close(th.done)
	}

}
