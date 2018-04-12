package main

import (
	"log"
	"github.com/ISchengchen/go_scheduler/store"
	"github.com/ISchengchen/go_scheduler"
	"time"
)

func TaskWithoutArgs() {
	log.Println("task without args is excuted")
}

func TaskWithArgs(message string, num int) {
	log.Println("task with args is excuted. message: num:", message, num)
}
func main() {
	storage := store.NewMemoryStorage()
	s := go_scheduler.NewScheduler(storage)
	if _, err := s.RunEvery(1*time.Second, TaskWithArgs, "test", 1); err != nil {
		log.Fatal(err)
	}
	s.Start()
	s.Wait()
}
