package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Task struct {
	ID        int
	Title     string
	Completed bool
	DueDate   time.Time
}

func parseCommandLineArgs() (string, string, bool, int) {
	var title string
	var dueDate string
	var listTasksFlag bool
	var id int

	flag.StringVar(&title, "t", "", "Title of task")
	flag.StringVar(&dueDate, "d", "", "Specify the date in DD-MM-YYYY format")
	flag.BoolVar(&listTasksFlag, "l", false, "List all the tasks")
	flag.IntVar(&id, "i", 0, "Id of task")

	flag.Usage = func() {
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	return title, dueDate, listTasksFlag, id
}

func addTask(task *[]Task, title string, parsedDate time.Time) {
	newTask := Task{
		ID:        len(*task) + 1,
		Title:     title,
		Completed: false,
		DueDate:   parsedDate,
	}

	*task = append(*task, newTask)
	err := saveTask(*task)

	if err != nil {
		println("Error while saving data:", err)
		return
	}

	fmt.Printf("Task added: %s (Due: %s)\n", title, parsedDate.Format("02-01-2006"))
}

func loadTask() ([]Task, error) {
	data, err := os.ReadFile("tasks.json")

	if err != nil {
		return nil, err
	}

	var task []Task
	err = json.Unmarshal(data, &task)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func markTaskCompleted(id int) {
	tasks, err := loadTask()

	if err != nil {
		println("Error loading task:", err)
		return
	}

	for i, task := range tasks {
		if id == task.ID {
			task.Completed = true
			fmt.Printf("ID: %d | Title: %s | Due Date: %s | Completed: %v\n", task.ID, task.Title, task.DueDate.Format("02-01-2006"), task.Completed)
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTask(tasks)
			return
		}
	}
	println("Task not found with ID:", id)
}

func saveTask(tasks []Task) error {
	data, err := json.Marshal(tasks)

	if err != nil {
		return err
	}
	err = os.WriteFile("tasks.json", data, 0644)

	if err != nil {
		return err
	}

	return nil
}

func listTask() error {
	var tasks []Task
	tasks, err := loadTask()

	if err != nil {
		return err
	}

	for _, task := range tasks {
		fmt.Printf("ID: %d | Title: %s | Due Date: %s | Completed: %v\n", task.ID, task.Title, task.DueDate.Format("02-01-2006"), task.Completed)
	}
	return nil
}

func main() {

	title, dueDate, listTasksFlag, id := parseCommandLineArgs()

	if listTasksFlag {
		listTask()
	}

	dateLayout := "02-01-2006"
	var parsedDate time.Time

	if dueDate == "" {
		parsedDate = time.Now()
	} else {
		var err error
		parsedDate, err = time.Parse(dateLayout, dueDate)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if title == "" && listTasksFlag {
		err := listTask()
		if err != nil {
			print("Error while listing tasks:", err)
			return
		}
	}

	if title == "" && !listTasksFlag && id == 0 {
		println("Please provide title")
		flag.Usage()
		os.Exit(1)
	}

	if title != "" {
		task, err := loadTask()
		if err != nil {
			println("Error while loading tasks:", err)
		}
		addTask(&task, title, parsedDate)
	}

	if title == "" && id != 0 {
		markTaskCompleted(id)
	}
}
