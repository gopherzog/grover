package main

import (
	_ "bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	APPNAME          string = "grover"
	VERSION          string = "0.0.3"
	defaultVerbosity bool   = false
	defaultLogPath   string = "logs/xlog.txt"
	defaultTaskPath  string = "tasks-out.dat"
)

func version() string {
	return VERSION
}

func OpenFileLogger(filename string) *log.Logger {
	outF := WriterForLogging(filename)
	return log.New(outF, "", log.LstdFlags)
}

func ShowPrompt() {
	fmt.Printf("Cmd> ")
}

func main() {
	var (
		loggerPath string
		verbose    bool
		taskPath   string
	)
	flag.StringVar(&taskPath, "taskpath", defaultTaskPath, "path for the task file")
	flag.StringVar(&loggerPath, "logpath", defaultLogPath, "path for the log file")
	flag.BoolVar(&verbose, "verbose", defaultVerbosity, "true for verbose logging")
	flag.Parse()
	xLog := OpenFileLogger(loggerPath)
	if verbose {
		xLog.Printf("%s %s\n", APPNAME, VERSION)
	}
	tasks, err := ReadTasks(taskPath)
	if err != nil {
		xLog.Printf("Unable to read tasks: %v", err)
		return
	}
	REPL(tasks, xLog, taskPath, verbose)
}

func REPL(tasks []Task, xLog *log.Logger, taskPath string, verbose bool) {
	for working := true; working; {
		ShowPrompt()
		command, args := ParseCommand()
		switch command {
		case "nop":
			fmt.Println("\n")

		case "exit", "quit", "q":
			working = false

		case "list":
			ListTasks(tasks)

		case "delete", "del":
			newTasks, err := HandleDeleteTask(args, tasks)
			tasks = HandlePersistTasks(tasks, newTasks, err, xLog, taskPath)

		case "modify", "mod", "update", "upd":
			newTasks, err := HandleModifyTask(args, tasks)
			tasks = HandlePersistTasks(tasks, newTasks, err, xLog, taskPath)

		case "create", "new":
			newTasks, err := HandleAddTask(args, tasks)
			tasks = HandlePersistTasks(tasks, newTasks, err, xLog, taskPath)
		}
	}
}

func WriterForLogging(logPath string) io.Writer {
	logf := os.Stderr
	if logPath != "" {
		os.MkdirAll(filepath.Dir(logPath), os.ModePerm)
		openedFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			logf = openedFile
		}
	}
	return logf
}
