package utils

import (
	"giveaway/data/tasks"
	"giveaway/data/сontainers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type Repository struct {
	col *mongo.Collection
}

func (r *Repository) FindCommentsTaskById(id interface{}) (*tasks.CommentsTask, error) {
	a := &tasks.CommentsTask{}
	e := r.col.FindOne(nil, bson.M{"_id": id}).Decode(a)
	if e != nil {
		return nil, e
	}
	return a, nil
}

func (r *Repository) FindHashTagTaskById(id interface{}) (*tasks.HashTagTask, error) {
	a := &tasks.HashTagTask{}
	e := r.col.FindOne(nil, bson.M{"_id": id}).Decode(a)
	if e != nil {
		return nil, e
	}
	return a, nil
}

func (r *Repository) FindHashTagStoryTaskById(id interface{}) (*tasks.StoriesTask, error) {
	a := &tasks.StoriesTask{}
	e := r.col.FindOne(nil, bson.M{"_id": id}).Decode(a)
	if e != nil {
		return nil, e
	}
	return a, nil
}

func (r *Repository) Save(task сontainers.HasKey) error {
	opts := options.UpdateOptions{}
	opts.SetUpsert(true)
	_, err := r.col.UpdateOne(nil, bson.M{"_id": task.GetKey()}, bson.M{"$set": task}, &opts)
	return err
}

func NewTasksRepository(table string) *Repository {
	r := &Repository{}
	r.col = Database().Collection(table)

	return r
}

func GetTasksRepositoryInstance() *Repository {
	return GetNamedTasksRepositoryInstance("Tasks")
}

var repos map[string]*Repository = make(map[string]*Repository)
var taskRepositorySingletonMux = sync.Mutex{}

func GetNamedTasksRepositoryInstance(name string) *Repository {
	taskRepositorySingletonMux.Lock()
	repo, ok := repos[name]
	if !ok {
		repo = NewTasksRepository(name)
		repos[name] = repo
	}
	taskRepositorySingletonMux.Unlock()
	return repo
}
