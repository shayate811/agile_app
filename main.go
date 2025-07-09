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
		title := os.Args[2]
		AddTask(title)
	case "list":
		ListTasks()
	case "complete":
		id, _ := strconv.Atoi(os.Args[2])
		CompleteTask(id)
	case "delete":
		id, _ := strconv.Atoi(os.Args[2])
		DeleteTask(id)
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
