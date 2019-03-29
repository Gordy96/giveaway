package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type ConnectionOptions struct {
	DBName string
	DBHost string
	DBPort int
}

var defaultOptions = &ConnectionOptions{
	DBHost: "localhost",
	DBName: "giveaway",
	DBPort: 27017,
}

func mergeOptions(opts ...*ConnectionOptions) *ConnectionOptions {
	d := *defaultOptions
	for _, o := range opts {
		if o.DBName != "" {
			d.DBName = o.DBName
		}
		if o.DBHost != "" {
			d.DBHost = o.DBHost
		}
		if o.DBPort > 0 {
			d.DBPort = o.DBPort
		}
	}
	return &d
}

func compareOptions(opts ...*ConnectionOptions) bool {
	if cl == nil {
		return false
	}
	for _, o := range opts {
		if instanceOpts.DBPort != o.DBPort {
			return false
		}
		if instanceOpts.DBHost != o.DBHost {
			return false
		}
		if instanceOpts.DBName != o.DBName {
			return false
		}
	}
	return true
}

var cl *mongo.Client = nil
var instanceOpts *ConnectionOptions = nil

func Connect(opts ...*ConnectionOptions) *mongo.Client {
	o := mergeOptions(opts...)
	cl, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", o.DBHost, o.DBPort)))
	if err != nil {
		panic(err)
	}
	cl.Connect(nil)
	return cl
}

func Instance(opts ...*ConnectionOptions) *mongo.Client {
	mux := sync.Mutex{}
	mux.Lock()
	defer mux.Unlock()
	if !compareOptions(opts...) {
		if cl != nil {
			cl.Disconnect(nil)
		}
		instanceOpts = mergeOptions(opts...)
		cl = Connect(instanceOpts)
	}
	return cl
}

func Database(db ...string) *mongo.Database {
	cl := Instance()
	if len(db) == 1 {
		return cl.Database(db[0])
	}
	return cl.Database(instanceOpts.DBName)
}
