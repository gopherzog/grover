package main

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Task struct {
	Uuid         string    `json:"id"`
	Name         string    `json:"name"`
	Descr        string    `json:"descr"`
	CompleteDate time.Time `json:"complete_date"`
	StartDate    time.Time `json:"start_date"`
	PercentDone  int       `json:"percent_done"`
}

func NewTask(name, descr string) *Task {
	var task *Task
	id := uuid.NewV4()
	task = &Task{
		Uuid:        id.String(),
		Name:        name,
		Descr:       descr,
		StartDate:   time.Now(),
		PercentDone: 0,
	}
	return task
}

func (task Task) String() string {
	spacer := "       "
	return fmt.Sprintf("%s\n[%3.1d%%] %s\n%s%s\n%sStarted:   %s\n%sCompleted: %s",
		task.Uuid, task.PercentDone, task.Name,
		spacer, task.Descr,
		spacer, task.StartDate.Format(time.RFC1123),
		spacer, task.CompleteDate.Format(time.RFC1123))
}

func AddTask(tasks []Task, task Task) ([]Task, error) {
	var err error
	tasks = append(tasks, task)
	return tasks, err
}

func WriteTasks(tasks []Task, filename string) error {
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	openedFile, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		bx, err := json.Marshal(&tasks)
		if err != nil {
			return err
		}
		_, err = openedFile.Write(bx)
		if err != nil {
			return err
		}
		err = openedFile.Close()
	}
	return err
}

func ReadTasks(filename string) ([]Task, error) {
	var tasks []Task
	bx, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(bx, &tasks)
	}
	return tasks, err
}

func FindTaskById(tasks []Task, itemId string) (*Task, error) {
	for _, task := range tasks {
		if strings.HasPrefix(task.Uuid, itemId) {
			return &task, nil
		}
	}
	return nil, fmt.Errorf("task not found")
}

func DeleteTask(tasks []Task, task Task) ([]Task, error) {
	newTasks := make([]Task, 0)
	didFind := false
	for _, oldTask := range tasks {
		if oldTask.Uuid != task.Uuid {
			newTasks = append(newTasks, oldTask)
		} else {
			didFind = true
		}
	}
	var err error
	if !didFind {
		err = fmt.Errorf("Unable to update task")
	}
	return newTasks, err
}

func ReplaceTask(tasks []Task, task Task) ([]Task, error) {
	newTasks := make([]Task, 0)
	didFind := false
	for _, oldTask := range tasks {
		if oldTask.Uuid != task.Uuid {
			newTasks = append(newTasks, oldTask)
		} else {
			didFind = true
			newTasks = append(newTasks, task)
		}
	}
	var err error
	if !didFind {
		err = fmt.Errorf("Unable to update task")
	}
	return newTasks, err
}
