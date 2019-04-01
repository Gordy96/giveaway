package utils

import (
	"giveaway/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type Repository struct {
	col *mongo.Collection
}

func (r *Repository) FindCommentsTaskById(id interface{}) (*data.CommentsTask, error) {
	a := &data.CommentsTask{}
	e := r.col.FindOne(nil, bson.M{"_id": id}).Decode(a)
	if e != nil {
		return nil, e
	}
	return a, nil
}

func (r *Repository) FindHashTagTaskById(id interface{}) (*data.HashTagTask, error) {
	a := &data.HashTagTask{}
	e := r.col.FindOne(nil, bson.M{"_id": id}).Decode(a)
	if e != nil {
		return nil, e
	}
	return a, nil
}

func (r *Repository) Save(task data.HasKey) error {
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
	mux := sync.Mutex{}
	mux.Lock()
	repo, ok := repos["Tasks"]
	if !ok {
		repo = NewTasksRepository("Tasks")
		repos["Tasks"] = repo
	}
	mux.Unlock()
	return repo
}

var repos map[string]*Repository = make(map[string]*Repository)

func GetNamedTasksRepositoryInstance(name string) *Repository {
	mux := sync.Mutex{}
	mux.Lock()
	repo, ok := repos[name]
	if !ok {
		repo = NewTasksRepository(name)
		repos[name] = repo
	}
	mux.Unlock()
	return repo
}
