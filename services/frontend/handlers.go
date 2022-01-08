package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func (fe *frontendServer) homeHandler(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("services/frontend/templates/index.html"))

	data := TodoPageData{
		PageTitle: fe.getOrders(),
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}

	tmpl.Execute(c.Writer, data)
}
