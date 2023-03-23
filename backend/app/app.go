package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	redis "github.com/redis/go-redis/v9"
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
	"strconv"
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
func cacheSetup(ctx context.Context, url string) (*redis.Client, error) {
	cacheOptions, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(cacheOptions)
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
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
func repoSetup(client *mongo.Client, redis *redis.Client, redisRequired bool) {
	repos = store.Store{
		TodoRepo: todo.Init(client, redis, redisRequired),
	}
}

func Start() {
	// env declarations
	mongoUrl := os.Getenv("mongo_url")
	envPort := os.Getenv("port")
	host := os.Getenv("host")
	redisUrl := os.Getenv("redis_url")
	isredisRequired := os.Getenv("redis_required")
	ctx := context.Background()

	//connections setup
	//mongo connection
	mongoClient, err := dbSetup(mongoUrl, ctx)
	defer mongoClient.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// redis connection
	redisClient, err := cacheSetup(ctx, redisUrl)
	if err != nil {
		log.Fatal(err)
	}
	redisRequired, err := strconv.ParseBool(isredisRequired)
	if err != nil {
		log.Fatal(err)
	}

	//graphDb setup
	driver, err := neo4j.NewDriverWithContext("", neo4j.BasicAuth("", "", ""))
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)
	result, err := session.Run(ctx, "match n return n", map[string]interface{}{"message": "done"})
	if err != nil {
		fmt.Printf(err.Error())
	}
	some := make([]interface{}, 0)
	v, _ := json.Marshal(result.Record().Values)
	json.Unmarshal(v, &some)

	//dependency  setup
	repoSetup(mongoClient, redisClient, redisRequired)
	serviceSetup(repos)
	handlerSetup()
	runserver(host, envPort, h)
}
