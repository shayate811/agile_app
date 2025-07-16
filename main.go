package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo [add|list|complete|delete] ...")
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "add":
		if len(os.Args) < 5 {
			fmt.Println("Usage: todo add <title> <sprintNumber> <taskWeight>")
			return
		}
		title := os.Args[2]
		sprintNumber, err1 := strconv.Atoi(os.Args[3])
		taskWeight, err2 := strconv.Atoi(os.Args[4])
		if err1 != nil || err2 != nil {
			fmt.Println("sprintNumberとtaskWeightは数値で指定してください")
			return
		}
		AddTask(title, sprintNumber, taskWeight)
	case "list":
		ListTasks()
	case "assign":
		id, _ := strconv.Atoi(os.Args[2])
		name := ""
		if len(os.Args) >= 4 {
			name = os.Args[3]
		}
		AssignTask(id, name)
	case "complete":
		id, _ := strconv.Atoi(os.Args[2])
		CompleteTask(id)
	case "delete":
		id, _ := strconv.Atoi(os.Args[2])
		DeleteTask(id)
	case "timerstart":
		TimerStartSprint()
	case "timersetting":
		if len(os.Args) < 5 {
			fmt.Println("Usage: todo timersetting <planningTime> <developmentTime> <reviewTime>")
			return
		}
		planningTime, _ := strconv.Atoi(os.Args[2])
		developmentTime, _ := strconv.Atoi(os.Args[3])
		reviewTime, _ := strconv.Atoi(os.Args[4])
		TimerSetting(planningTime, developmentTime, reviewTime)
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
