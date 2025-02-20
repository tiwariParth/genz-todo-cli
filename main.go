package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Todo struct {
    ID          int       `json:"id"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
}

var todos []Todo
var nextID = 1

type Streak struct {
    Count     int       `json:"count"`
    LastCheck time.Time `json:"last_check"`
}

var streak Streak

func loadTodos() {
    data, err := ioutil.ReadFile("todos.json")
    if err != nil {
        fmt.Println("No existing todos found. Starting fresh! 🌟")
        return
    }
    err = json.Unmarshal(data, &todos)
    if err != nil {
        fmt.Println("Error loading todos:", err)
        return
    }
    if len(todos) > 0 {
        nextID = todos[len(todos)-1].ID + 1
    }
}

func saveTodos() {
    data, err := json.Marshal(todos)
    if err != nil {
        fmt.Println("Error saving todos:", err)
        return
    }
    err = ioutil.WriteFile("todos.json", data, 0644)
    if err != nil {
        fmt.Println("Error writing todos to file:", err)
    }
}

func addTodo(description string) {
    todo := Todo{
        ID:          nextID,
        Description: description,
        Completed:   false,
        CreatedAt:   time.Now(),
    }
    todos = append(todos, todo)
    nextID++
    saveTodos()
    fmt.Printf("✅ Added: %s\n", description)
}

func listTodos() {
    if len(todos) == 0 {
        fmt.Println("🎉 No todos yet! You're all caught up!")
        return
    }
    fmt.Println("📋 Your Todos:")
    for _, todo := range todos {
        status := "❌"
        if todo.Completed {
            status = "✅"
        }
        fmt.Printf("%d. [%s] %s (%s)\n", todo.ID, status, todo.Description, todo.CreatedAt.Format("Jan 2"))
    }
}

func completeTodo(id int) {
    for i, todo := range todos {
        if todo.ID == id {
            todos[i].Completed = true
            saveTodos()
            fmt.Printf("🌟 Completed: %s\n", todo.Description)
            
            // Update streak
            loadStreak()
            updateStreak()
            
            return
        }
    }
    fmt.Println("🤔 Todo not found!")
}

func deleteTodo(id int) {
    for i, todo := range todos {
        if todo.ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            saveTodos()
            fmt.Printf("🗑️ Deleted: %s\n", todo.Description)
            return
        }
    }
    fmt.Println("🤔 Todo not found!")
}

func loadStreak() {
    data, err := ioutil.ReadFile("streak.json")
    if err != nil {
        fmt.Println("Starting new streak! 🔥")
        streak = Streak{Count: 0}
        return
    }
    err = json.Unmarshal(data, &streak)
    if err != nil {
        fmt.Println("Error loading streak:", err)
    }
}

func saveStreak() {
    data, err := json.Marshal(streak)
    if err != nil {
        fmt.Println("Error saving streak:", err)
        return
    }
    err = ioutil.WriteFile("streak.json", data, 0644)
    if err != nil {
        fmt.Println("Error writing streak to file:", err)
    }
}

func updateStreak() {
    now := time.Now()
    if now.Day() != streak.LastCheck.Day() {
        if now.Sub(streak.LastCheck).Hours() < 48 {
            streak.Count++
        } else {
            streak.Count = 1
        }
        streak.LastCheck = now
        saveStreak()
        fmt.Printf("🔥 Current Streak: %d days\n", streak.Count)
    }
}

func main() {
    loadTodos()
    loadStreak()

    if len(os.Args) < 2 {
        fmt.Println("Usage: genz-todo [add|list|complete|delete] [args]")
        return
    }

    command := os.Args[1]

    switch command {
    case "add":
        if len(os.Args) < 3 {
            fmt.Println("Usage: genz-todo add <description>")
            return
        }
        description := os.Args[2]
        addTodo(description)
    case "list":
        listTodos()
    case "complete":
        if len(os.Args) < 3 {
            fmt.Println("Usage: genz-todo complete <id>")
            return
        }
        id := parseInt(os.Args[2])
        completeTodo(id)
    case "delete":
        if len(os.Args) < 3 {
            fmt.Println("Usage: genz-todo delete <id>")
            return
        }
        id := parseInt(os.Args[2])
        deleteTodo(id)
    default:
        fmt.Println("Unknown command. Use 'add', 'list', 'complete', or 'delete'.")
    }
}

func parseInt(s string) int {
    var id int
    _, err := fmt.Sscanf(s, "%d", &id)
    if err != nil {
        fmt.Println("Invalid ID:", s)
        os.Exit(1)
    }
    return id
}