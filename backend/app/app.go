package app

import (
	"context"
	handl "github.com/suresh024/MyTodo/handler"
	todo_handl "github.com/suresh024/MyTodo/handler/todo"
	store "github.com/suresh024/MyTodo/repository"
	"github.com/suresh024/MyTodo/repository/todo"
	serv "github.com/suresh024/MyTodo/service"
	todo_serv "github.com/suresh024/MyTodo/service/todo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var repos store.Store
var service serv.Store
var h handl.Store

func dbSetup(url string, ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	return client, nil
}
func handlerSetup() {
	h = handl.Store{
		TodoHandler: todo_handl.New(service),
	}
}
func serviceSetup(repo store.Store) {
	service = serv.Store{
		TodoServ: todo_serv.New(repos),
	}
}
func repoSetup(client *mongo.Client) {
	repos = store.Store{
		TodoRepo: todo.Init(client),
	}
}

func Start() {
	mongoUrl := os.Getenv("mongo_url")
	envPort := os.Getenv("port")
	host := os.Getenv("host")
	ctx := context.Background()
	client, err := dbSetup(mongoUrl, ctx)
	defer client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	repoSetup(client)
	serviceSetup(repos)
	handlerSetup()
	runserver(host, envPort, h)
}
