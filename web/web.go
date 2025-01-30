package web

import (
	"encoding/json"
	"net/http"
	"os"
	"text/template"
)

type TodoItem struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

func StaticPage() {
	fs := http.FileServer(http.Dir("web/static/"))
	http.Handle("/about/", http.StripPrefix("/about/", fs))
	http.ListenAndServe(":8080", nil)
}


func DynamicPage() {
	http.HandleFunc("/list", listHandler)
	http.ListenAndServe(":8080", nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	todoItems := loadTodoItems()
	if todoItems == nil {
		http.Error(w, "Unable to load to-do items", http.StatusInternalServerError)
		return
	}

	// Define the template
	tmpl := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>To-Do List</title>
		</head>
		<body>
			<h1>To-Do List</h1>
			<ul>
				{{range .}}
					<li>
						<strong>{{.Desc}}</strong> (Status: {{.Status}})
					</li>
				{{else}}
					<li>No to-do items found</li>
				{{end}}
			</ul>
		</body>
		</html>
	`

	// Parse the template
	t, err := template.New("list").Parse(tmpl)
	if err != nil {
		http.Error(w, "Unable to parse template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, todoItems)
	if err != nil {
		http.Error(w, "Unable to execute template", http.StatusInternalServerError)
		return
	}
}

func loadTodoItems() []TodoItem {
	file, err := os.Open("list.json")
	if err != nil {
		return nil
	}
	defer file.Close()

	var todoItems []TodoItem
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&todoItems)
	if err != nil {
		return nil
	}
	return todoItems
}