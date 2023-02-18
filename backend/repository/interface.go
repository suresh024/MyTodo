package repository

import "github.com/suresh024/MyTodo/repository/todo"

type Store struct {
	TodoRepo todo.Repository
}
