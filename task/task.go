package task

import (
	"time"
	"reflect"
	"crypto/sha1"
	"io"
	"fmt"
)

type ID string
type Function interface{}
type Param interface{}

type Trigger struct {
	IsRecurring bool
	LastRun     time.Time
	NextRun     time.Time
	Duration    time.Duration
}
type FuncMeta struct {
	Name   string
	Func   Function
	Params []Param
}

type Task struct {
	Tg   Trigger
	Func FuncMeta
	Id   ID
}

func NewTask(function FuncMeta, trigger Trigger) *Task {
	task := &Task{
		Func: function,
		Tg:   trigger,
	}
	task.Id = task.hash()
	return task
}

func (task *Task) IsDue() bool {
	timeNow := time.Now()
	return timeNow == task.Tg.NextRun || timeNow.After(task.Tg.NextRun)
}

func (task *Task) Run() {
	task.scheduleNextRun()
	function := reflect.ValueOf(task.Func.Func)
	params := make([]reflect.Value, len(task.Func.Params))
	for i, param := range task.Func.Params {
		params[i] = reflect.ValueOf(param)
	}
	function.Call(params)
}

func (task *Task) hash() ID {
	hash := sha1.New()
	_, _ = io.WriteString(hash, task.Func.Name)
	_, _ = io.WriteString(hash, fmt.Sprintf("%+v", task.Func.Params))
	_, _ = io.WriteString(hash, fmt.Sprintf("%s", task.Tg.Duration))
	_, _ = io.WriteString(hash, fmt.Sprintf("%t", task.Tg.IsRecurring))
	return ID(fmt.Sprintf("%x", hash.Sum(nil)))
}

func (task *Task) scheduleNextRun() {
	if !task.Tg.IsRecurring {
		return
	}
	task.Tg.LastRun = task.Tg.NextRun
	task.Tg.NextRun = task.Tg.NextRun.Add(task.Tg.Duration)
}
