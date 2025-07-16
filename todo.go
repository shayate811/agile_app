package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"time"
	"strconv"
)

type Task struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Done         bool   `json:"done"`
	SprintNumber int    `json:"sprint_number,omitempty"` // Optional field for sprint number
	TaskWeight   int    `json:"task_weight,omitempty"`   // Optional field for task weight
}

const dataFile = "todo.json"
const timersettingFile = "timer_setting.json"

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

func loadTimerSettings() (map[string]int, error) {
	file, err := os.Open(timersettingFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No settings file, return nil
		}
		return nil, err
	}
	defer file.Close()

	var settings map[string]int
	if err := json.NewDecoder(file).Decode(&settings); err != nil {
		return nil, err
	}
	return settings, nil
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

func AddTask(title string, sprintNumber int, taskWeight int) {
	tasks, err := loadTasks()
	if err != nil {
		panic(err)
	}

	newTask := Task{
		ID:           nextID(tasks),
		Title:        title,
		Done:         false,
		SprintNumber: sprintNumber, // Default value for sprint number
		TaskWeight:   taskWeight,   // Default value for task weight
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Sprint_Number", "Task_Weight", "Status"})

	for _, task := range tasks {
		status := "[ ]"
		if task.Done {
			status = "[x]"
		}
		row := []string{
			strconv.Itoa(task.ID),
			task.Title,
			strconv.Itoa(task.SprintNumber),
			strconv.Itoa(task.TaskWeight),
			status,
		}
		table.Append(row)
	}

	table.Render()
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

func TimerStartSprint() {
    //jsonの読み込み
    settings, err := loadTimerSettings()
    if err != nil {
        panic(err)
    }	
    if settings == nil {
        // デフォルトのタイマー設定を使用
        settings = map[string]int{
            "planning":    15, // スプリント計画: 15分
            "development": 60, // 開発: 60分
            "review":      15, // スプリントレビュー＋振り返り: 15分
        }
        fmt.Println("タイマー設定ファイルが見つからないため、デフォルト値を使用します。")
    } else {
        fmt.Printf("スプリント計画: %d分, 開発: %d分, スプリントレビュー＋振り返り: %d分\n",
            settings["planning"], settings["development"], settings["review"])
    }

	fmt.Printf("スプリントプランニング（%d分）を開始します\n", settings["planning"])
    timerMinutes(settings["planning"])
    fmt.Println("スプリント計画が終了しました")

    fmt.Println("開発（%d分）を開始します", settings["development"])
    timerMinutes(settings["development"])
    fmt.Println("開発が終了しました")

    fmt.Println("スプリントレビュー＋振り返り（%d分）を開始します", settings["review"])
    timerMinutes(settings["review"])
    fmt.Println("スプリントレビュー＋振り返りが終了しました")

    fmt.Println("=== スプリントタイムボックス終了 ===")
}

// 分数を受け取ってタイマー表示する補助関数
func timerMinutes(min int) {
    totalSec := min * 60
    for i := totalSec; i > 0; i-- {
        fmt.Printf("\r残り: %d分%d秒", i/60, i%60)
        time.Sleep(1 * time.Second)
    }
    fmt.Println("\nタイマー終了")
}

func TimerSetting(planningTime, developmentTime, reviewTime int) {
	timerSettings := map[string]int{
		"planning":  planningTime,
		"development": developmentTime,
		"review":     reviewTime,
	}

	file, err := os.Create(timersettingFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(timerSettings); err != nil {
		panic(err)
	}
}