package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func ParseCommand() (string, []string) {
	args := make([]string, 0)
	consoleReader := bufio.NewReader(os.Stdin)
	rawLine, err := consoleReader.ReadString('\n')
	if err != nil {
		return "nop", args
	}
	commandLine := strings.TrimSpace(rawLine)
	allArgs := strings.Split(commandLine, " ")
	if len(allArgs) == 0 {
		return "nop", args
	}
	return allArgs[0], allArgs[1:]
}

func ListTasks(tasks []Task) {
	for _, task := range tasks {
		fmt.Println(task, "\n")
	}
}

func HandleAddTask(args []string, tasks []Task) ([]Task, error) {
	var name string
	var descr string
	if len(args) > 0 {
		name = args[0]
	}

	if len(args) > 1 {
		descr = strings.Join(args[1:], " ")
	}

	newTask := *NewTask(name, descr)
	newTasks, err := AddTask(tasks, newTask)
	if err != nil {
		return tasks, err
	}
	return newTasks, nil
}

func HandleDeleteTask(args []string, tasks []Task) ([]Task, error) {
	var itemId string
	if len(args) > 0 {
		itemId = args[0]
	}
	oldTask, err := FindTaskById(tasks, itemId)
	if err != nil {
		return tasks, err
	}
	newTasks, err := DeleteTask(tasks, *oldTask)
	if err != nil {
		return tasks, err
	}
	return newTasks, nil

}

func HandleModifyTask(args []string, tasks []Task) ([]Task, error) {
	var itemId string
	if len(args) > 0 {
		itemId = args[0]
	}
	oldTask, err := FindTaskById(tasks, itemId)
	if err != nil {
		return tasks, err
	}
	argLine := strings.Join(args[1:], " ")
	idx := strings.IndexRune(argLine, '=')
	if idx == -1 {
		return tasks, fmt.Errorf("Update expression requires '='")
	}
	fieldName := strings.ToLower(TextBefore(argLine, idx))
	fieldValue := TextAfter(argLine, idx)
	switch fieldName {
	case "name":
		oldTask.Name = fieldValue

	case "descr":
		oldTask.Descr = fieldValue

	case "percent", "pct":
		iValue, err := strconv.ParseInt(fieldValue, 10, 8)
		if err != nil {
			return tasks, fmt.Errorf("needs to be a number: %s", fieldValue)
		}
		if iValue < 0 || iValue > 100 {
			return tasks, fmt.Errorf("needs to be in the range 0..100 inclusive: %d", iValue)
		}
		oldTask.PercentDone = int(iValue)
	}
	newTasks, err := ReplaceTask(tasks, *oldTask)
	if err != nil {
		return tasks, err
	}
	return newTasks, err
}

func HandlePersistTasks(tasks []Task, newTasks []Task, err error, xLog *log.Logger, taskPath string) []Task {
	if err != nil {
		xLog.Printf("%s", err)
		return tasks
	}
	err = WriteTasks(newTasks, taskPath)
	if err != nil {
		xLog.Printf("Unable to write tasks")
		return tasks
	}
	return newTasks
}

func TextBefore(s string, idx int) string {
	return strings.TrimSpace(string([]rune(s)[:idx]))
}

func TextAfter(s string, idx int) string {
	from := idx + 1
	return strings.TrimSpace(string([]rune(s)[from:]))
}
