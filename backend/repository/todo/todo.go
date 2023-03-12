package todo

import (
	"context"
	"fmt"
	"github.com/lucsky/cuid"
	"github.com/suresh024/MyTodo/consts"
	"github.com/suresh024/MyTodo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Init(client *mongo.Client) Repository {
	repo := client.Database(consts.Database).Collection(consts.TodoCollection)
	return &todoRepo{
		todo: repo,
	}
}

type todoRepo struct {
	todo *mongo.Collection
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

	return myTodos, nil
}

func (r *todoRepo) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	myTodos := make([]models.Todo, 0)

	cursor, err := r.todo.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	if err != nil {
		return myTodos, err
	}

	err = cursor.All(ctx, &myTodos)
	if err != nil {
		return myTodos, err
	}

	return myTodos, nil
}

func (r *todoRepo) CreateTodo(ctx context.Context, todo models.Todo) (models.Todo, error) {
	todo.ID = cuid.New()
	_, err := r.todo.InsertOne(ctx, todo)
	if err != nil {
		return todo, err
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
	return todo, nil
}
