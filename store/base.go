package store

import (
	"github.com/ISchengchen/go_scheduler/task"
)

type TaskStore interface {
	Add(task task.Task) error
	GetDueJobs() ([]task.Task)
	Remove(task task.Task) error
}
