package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"os"
	"sort"
	"strconv"
	"time"
)

type Task struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Done         bool   `json:"done"`
	SprintNumber int    `json:"sprint_number,omitempty"` // Optional field for sprint number
	TaskWeight   int    `json:"task_weight,omitempty"`   // Optional field for task weight
	Assignees    string `json:"assignees"`
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
	table.SetHeader([]string{"ID", "Title", "Sprint_Number", "Task_Weight", "Assignees", "Status"})

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
			task.Assignees,
			status,
		}
		table.Append(row)
	}

	table.Render()
}

func AssignTask(id int, name string) {
	tasks, err := loadTasks()
	if err != nil {
		panic(err)
	}

	isExist := false

	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Assignees = name
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

func ShowProgress() {
    tasks, err := loadTasks()
    if err != nil {
        panic(err)
    }

    // assigneeごとに重みを集計
    type progress struct {
        doneWeight  int
        totalWeight int
    }
    progressMap := make(map[string]*progress)

    for _, task := range tasks {
        name := task.Assignees
        if name == "" {
            continue // 未割り当てタスクは集計しない
        }
        if _, ok := progressMap[name]; !ok {
            progressMap[name] = &progress{}
        }
        progressMap[name].totalWeight += task.TaskWeight
        if task.Done {
            progressMap[name].doneWeight += task.TaskWeight
        }
    }

    // 表示用にソート
    names := make([]string, 0, len(progressMap))
    for name := range progressMap {
        names = append(names, name)
    }
    sort.Strings(names)

    // テーブル表示
    fmt.Println("作業者\t完了重み/担当重み\t進捗率")
    fmt.Println("-------------------------------------")
    for _, name := range names {
        p := progressMap[name]
        rate := 0
        if p.totalWeight > 0 {
            rate = p.doneWeight * 100 / p.totalWeight
        }
        fmt.Printf("%s\t%d/%d\t\t%d%%\n", name, p.doneWeight, p.totalWeight, rate)
    }

    // グラフ用データ作成
    values := make(plotter.Values, len(names))
    for i, name := range names {
        p := progressMap[name]
        var rate float64
        if p.totalWeight > 0 {
            rate = float64(p.doneWeight) / float64(p.totalWeight) * 100
        }
        values[i] = rate
    }

    // グラフ生成
    p := plot.New()
    p.Title.Text = "Progress by Assignee"
    p.Y.Label.Text = "Progress (%)"

    // 横軸に作業者名を表示
    p.NominalX(names...)

    bar, err := plotter.NewBarChart(values, vg.Points(30))
    if err != nil {
        fmt.Println("グラフ生成に失敗しました:", err)
        return
    }
    bar.LineStyle.Width = vg.Length(0)
    p.Add(bar)
    p.Y.Max = 100

    // 横軸ラベルの角度を調整
    p.X.Tick.Label.Rotation = 0.5 // 0.5ラジアン（約30度）傾ける

    // 余白を設定（左右に0.2インチずつ余白を追加）
    p.X.Padding = vg.Points(40)
    p.X.Min = -0.5
    p.X.Max = float64(len(names)) - 0.5

    // グラフ画像として保存
    if err := p.Save(8*vg.Inch, 4*vg.Inch, "progress.png"); err != nil {
        fmt.Println("グラフ画像の保存に失敗しました:", err)
        return
    }
    fmt.Println("進捗グラフ(progress.png)を出力しました。")
}