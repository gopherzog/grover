package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_version(t *testing.T) {
	v := version()
	expected := VERSION
	if v != expected {
		t.Errorf("Expected %s, received %s", expected, v)
	}
}

func Test_AddTaskOk(t *testing.T) {
	task := Task{
		Name:        "TestTask",
		Descr:       "Ensure we can create a task",
		StartDate:   time.Now(),
		PercentDone: 22,
	}
	var oldTasks []Task
	newTasks, err := AddTask(oldTasks, task)
	if err != nil {
		t.Errorf("Unable to AddTask: %v", err)
		return
	}

	taskCount := len(newTasks)
	if taskCount != 1 {
		t.Errorf("New Task not included in return value")
		return
	}

	actualTask := newTasks[0]
	if actualTask != task {
		t.Errorf("Added task not found: %s", task)
	}
	fmt.Printf("%s\n", actualTask)
}
