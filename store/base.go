package store

import (
	"go_scheduler/task"
)

type TaskStore interface {
	Add(task task.Task) error
	GetDueJobs() ([]task.Task)
	Remove(task task.Task) error
}
