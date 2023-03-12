package todo

import (
	"context"
	"github.com/suresh024/MyTodo/models"
	store "github.com/suresh024/MyTodo/repository"
	"github.com/suresh024/MyTodo/repository/todo"
	"time"
)

type todoServ struct {
	todoRepo todo.Repository
}

func New(store store.Store) Service {
	return &todoServ{
		todoRepo: store.TodoRepo,
	}
}

type Service interface {
	GetTodoByIDServ(ctx context.Context, todoIds []string) ([]models.Todo, error)
	CreateTodoServ(ctx context.Context, todo models.Todo) (models.Todo, error)
	GetTodos(ctx context.Context) ([]models.Todo, error)
	UpdateTodoServ(ctx context.Context, todo models.Todo) (models.Todo, error)
}

func (s *todoServ) GetTodoByIDServ(ctx context.Context, todoIds []string) ([]models.Todo, error) {
	myTodos, err := s.todoRepo.GetTodoByID(ctx, todoIds)
	if err != nil {
		return myTodos, err
	}
	return myTodos, nil
}

func (s *todoServ) GetTodos(ctx context.Context) ([]models.Todo, error) {
	myTodos, err := s.todoRepo.GetAllTodos(ctx)
	if err != nil {
		return myTodos, err
	}
	return myTodos, nil
}

func (s *todoServ) CreateTodoServ(ctx context.Context, todo models.Todo) (models.Todo, error) {
	currentTime := int(time.Now().Unix())
	todo.Audit.CreatedAt, todo.Audit.UpdateAt = currentTime, currentTime
	myTodos, err := s.todoRepo.CreateTodo(ctx, todo)
	if err != nil {
		return myTodos, err
	}
	return myTodos, nil
}

func (s *todoServ) UpdateTodoServ(ctx context.Context, todo models.Todo) (models.Todo, error) {
	currentTime := int(time.Now().Unix())
	todo.Audit.UpdateAt = currentTime
	myTodos, err := s.todoRepo.UpdateTodo(ctx, todo)
	if err != nil {
		return myTodos, err
	}
	return myTodos, nil
}
