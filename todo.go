package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/benoitmasson/plotters/piechart"
	"github.com/olekukonko/tablewriter"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"hash/fnv"
	"image/color"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
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

type Timer struct {
	Plannning    int `json:"planning"`
	Development  int `json:"development"`
	Review       int `json:"review"`
	SprintNumber int `json:"sprint_number"`
}

const dataFile = "todo.json"
const timersettingFile = "timer_setting.json"

// ColorFromName は人名を安定した色(RGBA)に変換します。
// 同じ name を渡せば常に同じ色が返ります。
func ColorFromName(name string) color.RGBA {
	// 1) 64bit FNV ハッシュ
	h := fnv.New64a()
	h.Write([]byte(name))
	hash := h.Sum64()

	// 2) ハッシュ値 → 0–359 の Hue
	hue := float64(hash % 360)
	sat := 0.55   // 彩度
	light := 0.60 // 輝度

	return hslToRGBA(hue, sat, light)
}

// --- 内部関数: HSL → RGBA -----------------------------------------
func hslToRGBA(h, s, l float64) color.RGBA {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2
	var r, g, b float64
	switch {
	case 0 <= h && h < 60:
		r, g, b = c, x, 0
	case 60 <= h && h < 120:
		r, g, b = x, c, 0
	case 120 <= h && h < 180:
		r, g, b = 0, c, x
	case 180 <= h && h < 240:
		r, g, b = 0, x, c
	case 240 <= h && h < 300:
		r, g, b = x, 0, c
	default: // 300‑360
		r, g, b = c, 0, x
	}
	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}

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

func loadTimerSettings() (*Timer, error) {
	file, err := os.Open(timersettingFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No settings file, return nil
		}
		return nil, err
	}
	defer file.Close()

	var settings *Timer
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

func ListDoingTasks(sprint int) {
	tasks, err := loadTasks()
	if err != nil {
		panic(err)
	}

	// フィルタ & グループ分け
	todo := []Task{}
	doing := []Task{}
	done := []Task{}
	for _, task := range tasks {
		if task.SprintNumber <= sprint {
			if task.Done == true {
				done = append(done, task)
			} else if strings.TrimSpace(task.Assignees) == "" {
				todo = append(todo, task)
			} else {
				doing = append(doing, task)
			}
		}
	}

	// ソート（どちらも同様に）
	sortTasks := func(ts []Task) {
		sort.Slice(ts, func(i, j int) bool {
			if ts[i].SprintNumber == ts[j].SprintNumber {
				return ts[i].ID < ts[j].ID
			}
			return ts[i].SprintNumber < ts[j].SprintNumber
		})
	}
	sortTasks(todo)
	sortTasks(doing)
	sortTasks(done)

	// 表示関数
	renderTable := func(title string, ts []Task) {
		fmt.Printf("\n=== %s ===\n", title)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Title", "Sprint", "Weight", "Assignees", "Status"})
		for _, t := range ts {
			status := "[ ]"
			if t.Done {
				status = "[x]"
			}
			row := []string{
				strconv.Itoa(t.ID),
				t.Title,
				strconv.Itoa(t.SprintNumber),
				strconv.Itoa(t.TaskWeight),
				t.Assignees,
				status,
			}
			table.Append(row)
		}
		table.Render()
	}

	// 出力
	renderTable("Todo", todo)
	renderTable("Doing", doing)
	renderTable("Done", done)
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
		settings = &Timer{
			Plannning:    15,
			Development:  60,
			Review:       15,
			SprintNumber: 0,
		}
		fmt.Println("タイマー設定ファイルが見つからないため、デフォルト値を使用します。")
	} else {
		fmt.Printf("スプリント番号 : %d, スプリント計画: %d分, 開発: %d分, スプリントレビュー＋振り返り: %d分\n",
			settings.SprintNumber, settings.Plannning, settings.Development, settings.Review)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 安全のため

	// ── 入力受付を並列実行
	go listenInput(ctx, cancel)

	go func() {

		fmt.Printf("スプリントプランニング（%d分）を開始します\n", settings.Plannning)
		timerMinutes(settings.Plannning)
		fmt.Println("スプリント計画が終了しました")

		fmt.Println("開発（%d分）を開始します", settings.Development)
		ListDoingTasks(settings.SprintNumber)
		timerMinutes(settings.Development)
		fmt.Println("開発が終了しました")

		fmt.Println("スプリントレビュー＋振り返り（%d分）を開始します", settings.Review)
		timerMinutes(settings.Review)
		fmt.Println("スプリントレビュー＋振り返りが終了しました")

		fmt.Println("=== スプリントタイムボックス終了 ===")

		settings.SprintNumber += 1

		file, err := os.Create(timersettingFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(settings); err != nil {
			panic(err)
		}
		cancel()

	}()

	<-ctx.Done()

}

// 分数を受け取ってタイマー表示する補助関数
func timerMinutes(min int) {
	totalSec := min * 60
	for i := totalSec; i > 0; i-- {
		fmt.Fprintf(os.Stderr, "\r[タイマー] 残り: %2d分%02d秒", i/60, i%60)
		time.Sleep(1 * time.Second)
	}
	fmt.Fprintln(os.Stderr, "\n[タイマー] タイマー終了")
}

func TimerSetting(planningTime, developmentTime, reviewTime int) {
	settings, err := loadTimerSettings()
	if err != nil {
		panic(err)
	}

	timerSettings := Timer{
		Plannning:   planningTime,
		Development: developmentTime,
		Review:      reviewTime,
	}

	if settings == nil {
		timerSettings.SprintNumber = 1
	} else {
		timerSettings.SprintNumber = settings.SprintNumber
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
	p := plot.New()
	p.Title.Text = "Progress by Assignee"
	p.Y.Label.Text = "Progress (%)"
	p.NominalX(names...)

	// 各作業者ごとに1本ずつBarChartを重ねて色分け
	for i, name := range names {
		prog := progressMap[name]
		var rate float64
		if prog.totalWeight > 0 {
			rate = float64(prog.doneWeight) / float64(prog.totalWeight) * 100
		}
		vals := make(plotter.Values, len(names))
		vals[i] = rate // 他は0
		bar, err := plotter.NewBarChart(vals, vg.Points(30))
		if err != nil {
			fmt.Println("グラフ生成に失敗しました:", err)
			return
		}
		bar.LineStyle.Width = vg.Length(0)
		bar.Color = ColorFromName(name) // 名前から色を取得
		p.Add(bar)
	}
	p.Y.Max = 100

	// 横軸ラベルの角度を調整
	p.X.Tick.Label.Rotation = 0.5 // 0.5ラジアン（約30度）傾ける

	// 余白を設定
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

func ShowContribution() {
	// ==== 1. タスク読み込み & 集計 ==========================================
	tasks, err := loadTasks()
	if err != nil {
		log.Fatalf("loadTasks 失敗: %v", err)
	}

	contrib := make(map[string]float64)
	var unfinishedWeight float64

	for _, t := range tasks {
		if t.Done {
			name := t.Assignees
			if name == "" {
				name = "Unassigned"
			}
			contrib[name] += float64(t.TaskWeight)
		} else {
			unfinishedWeight += float64(t.TaskWeight)
		}
	}

	// ==== 2. 円グラフ用データ作成 ==========================================
	labels := make([]string, 0, len(contrib)+1)
	values := make(plotter.Values, 0, len(contrib)+1)

	for n, w := range contrib {
		labels = append(labels, n)
		values = append(values, w)
	}
	if unfinishedWeight > 0 {
		labels = append(labels, "Unfinished")
		values = append(values, unfinishedWeight)
	}

	// ==== 3. グラフベース生成 ==============================================
	p := plot.New()
	p.Title.Text = "Task Contribution"
	p.HideAxes() // 円グラフなので軸は非表示

	// ==== 4. スライスごとに PieChart を生成して色分け =======================
	total := 0.0
	for _, v := range values {
		total += v
	}

	offset := 0.0

	for i, v := range values {
		// 4‑1) 1スライスだけをもつ PieChart を生成
		pc, err := piechart.NewPieChart(plotter.Values{v})
		if err != nil {
			log.Fatalf("piechart 生成失敗: %v", err)
		}

		// 4‑2) 色と開始位置・合計値を設定
		pc.Color = ColorFromName(labels[i])
		pc.Offset.Value = offset
		pc.Total = total

		// 4‑3) ラベル表示設定
		pc.Labels.Show = true
		pc.Labels.Nominal = []string{labels[i]}
		pc.Labels.Values.Show = true
		pc.Labels.Values.Percentage = true // 割合表示

		// 4‑4) 追加
		p.Add(pc)
		offset += v
	}

	// ==== 5. 保存 ==========================================================
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "contribution.png"); err != nil {
		log.Fatalf("画像保存失敗: %v", err)
	}
	fmt.Println("貢献度円グラフを出力しました → contribution.png")
}

// defaultColors は必要数だけ色を返す簡易パレット
func defaultColors(n int) []color.Color {
	base := []color.Color{
		color.RGBA{255, 99, 132, 255},  // 赤
		color.RGBA{54, 162, 235, 255},  // 青
		color.RGBA{255, 206, 86, 255},  // 黄
		color.RGBA{75, 192, 192, 255},  // 緑
		color.RGBA{153, 102, 255, 255}, // 紫
		color.RGBA{255, 159, 64, 255},  // オレンジ
	}
	out := make([]color.Color, n)
	for i := 0; i < n; i++ {
		out[i] = base[i%len(base)]
	}
	return out
}

func listenInput(ctx context.Context, cancel context.CancelFunc) {
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nコマンド > ")
		if !sc.Scan() {
			cancel()
			return
		}
		inputs := strings.Split(sc.Text(), " ")
		if len(inputs) < 1 {
			fmt.Println("Help Command: [help]")
			return
		}

		switch inputs[0] {
		case "add":
			if len(inputs) < 4 {
				fmt.Println("Usage: add <title> <sprintNumber> <taskWeight>")
				return
			}
			title := inputs[1]
			sprintNumber, err1 := strconv.Atoi(inputs[2])
			taskWeight, err2 := strconv.Atoi(inputs[3])
			if err1 != nil || err2 != nil {
				fmt.Println("sprintNumberとtaskWeightは数値で指定してください")
				return
			}
			AddTask(title, sprintNumber, taskWeight)
		case "list":
			ListTasks()
		case "assign":
			id, _ := strconv.Atoi(inputs[1])
			name := ""
			if len(inputs) >= 3 {
				name = inputs[2]
			}
			AssignTask(id, name)
		case "complete":
			id, _ := strconv.Atoi(inputs[1])
			CompleteTask(id)
		case "delete":
			id, _ := strconv.Atoi(inputs[1])
			DeleteTask(id)
		case "exit":
			fmt.Println("exit sprint.")
			cancel()
			return
		case "help":
			fmt.Println("<Usage>\nAddTask : add <title> <sprintNumber> <taskWeight>\nListTasks :  list\nAssignTask : assign <TaskID> <UserName>\nCompleteTask : complete <TaskID>\nDeleteTask : delete <TaskID>\nExitSprint : exit ")
		default:
			fmt.Println("不明なコマンド")
		}

		// タイマーが先に終わっていないかチェック
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
