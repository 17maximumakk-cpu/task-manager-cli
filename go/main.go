package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
}

const tasksFile = "tasks.json"

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tasksFile, data, 0644)
}

func colorize(text, color string) string {
	colors := map[string]string{
		"red":    "\033[91m",
		"green":  "\033[92m",
		"yellow": "\033[93m",
		"blue":   "\033[94m",
		"reset":  "\033[0m",
	}
	return colors[color] + text + colors["reset"]
}

func addTask(desc, priority, due string) {
	tasks, _ := loadTasks()
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	newID := maxID + 1
	newTask := Task{
		ID:          newID,
		Description: desc,
		Status:      "pending",
		Priority:    priority,
		DueDate:     due,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	tasks = append(tasks, newTask)
	saveTasks(tasks)
	fmt.Println(colorize(fmt.Sprintf("✓ Задача %d добавлена.", newID), "green"))
}

func listTasks(statusFilter string) {
	tasks, _ := loadTasks()
	if statusFilter != "" {
		filtered := []Task{}
		for _, t := range tasks {
			if t.Status == statusFilter {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}
	if len(tasks) == 0 {
		fmt.Println(colorize("Нет задач.", "yellow"))
		return
	}
	for _, t := range tasks {
		statusIcon := "◻"
		if t.Status == "done" {
			statusIcon = "✓"
		}
		prioColor := map[string]string{"high": "red", "medium": "yellow", "low": "green"}[t.Priority]
		due := ""
		if t.DueDate != "" {
			due = " [до " + t.DueDate + "]"
		}
		fmt.Printf("%s [%d] %s%s\n", statusIcon, t.ID, colorize(t.Description, prioColor), due)
	}
}

func completeTask(id int) {
	tasks, _ := loadTasks()
	found := false
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Status = "done"
			found = true
			break
		}
	}
	if !found {
		fmt.Println(colorize(fmt.Sprintf("Ошибка: задача %d не найдена.", id), "red"))
		return
	}
	saveTasks(tasks)
	fmt.Println(colorize(fmt.Sprintf("✓ Задача %d выполнена.", id), "green"))
}

func deleteTask(id int) {
	tasks, _ := loadTasks()
	newTasks := []Task{}
	found := false
	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue
		}
		newTasks = append(newTasks, t)
	}
	if !found {
		fmt.Println(colorize(fmt.Sprintf("Ошибка: задача %d не найдена.", id), "red"))
		return
	}
	saveTasks(newTasks)
	fmt.Println(colorize(fmt.Sprintf("✓ Задача %d удалена.", id), "green"))
}

func exportCSV(filename string) {
	tasks, _ := loadTasks()
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Ошибка создания CSV:", err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"id", "description", "status", "priority", "due_date", "created_at"})
	for _, t := range tasks {
		writer.Write([]string{
			strconv.Itoa(t.ID),
			t.Description,
			t.Status,
			t.Priority,
			t.DueDate,
			t.CreatedAt,
		})
	}
	fmt.Println(colorize(fmt.Sprintf("Экспортировано %d задач в %s", len(tasks), filename), "blue"))
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println(`
Доступные команды:
  add "описание" [--priority high|medium|low] [--due YYYY-MM-DD]
  list [--status pending|done]
  done <id>
  rm <id>
  export [filename.csv]
`)
		return
	}
	switch args[0] {
	case "add":
		if len(args) < 2 {
			fmt.Println("Укажите описание задачи")
			return
		}
		desc := args[1]
		priority := "medium"
		due := ""
		for i := 2; i < len(args); i++ {
			if args[i] == "--priority" && i+1 < len(args) {
				priority = args[i+1]
				i++
			}
			if args[i] == "--due" && i+1 < len(args) {
				due = args[i+1]
				i++
			}
		}
		addTask(desc, priority, due)
	case "list":
		status := ""
		if len(args) > 2 && args[1] == "--status" {
			status = args[2]
		}
		listTasks(status)
	case "done":
		if len(args) < 2 {
			fmt.Println("Укажите ID задачи")
			return
		}
		id, _ := strconv.Atoi(args[1])
		completeTask(id)
	case "rm":
		if len(args) < 2 {
			fmt.Println("Укажите ID задачи")
			return
		}
		id, _ := strconv.Atoi(args[1])
		deleteTask(id)
	case "export":
		filename := "tasks.csv"
		if len(args) > 1 {
			filename = args[1]
		}
		exportCSV(filename)
	default:
		fmt.Println("Неизвестная команда")
	}
}
