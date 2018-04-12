package go_scheduler

import (
	"time"
	"reflect"
	"fmt"
	"runtime"
	"os"
	"os/signal"
	"syscall"
	"github.com/ISchengchen/go_scheduler/store"
	"github.com/ISchengchen/go_scheduler/task"
)

type Scheduler struct {
	stopChan  chan bool
	taskStore store.TaskStore
}

func (s *Scheduler) RunAt(time time.Time, function task.Function, params ...task.Param) (task.ID, error) {
	funcMeta, err := ResoleveFunc(function, params)
	if err != nil {
		return "", nil
	}
	tg := task.Trigger{
		NextRun: time,
	}
	t := task.NewTask(funcMeta, tg)
	return t.Id, nil
}

func ResoleveFunc(function task.Function, params []task.Param) (task.FuncMeta, error) {
	funcValue := reflect.ValueOf(function)
	if funcValue.Kind() != reflect.Func {
		return task.FuncMeta{}, fmt.Errorf("参数必须为一个参数")
	}
	name := runtime.FuncForPC(funcValue.Pointer()).Name()
	funcMeta := task.FuncMeta{
		Name:   name,
		Func:   function,
		Params: params,
	}
	return funcMeta, nil
}

func NewScheduler(store store.TaskStore) Scheduler {
	return Scheduler{
		stopChan:  make(chan bool),
		taskStore: store,
	}
}

func (s *Scheduler) RunAfter(duration time.Duration, function task.Function, params ...task.Param) (task.ID, error) {
	return s.RunAt(time.Now().Add(duration), function, params...)
}

func (s *Scheduler) RunEvery(duration time.Duration, function task.Function, params ...task.Param) (task.ID, error) {
	funcMeta, err := ResoleveFunc(function, params)
	if err != nil {
		return "", nil
	}
	tg := task.Trigger{
		IsRecurring: true,
		NextRun:     time.Now().Add(duration),
		Duration:    duration,
	}
	t := task.NewTask(funcMeta, tg)
	s.taskStore.Add(*t)
	return t.Id, nil
}

func (s *Scheduler) Start() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				s.processTasks()
			case <-sigChan:
				s.stopChan <- true
			case <-s.stopChan:
				close(s.stopChan)
			}
		}
	}()
	return nil
}

func (s *Scheduler) processTasks() {
	tasks := s.taskStore.GetDueJobs()
	for _, t := range tasks {
		t.Run()
	}
}

func (s *Scheduler) Wait() {
	<-s.stopChan
}
