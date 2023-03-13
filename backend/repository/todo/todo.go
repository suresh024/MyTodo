package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lucsky/cuid"
	"github.com/redis/go-redis/v9"
	"github.com/suresh024/MyTodo/consts"
	"github.com/suresh024/MyTodo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(db *mongo.Client, redis *redis.Client, redisRequired bool) Repository {
	repo := db.Database(consts.Database).Collection(consts.TodoCollection)
	return &todoRepo{
		todo:          repo,
		todoCache:     redis,
		redisRequired: redisRequired,
	}
}

type todoRepo struct {
	todo          *mongo.Collection
	todoCache     *redis.Client
	redisRequired bool
}

type Repository interface {
	GetTodoByID(ctx context.Context, todoID []string) ([]models.Todo, error)
	GetAllTodos(ctx context.Context) ([]models.Todo, error)
	CreateTodo(ctx context.Context, todo models.Todo) (models.Todo, error)
	UpdateTodo(ctx context.Context, todo models.Todo) (models.Todo, error)
}

func (r *todoRepo) GetTodoByID(ctx context.Context, todoID []string) ([]models.Todo, error) {
	myTodos := make([]models.Todo, 0)
	filter := bson.M{"id": bson.M{"$in": todoID}}

	if r.redisRequired {
		val, err := r.todoCache.HMGet(ctx, consts.Todo, todoID...).Result()
		if err == nil && len(val) > 0 {
			for _, value := range val {
				var todo models.Todo
				json.Unmarshal([]byte(value.(string)), &todo)
				myTodos = append(myTodos, todo)
			}
			return myTodos, nil
		}
	}
	cursor, err := r.todo.Find(ctx, filter)
	defer cursor.Close(ctx)
	fmt.Println(filter)
	if err != nil {
		return myTodos, err
	}

	err = cursor.All(ctx, &myTodos)
	if err != nil {
		return myTodos, err
	}
	if r.redisRequired {
		for _, value := range myTodos {
			valueData, _ := json.Marshal(value)
			r.todoCache.HMSet(ctx, consts.Todo, value.ID, string(valueData))
		}
	}

	return myTodos, nil
}

func (r *todoRepo) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	myTodos := make([]models.Todo, 0)

	if r.redisRequired {
		result, err := r.todoCache.HGetAll(ctx, consts.Todo).Result()
		if err == nil && len(result) > 0 {
			for _, value := range result {
				var todo models.Todo
				json.Unmarshal([]byte(value), &todo)
				myTodos = append(myTodos, todo)
			}
			return myTodos, nil
		}
	}
	cursor, err := r.todo.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	if err != nil {
		return myTodos, err
	}

	err = cursor.All(ctx, &myTodos)
	if err != nil {
		return myTodos, err
	}

	if r.redisRequired {
		for _, value := range myTodos {
			valueData, _ := json.Marshal(value)
			r.todoCache.HMSet(ctx, consts.Todo, value.ID, string(valueData))
		}

	}

	return myTodos, nil
}

func (r *todoRepo) CreateTodo(ctx context.Context, todo models.Todo) (models.Todo, error) {
	todo.ID = cuid.New()
	_, err := r.todo.InsertOne(ctx, todo)
	if err != nil {
		return todo, err
	}
	if r.redisRequired {
		value, err := json.Marshal(todo)
		if err == nil && len(value) > 0 {
			r.todoCache.HSet(ctx, consts.Todo, todo.ID, value)
		}
	}
	return todo, nil
}

func (r *todoRepo) UpdateTodo(ctx context.Context, todo models.Todo) (models.Todo, error) {
	filter := bson.M{"id": todo.ID}
	update := bson.D{{Key: "$set", Value: todo}}
	result, err := r.todo.UpdateOne(ctx, filter, update)
	if err != nil {
		return todo, err
	}
	if result.MatchedCount == 0 {
		return todo, fmt.Errorf("no documents with this id")
	}
	if r.redisRequired {
		value, err := json.Marshal(todo)
		if err == nil {
			r.todoCache.HMSet(ctx, consts.Todo, todo.ID, value)
		}
	}
	return todo, nil
}
