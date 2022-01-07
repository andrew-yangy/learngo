package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	fmt.Println("Frontend server started")
	//resp, err := http.Get("http://localhost:8080")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer resp.Body.Close()
	//b, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println(err)
	//}
	tmpl := template.Must(template.ParseFiles("services/frontend/templates/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := TodoPageData{
			PageTitle: "My TODO list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
			},
		}

		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":3000", nil)
}
