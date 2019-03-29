package repository

import (
	"giveaway/instagram/account"
	"giveaway/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type Repository struct {
	col *mongo.Collection
}

func (r *Repository) FindById(id string) *account.Account {
	a := &account.Account{}
	r.col.FindOne(nil, bson.M{"_id": id}).Decode(a)
	return a
}

func (r *Repository) FindByUsername(username string) *account.Account {
	a := &account.Account{}
	r.col.FindOne(nil, bson.M{"username": username}).Decode(a)
	return a
}

func (r *Repository) Save(acc *account.Account) error {
	opts := options.UpdateOptions{}
	opts.SetUpsert(true)
	_, err := r.col.UpdateOne(nil, bson.M{"username": acc.Username}, bson.M{"$set": acc}, &opts)
	return err
}

func (r *Repository) GetRandom() *account.Account {
	ac := &account.Account{}
	cursor, err := r.col.Aggregate(nil, []interface{}{
		bson.M{"$sample": bson.M{"size": 1}},
	})
	if err == nil {
		cursor.Decode(ac)
	}
	return ac
}

func (r *Repository) GetOldestUsed() *account.Account {
	ac := &account.Account{}
	mux := sync.Mutex{}
	mux.Lock()
	defer mux.Unlock()
	filter := bson.M{
		"status": bson.M{
			"$nin": []account.AccountStatus{
				account.Maintenance,
				account.CheckPoint,
				account.Error,
			},
		},
	}
	res := r.col.FindOne(nil, filter, &options.FindOneOptions{Sort: bson.M{"updated_at": 1}})
	if res.Err() != nil {
		panic(res.Err())
	}

	res.Decode(ac)
	if ac.Username == "" {
		return nil
	}
	ac.UpdatedAt = time.Now().UnixNano()
	r.col.UpdateOne(nil, bson.M{"username": ac.Username}, bson.M{"$set": bson.M{"updated_at": ac.UpdatedAt}})
	return ac
}

var defaultRepo *Repository = nil

func NewRepository(table string) *Repository {
	r := &Repository{}
	r.col = utils.Database().Collection(table)

	return r
}

func GetRepositoryInstance() *Repository {
	mux := sync.Mutex{}
	mux.Lock()
	if defaultRepo == nil {
		defaultRepo = NewRepository("Accounts")
	}
	mux.Unlock()
	return defaultRepo
}
