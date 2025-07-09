package main

import (
	"encoding/json"
	"fmt"
	"os"
	"fmt"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

const dataFile = "todo.json"

func loadTasks() ([]Task, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var tasks []Task
	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(tasks)
}

func nextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

func AddTask(title string) {
	tasks, err := loadTasks()
	if err != nil {
		panic(err)
	}

	newTask := Task{
		ID:    nextID(tasks),
		Title: title,
		Done:  false,
	}
	tasks = append(tasks, newTask)

	if err := saveTasks(tasks); err != nil {
		panic(err)
	}
}

func ListTasks() {
	tasks, err := loadTasks()
	if err != nil {
		panic(err)
	}

	for _, task := range tasks {
		status := "[ ]"
		if task.Done {
			status = "[x]"
		}
		fmt.Printf("%d: %s %s\n", task.ID, task.Title, status)
	}
}

func CompleteTask(id int) {
	tasks, err := loadTasks()
	if err != nil {
		panic(err)
	}

	isExist := false

	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Done = true
			isExist = true
			break
		}
	}

	if !isExist {
		fmt.Print("task not found")
	}

	if err := saveTasks(tasks); err != nil {
		panic(err)
	}
}

func DeleteTask(id int) {
	tasks, err := loadTasks()
	if err != nil {
		panic(err)
	}

	newTasks := make([]Task, 0, len(tasks))
	isExist := false

	for _, t := range tasks {
		if t.ID == id {
			isExist = true
			continue // このタスクを削除（スキップ）
		}
		newTasks = append(newTasks, t)
	}

	if !isExist {
		fmt.Println("task not found")
		return
	}

	if err := saveTasks(newTasks); err != nil {
		panic(err)
	}

}
