package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/crowdersoup/todo/store"
	"github.com/crowdersoup/todo/todos"
)

type addTodoForm struct {
	Text string
}

type urlParams struct {
	ID string `uri:"id"`
}

func main() {
	router := gin.Default()

	// Set up session store middleware
	router.Use(store.InitSession())

	// Routes for favicon / stylesheet
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.StaticFile("/style.css", "./resources/style.css")

	// Load templates
	router.LoadHTMLGlob("./templates/*")

	// Index Route
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/todos", func(c *gin.Context) {
		s := store.NewStore(c)
		t := todos.NewTodoer(s)

		c.HTML(http.StatusOK, "todos.html", gin.H{
			todos.STORE_KEY: t.GetAll(),
		})
	})

	router.POST("/todos", func(c *gin.Context) {
		s := store.NewStore(c)
		t := todos.NewTodoer(s)

		var form addTodoForm
		c.Bind(&form)

		newTodo := todos.Todo{
			ID:   uuid.NewString(),
			Text: form.Text,
		}

		t.AddOrUpdate(newTodo)

		c.HTML(http.StatusOK, "todos.html", gin.H{
			todos.STORE_KEY: t.GetAll(),
		})
	})

	router.PATCH("/todos/:id", func(c *gin.Context) {
		var params urlParams
		if err := c.ShouldBindUri(&params); err != nil {
			panic("bad request")
		}

		s := store.NewStore(c)
		t := todos.NewTodoer(s)

		todo := t.Get(params.ID)
		todo.Done = !todo.Done

		t.AddOrUpdate(todo)

		c.HTML(http.StatusOK, "todos.html", gin.H{
			todos.STORE_KEY: t.GetAll(),
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
