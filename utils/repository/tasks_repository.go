package repository

import (
	"giveaway/data"
	"giveaway/data/сontainers"
	"giveaway/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type Repository struct {
	col *mongo.Collection
}

func (r *Repository) FindTaskById(id interface{}, target interface{}) error {
	return r.col.FindOne(nil, bson.M{"_id": id}).Decode(target)
}

func (r *Repository) FindStory(id interface{}) (*data.StoryContainer, error) {
	a := &data.StoryContainer{}
	e := r.col.FindOne(nil, bson.M{"_id": id}).Decode(a)
	if e != nil {
		return nil, e
	}
	return a, nil
}

func (r *Repository) FindAll(criteria bson.M) (*mongo.Cursor, error) {
	c, e := r.col.Find(nil, criteria)
	if e != nil {
		return nil, e
	}
	return c, nil
}

func (r *Repository) Save(s сontainers.HasKey) error {
	opts := options.UpdateOptions{}
	opts.SetUpsert(true)
	_, err := r.col.UpdateOne(nil, bson.M{"_id": s.GetKey()}, bson.M{"$set": s}, &opts)
	return err
}

func (r *Repository) Insert(s сontainers.HasKey) (int64, error) {
	opts := options.UpdateOptions{}
	opts.SetUpsert(true)
	ur, err := r.col.UpdateOne(nil, bson.M{"_id": s.GetKey()}, bson.M{"$set": s}, &opts)
	return ur.UpsertedCount, err
}

func (r *Repository) DeleteOne(s сontainers.HasKey) error {
	_, err := r.col.DeleteOne(nil, bson.M{"_id": s.GetKey()})
	return err
}

func (r *Repository) DeleteMany(criteria bson.M) error {
	_, err := r.col.DeleteOne(nil, criteria)
	return err
}

func NewTasksRepository(table string) *Repository {
	r := &Repository{}
	r.col = utils.Database().Collection(table)

	return r
}

func GetTasksRepositoryInstance() *Repository {
	return GetNamedRepositoryInstance("Tasks")
}

var repos map[string]*Repository = make(map[string]*Repository)
var taskRepositorySingletonMux = sync.Mutex{}

func GetNamedRepositoryInstance(name string) *Repository {
	taskRepositorySingletonMux.Lock()
	repo, ok := repos[name]
	if !ok {
		repo = NewTasksRepository(name)
		repos[name] = repo
	}
	taskRepositorySingletonMux.Unlock()
	return repo
}
