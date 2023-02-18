package app

import (
	"fmt"
	"github.com/gorilla/mux"
	handl "github.com/suresh024/MyTodo/handler"
	"log"
	"net/http"
)

func runserver(host, port string, h handl.Store) {
	r := mux.NewRouter()
	//r := router.Host(host).Subrouter()
	r = r.StrictSlash(true)

	//health check
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "alive")
	}).Methods("GET")

	r.HandleFunc("/todo/create", h.TodoHandler.CreateTodo).Methods("POST")
	r.HandleFunc("/todo/get/{todo_id}", h.TodoHandler.GetTodoByID).Methods("GET")
	r.HandleFunc("/todo/getall", h.TodoHandler.GetAllTodos).Methods("GET")
	r.HandleFunc("/todo/{todo_id}", h.TodoHandler.UpdateTodo).Methods("PUT")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
