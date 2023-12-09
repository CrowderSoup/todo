package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/crowdersoup/todo/store"
)

const STORE_KEY string = "todos"

type Todo struct {
	ID   string
	Text string
	Done bool
}

type addTodoForm struct {
	Text string
}

type urlParams struct {
	ID string `uri:"id"`
}

type markAsDoneForm struct {
	ID string
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

		var todos []Todo
		t, _ := s.Get(STORE_KEY)
		if t == nil {
			todos = []Todo{}
		} else {
			todos = t.([]Todo)
		}

		err := s.Set(STORE_KEY, todos)
		if err != nil {
			panic(err)
		}

		c.HTML(http.StatusOK, "todos.html", gin.H{
			STORE_KEY: todos,
		})
	})

	router.POST("/todos", func(c *gin.Context) {
		s := store.NewStore(c)
		var form addTodoForm

		c.Bind(&form)

		var todos []Todo
		t, _ := s.Get(STORE_KEY)
		if t == nil {
			todos = []Todo{}
		} else {
			todos = t.([]Todo)
		}

		todos = append(todos, Todo{
			ID:   uuid.NewString(),
			Text: form.Text,
		})

		// Ensure we save the current list of todos
		err := s.Set(STORE_KEY, todos)
		if err != nil {
			panic(err)
		}

		c.HTML(http.StatusOK, "todos.html", gin.H{
			STORE_KEY: todos,
		})
	})

	router.POST("/todos/:id", func(c *gin.Context) {
		var params urlParams
		if err := c.ShouldBindUri(&params); err != nil {
			panic("bad request")
		}
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
