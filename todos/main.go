package todos

import (
	"github.com/crowdersoup/todo/store"
)

const STORE_KEY string = "todos"

type Todo struct {
	ID   string
	Text string
	Done bool
}

type Todoer interface {
	GetAll() map[string]Todo
	Get(string) Todo
	AddOrUpdate(Todo)
}

type todoer struct {
	store store.Store
}

func NewTodoer(s store.Store) Todoer {
	return &todoer{
		store: s,
	}
}

func (t todoer) GetAll() map[string]Todo {
	var todos map[string]Todo
	rawTodos, _ := t.store.Get(STORE_KEY)
	if rawTodos != nil {
		todos = rawTodos.(map[string]Todo)
	} else {
		todos = map[string]Todo{}
		t.store.Set(STORE_KEY, todos)
	}

	return todos
}

func (t todoer) Get(id string) Todo {
	todos := t.GetAll()

	return todos[id]
}

func (t todoer) AddOrUpdate(todo Todo) {
	todos := t.GetAll()

	// Add todo to the map
	todos[todo.ID] = todo

	// Save the updated map
	t.store.Set(STORE_KEY, todos)
}
