package store

import (
	"github.com/ISchengchen/go_scheduler/task"
)

type MemoryStorage struct {
	tasks []task.Task
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (memStore *MemoryStorage) Add(task task.Task) error {
	memStore.tasks = append(memStore.tasks, task)
	return nil
}

func (memStore *MemoryStorage) GetDueJobs() ([]task.Task) {
	var tasks []task.Task
	for _, t := range memStore.tasks {
		if t.IsDue() {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

//TODO improve remove logic
func (memStore *MemoryStorage) Remove(t task.Task) error {
	var newTasks []task.Task
	for _, existingTask := range memStore.tasks {
		if existingTask.Id == t.Id {
			continue
		}
		newTasks = append(newTasks, existingTask)
	}
	memStore.tasks = newTasks
	return nil
}
