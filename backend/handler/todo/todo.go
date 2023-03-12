package todo

import (
	"context"
	"encoding/json"
	"github.com/suresh024/MyTodo/models"
	"github.com/suresh024/MyTodo/service"
	"github.com/suresh024/MyTodo/service/todo"
	"github.com/suresh024/MyTodo/utils"
	"net/http"
	"strings"
	"time"
)

type todoHandler struct {
	todoServ todo.Service
}

func New(service service.Store) Handler {
	return &todoHandler{
		todoServ: service.TodoServ,
	}
}

type Handler interface {
	GetTodoByID(w http.ResponseWriter, r *http.Request)
	CreateTodo(w http.ResponseWriter, r *http.Request)
	GetAllTodos(w http.ResponseWriter, r *http.Request)
	UpdateTodo(w http.ResponseWriter, r *http.Request)
}

func (h *todoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Second)
	defer cancel()
	todoID, err := utils.GetUrlParam(r, "todo_id")
	if err != nil {
		utils.ErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	todoIDS := strings.Split(todoID, ",")
	myTodos, err := h.todoServ.GetTodoByIDServ(ctx, todoIDS)
	if err != nil {
		utils.ErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, myTodos)
	return
}

func (h *todoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	myTodos, err := h.todoServ.GetTodos(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, myTodos)
	return
}

func (h *todoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		utils.ErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	if val, ok := r.Header["Email"]; ok {
		todo.Audit.UpdateBy = val[0]
	}
	myTodos, err := h.todoServ.UpdateTodoServ(ctx, todo)
	if err != nil {
		utils.ErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, myTodos)
	return
}

func (h *todoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		utils.ErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	if val, ok := r.Header["Email"]; ok {
		todo.Audit.CreatedBy, todo.Audit.UpdateBy = val[0], val[0]
	}
	myTodos, err := h.todoServ.CreateTodoServ(ctx, todo)
	if err != nil {
		utils.ErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, myTodos)
	return
}
