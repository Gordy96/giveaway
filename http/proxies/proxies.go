package proxies

import (
	"giveaway/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type ProxyService interface {
	AddMany([]string)
	Add(string)
	GetNext() string
	IncrementErrorCounter(string)
}

type StringProxyService struct {
	ProxyList []string
	Mux       sync.Mutex
	Index     int
	length    int
}

func NewStringProxyService() ProxyService {
	return &StringProxyService{
		ProxyList: []string{},
		Mux:       sync.Mutex{},
		Index:     0,
	}
}

func NewWithList(proxies []string) ProxyService {
	return &StringProxyService{
		ProxyList: proxies,
		Mux:       sync.Mutex{},
		Index:     0,
		length:    len(proxies),
	}
}

func (p *StringProxyService) AddMany(proxies []string) {
	p.Mux.Lock()
	defer p.Mux.Unlock()
	p.ProxyList = append(p.ProxyList, proxies...)
	p.length = len(p.ProxyList)
}

func (p *StringProxyService) Add(proxy string) {
	p.Mux.Lock()
	defer p.Mux.Unlock()
	p.ProxyList = append(p.ProxyList, proxy)
	p.length = len(p.ProxyList)
}

func (p *StringProxyService) GetNext() string {
	p.Mux.Lock()
	proxy := p.ProxyList[p.Index]
	if p.Index < p.length-1 {
		p.Index++
	} else {
		p.Index = 0
	}
	p.Mux.Unlock()
	return proxy
}
func (p *StringProxyService) IncrementErrorCounter(proxy string) {

}

type DataBaseProxyService struct {
	col *mongo.Collection
	mux sync.Mutex
}

type proxyEntry struct {
	Proxy      string `json:"proxy" bson:"proxy"`
	LastUsedAt int64  `json:"last_used_at" bson:"last_used_at"`
	ErrorCount int    `json:"error_count" bson:"error_count"`
}

func (d *DataBaseProxyService) AddMany(proxies []string) {
	l := len(proxies)
	docs := make([]interface{}, l)
	for i := 0; i < l; i++ {
		docs[i] = proxyEntry{
			Proxy:      proxies[i],
			LastUsedAt: time.Now().UnixNano(),
			ErrorCount: 0,
		}
	}
	d.col.InsertMany(nil, docs)
}

func (d *DataBaseProxyService) Add(proxy string) {
	doc := proxyEntry{
		Proxy:      proxy,
		LastUsedAt: time.Now().UnixNano(),
		ErrorCount: 0,
	}
	d.col.InsertOne(nil, doc)
}

func (d *DataBaseProxyService) GetNext() string {
	d.mux.Lock()
	defer d.mux.Unlock()
	opt := &options.FindOneAndUpdateOptions{}
	opt.SetSort(bson.M{
		"last_used_at": 1,
	})
	upd := bson.M{
		"$set": bson.M{
			"last_used_at": time.Now().UnixNano(),
		},
	}
	r := d.col.FindOneAndUpdate(nil, bson.M{}, upd, opt)
	doc := proxyEntry{}
	r.Decode(&doc)
	return doc.Proxy
}
func (d *DataBaseProxyService) IncrementErrorCounter(proxy string) {
	opt := &options.UpdateOptions{}
	opt.SetUpsert(true)
	d.col.UpdateOne(nil, bson.M{
		"proxy": proxy,
	}, bson.M{
		"$inc": bson.M{
			"error_count": 1,
		},
	}, opt)
}

func NewDatabaseProxyService() ProxyService {
	return &DataBaseProxyService{
		mux: sync.Mutex{},
		col: utils.Database().Collection("Proxies"),
	}
}

var globalProxyConveyor ProxyService = NewDatabaseProxyService()

func GetGlobalInstance() ProxyService {
	return globalProxyConveyor
}
