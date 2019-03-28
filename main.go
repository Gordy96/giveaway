package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"giveaway/client/api"
	"giveaway/client/web"
	"giveaway/data"
	"giveaway/data/errors"
	"giveaway/instagram/account"
	"giveaway/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"math/rand"
	"time"
)

type Entry struct {
	Key		string
	Value 	interface{}
}

type RandomEntryTask struct {
	data 		[]Entry
	dupes		map[string][]int
	keyGetter func(interface{}) string
}

func (t *RandomEntryTask) Get(i int) *Entry {
	if i < 0 {
		panic(fmt.Errorf("index < 0"))
	}
	return &t.data[i]
}

func (t *RandomEntryTask) GetRandomIndexNoDuplicates() int {
	l := len(t.data) - 1
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomizer.Seed(time.Now().UnixNano())
	var idx = -1
	for {
		idx = randomizer.Intn(l)
		entry := t.data[idx]
		if dupes := t.dupes[entry.Key]; len(dupes) == 1 && dupes[0] == idx {
			break
		}
	}
	return idx
}

func (t *RandomEntryTask) Add(value interface{}) {
	idx := len(t.data)

	entry := Entry{t.keyGetter(value), value}
	t.data = append(t.data, entry)
	if dup, in := t.dupes[entry.Key]; in {
		t.dupes[entry.Key] = append([]int{idx}, dup...)
	} else {
		t.dupes[entry.Key] = []int{idx}
	}
}

func (t *RandomEntryTask) Length() int {
	return len(t.data)
}

var logger *utils.Logger = nil

func GetLogger() *utils.Logger {
	if logger == nil {
		logger = utils.NewFileLogger()
	}
	return logger
}

func execPosts(task *data.HashTagTask, db *mongo.Database, rules utils.RuleCollection) {
	cl := web.NewWebClient(&utils.UserAgentGenerator{},"http://localhost:8888")
	cl.Init()

	ret := RandomEntryTask{make([]Entry, 0), make(map[string][]int, 0), func(e interface{}) string {
		return e.(*data.TagMedia).Owner.Id
	}}

	cl.QueryTag(task.HashTag, func(media data.TagMedia) bool {
		var shouldAdd = true
		var err error = nil
		for _, rule := range rules.Data() {
			shouldAdd, err = rule.Validate(&media)
			switch /*e := */err.(type) {
			case errors.ShouldStopIterationError:
				return false
			}
		}
		if shouldAdd {
			ret.Add(&media)
		}
		return true
	})
	winnerId := ret.GetRandomIndexNoDuplicates()
	winner := ret.Get(winnerId).Value.(*data.TagMedia)

	task.Post = winner
	task.Status = "complete"

	res, err := db.Collection("HashTagTasks").UpdateOne(nil, bson.M{"_id": bsonx.ObjectID(task.Id)}, bson.M{ "$set": task })
	if res.ModifiedCount == 0 {
		panic(fmt.Errorf("not found"))
	}
	if err != nil {
		panic(err)
	}
}

func execComments(task *data.CommentsTask, db *mongo.Database, rules utils.RuleCollection) {
	cl := web.NewWebClient(&utils.UserAgentGenerator{},"http://localhost:8888")
	cl.Init()

	ret := RandomEntryTask{make([]Entry, 0), make(map[string][]int, 0), func(e interface{}) string {
		return e.(*data.Comment).Owner.Id
	}}
	cl.QueryComments(task.ShortCode, func(comment data.Comment) bool {
		var shouldAdd = true
		var err error = nil
		for _, rule := range rules.Data() {
			shouldAdd, err = rule.Validate(&comment)
			switch /*e := */err.(type) {
			case errors.ShouldStopIterationError:
				return false
			}
		}
		if shouldAdd {
			ret.Add(&comment)
		}
		return true
	})
	winnerId := ret.GetRandomIndexNoDuplicates()
	winner := ret.Get(winnerId).Value.(*data.Comment)
	above := make([]*data.Comment, 0)
	below := make([]*data.Comment, 0)

	for i := winnerId - 1; i >= 0 && i >= winnerId - 2; i-- {
		above = append([]*data.Comment{ret.Get(i).Value.(*data.Comment)}, above...)
	}
	for i := winnerId + 1; i < len(ret.data) && i <= winnerId + 2; i++ {
		below = append(below, ret.Get(i).Value.(*data.Comment))
	}

	task.Winner = winner
	task.Above = above
	task.Below = below

	task.Position = winnerId
	task.Status = "complete"

	res, err := db.Collection("CommentTasks").UpdateOne(nil, bson.M{"_id": bsonx.ObjectID(task.Id)}, bson.M{ "$set": task })
	if res.ModifiedCount == 0 {
		panic(fmt.Errorf("not found"))
	}
	if err != nil {
		panic(err)
	}
}

type HasRulesJsonRequest struct {
	Rules	utils.RuleCollection	`json:"rules"`
}

type CommentTaskJsonRequest struct {
	HasRulesJsonRequest				`json:",inline" bson:",inline"`
	ShortCode	string				`json:"shortcode" bson:"shortcode"`
}

type HashTagTaskJsonRequest struct {
	HasRulesJsonRequest				`json:",inline" bson:",inline"`
	HashTag		string				`json:"hashtag" bson:"hashtag"`
}

type HasStatusJsonResponse struct {
	Status		bool				`json:"status" bson:"status"`
}

type NotFoundJsonResponse struct {
	HasStatusJsonResponse			`json:",inline"`
	Error		string				`json:"error"`
}

func NewNotFoundJsonResponse() NotFoundJsonResponse {
	r := NotFoundJsonResponse{}
	r.Status = false
	r.Error = "err_not_found"
	return r
}

type ValidationErrorJsonResponse struct {
	HasStatusJsonResponse			`json:",inline"`
	Error		string				`json:"error"`
}

func NewValidationErrorJsonResponse() ValidationErrorJsonResponse {
	r := ValidationErrorJsonResponse{}
	r.Status = false
	r.Error = "request_validation_error"
	return r
}

type SuccessfulCommentsTaskJsonResponse struct {
	HasStatusJsonResponse			`json:",inline"`
	Result		data.CommentsTask	`json:"result"`
}

func NewSuccessfulCommentsTaskJsonResponse(task data.CommentsTask) SuccessfulCommentsTaskJsonResponse {
	r := SuccessfulCommentsTaskJsonResponse{}
	r.Status = false
	r.Result = task
	return r
}

type SuccessfulHashTagTaskJsonResponse struct {
	HasStatusJsonResponse			`json:",inline"`
	Result		data.HashTagTask	`json:"result"`
}

func NewSuccessfulHashTagTaskJsonResponse(task data.HashTagTask) SuccessfulHashTagTaskJsonResponse {
	r := SuccessfulHashTagTaskJsonResponse{}
	r.Status = false
	r.Result = task
	return r
}

func SafeSend(ch chan int, value int) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = false
		}
	}()
	ch <- value  // panic if ch is closed
	return true // <=> closed = false; return
}

func main() {
	//GetLogger()
	//client := web.NewWebClient(&utils.UserAgentGenerator{},"http://localhost:8888")
	//exec("BvZUiaKps7N", client, task)
	//acc := account.NewAccount("sasai@protonmail.com", "1Qqwerty")
	//client := web.NewWebClient(&utils.UserAgentGenerator{}, GetLogger())
	//client.SetAccount(acc)
	//client.Init()
	//client.Login()
	acc := account.NewAccount("johndoe8365", "123qwerty")
	acc.DeviceId = "android-3815aa3061e066c6"
	acc.AdId = "ebb78380-02c4-4111-abf3-76a1960daf30"
	acc.GUID = "3809a356-7663-48ab-9454-3f5c97928253"
	acc.PhoneId = "01fe7828-9bad-467f-878d-57322c1d6337"
	client := api.NewApiClient("http://localhost:8888")
	client.SetAccount(acc)
	//client.QeSync()
	//client.LauncherSync()
	//client.Login()
	client.IsFollower(data.Owner{"3532042922", ""}, "25025320")

	return

	conn, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn.Connect(ctx)
	db := conn.Database("giveaway")
	if err != nil {
		panic(err)
	}
	app := gin.Default()
	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			tasks := v1.Group("/tasks")
			{
				comments := tasks.Group("/comments")
				{
					comments.GET("/:id", func(c *gin.Context) {
						task := data.CommentsTask{}
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, NewNotFoundJsonResponse())
							return
						}

						res := db.Collection("CommentTasks").FindOne(nil, bson.M{"_id": bsonx.ObjectID(id)})
						err = res.Decode(&task)

						if err != nil {
							c.JSON(404, NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, NewSuccessfulCommentsTaskJsonResponse(task))
					})
					comments.POST("/", func(c *gin.Context) {

						var req CommentTaskJsonRequest
						err := c.BindJSON(&req)

						if err != nil {
							c.JSON(400, NewValidationErrorJsonResponse())
							return
						}
						task := data.CommentsTask{}
						task.ShortCode = req.ShortCode
						task.Status = "in_progress"
						task.Id = primitive.NewObjectID()
						_ ,err = db.Collection("CommentTasks").InsertOne(nil, task)

						if err != nil {
							panic(err)
						}

						go execComments(&task, db, req.Rules)
						c.JSON(200, NewSuccessfulCommentsTaskJsonResponse(task))
					})
				}
				posts := tasks.Group("/posts")
				{
					posts.GET("/:id", func(c *gin.Context) {
						task := data.HashTagTask{}
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, NewNotFoundJsonResponse())
							return
						}

						res := db.Collection("HashTagTasks").FindOne(nil, bson.M{"_id": bsonx.ObjectID(id)})
						err = res.Decode(&task)

						if err != nil {
							c.JSON(404, NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, NewSuccessfulHashTagTaskJsonResponse(task))
					})
					posts.POST("/", func(c *gin.Context) {

						var req HashTagTaskJsonRequest
						err := c.BindJSON(&req)

						if err != nil {
							c.JSON(400, NewValidationErrorJsonResponse())
							return
						}
						task := data.HashTagTask{}
						task.HashTag = req.HashTag
						task.Status = "in_progress"
						task.Id = primitive.NewObjectID()
						_ ,err = db.Collection("HashTagTasks").InsertOne(nil, task)

						if err != nil {
							panic(err)
						}

						go execPosts(&task, db, req.Rules)
						c.JSON(200, NewSuccessfulHashTagTaskJsonResponse(task))
					})
				}
			}
		}
	}
	app.Run("0.0.0.0:80")
}
