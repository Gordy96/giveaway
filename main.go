package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"giveaway/client/validation"
	"giveaway/client/validation/rules"
	"giveaway/client/web"
	"giveaway/data"
	"giveaway/data/errors"
	"giveaway/http/requests"
	"giveaway/http/responses"
	"giveaway/instagram/solver"
	"giveaway/utils"
	"giveaway/utils/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"math/rand"
	"time"
)

type Entry struct {
	Key   string
	Value interface{}
}

type RandomEntryTask struct {
	data      []Entry
	dupes     map[string][]int
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

func (t *RandomEntryTask) LengthNoDuplicates() int {
	return len(t.dupes)
}

func filterWinner(ret *RandomEntryTask, rules []validation.IRule) (int, interface{}, error) {
	max := ret.LengthNoDuplicates()
	i := 0
	for {
		winnerId := ret.GetRandomIndexNoDuplicates()
		temp := ret.Get(winnerId).Value
		shouldChoose := true
		for _, rule := range rules {
			ruleResult, err := rule.Validate(temp)
			if err != nil {
				return -1, nil, err
			}
			if !ruleResult {
				shouldChoose = false
				break
			}
		}
		if shouldChoose {
			return winnerId, ret.Get(winnerId).Value, nil
		}
		i++
		if i >= max {
			return -1, nil, nil
		}
	}
}

func filterWinnerHashTag(ret *RandomEntryTask, rules []validation.IRule) (*data.TagMedia, error) {
	_, t, e := filterWinner(ret, rules)
	if t != nil {
		return t.(*data.TagMedia), nil
	}
	return nil, e
}

func filterWinnerComment(ret *RandomEntryTask, rules []validation.IRule) (int, *data.Comment, error) {
	i, t, e := filterWinner(ret, rules)
	if t != nil {
		return i, t.(*data.Comment), e
	}
	return -1, nil, e
}

func execPosts(task *data.HashTagTask, rules validation.RuleCollection) {
	cl := web.NewWebClient(&utils.UserAgentGenerator{}, "http://localhost:8888")
	cl.Init()

	ret := RandomEntryTask{make([]Entry, 0), make(map[string][]int, 0), func(e interface{}) string {
		return e.(*data.TagMedia).Owner.Id
	}}

	var err error = nil
	err = cl.QueryTag(task.HashTag, func(media data.TagMedia) (bool, error) {
		var shouldAdd = true
		for _, rule := range rules.AppendRules() {
			shouldAdd, err = rule.Validate(&media)
			if err != nil {
				switch /*e := */ err.(type) {
				case errors.ShouldStopIterationError:
					return false, err
				default:
					return false, err
				}
			}
			if !shouldAdd {
				break
			}
		}
		if shouldAdd {
			ret.Add(&media)
		}
		return true, nil
	})
	if err != nil {
		logger.DefaultLogger().Errorf("error: %v", err)
	}
	winner, err := filterWinnerHashTag(&ret, rules.SelectRules())

	if err != nil {
		logger.DefaultLogger().Errorf("error: %v", err)
	}

	if winner != nil {
		task.Post = winner
		task.Status = "complete"
	} else {
		task.Status = "cannot_decide_winner"
	}
	err = utils.GetNamedTasksRepositoryInstance("HashTagTasks").Save(task)
	if err != nil {
		logger.DefaultLogger().Errorf("error: %v", err)
	}
}

func execComments(task *data.CommentsTask, rules validation.RuleCollection) {
	cl := web.NewWebClient(&utils.UserAgentGenerator{}, "http://localhost:8888")
	cl.Init()

	ret := RandomEntryTask{make([]Entry, 0), make(map[string][]int, 0), func(e interface{}) string {
		return e.(*data.Comment).Owner.Id
	}}
	var err error = nil
	err = cl.QueryComments(task.ShortCode, func(comment data.Comment) (bool, error) {
		var shouldAdd = true
		for _, rule := range rules.AppendRules() {
			shouldAdd, err = rule.Validate(&comment)
			if err != nil {
				switch /*e := */ err.(type) {
				case errors.ShouldStopIterationError:
					return false, err
				default:
					return false, err
				}
			}
			if !shouldAdd {
				break
			}
		}
		if shouldAdd {
			ret.Add(&comment)
		}
		return true, nil
	})

	if err != nil {
		logger.DefaultLogger().Errorf("error: %v", err)
	}

	winnerId, winner, err := filterWinnerComment(&ret, rules.SelectRules())

	if err != nil {
		logger.DefaultLogger().Errorf("error: %v", err)
	}

	if winner != nil {
		above := make([]*data.Comment, 0)
		below := make([]*data.Comment, 0)

		for i := winnerId - 1; i >= 0 && i >= winnerId-2; i-- {
			above = append([]*data.Comment{ret.Get(i).Value.(*data.Comment)}, above...)
		}
		for i := winnerId + 1; i < len(ret.data) && i <= winnerId+2; i++ {
			below = append(below, ret.Get(i).Value.(*data.Comment))
		}

		task.Winner = winner
		task.Above = above
		task.Below = below

		task.Position = winnerId
		task.Status = "complete"
	} else {
		task.Status = "cannot_decide_winner"
	}

	err = utils.GetNamedTasksRepositoryInstance("CommentTasks").Save(task)
	if err != nil {
		logger.DefaultLogger().Errorf("error: %v", err)
	}
}

func main() {
	validation.RegisterRuleConstructor(validation.RuleConstructorMap{
		"FollowsRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			idi := i.(map[string]interface{})["id"].(interface{})
			rule := &rules.FollowsRule{"FollowsRule", idi.(string)}
			return validation.SelectRule, rule
		},
	})
	solv := solver.GetInstance()
	solv.Run()
	//repo := repository.GetRepositoryInstance()
	//for {
	//	ac := repo.GetOldestUsedRetries(5, 2 * time.Second)
	//	cl := api.NewApiClientWithAccount(ac)
	//	o := data.Owner{"4119227113", "ozcan198865"}
	//	id := "25025320"
	//	i, err := cl.IsFollower(&o, id)
	//	switch err.(type) {
	//	case errors.LoginRequired:
	//		solver.GetRunningInstance().Enqueue(commands.MakeNewReLoginCommand(ac))
	//	}
	//	if i {
	//		fmt.Printf("%s do follows %s", o.Username, id)
	//	} else {
	//		fmt.Printf("%s do not follows %s", o.Username, id)
	//	}
	//}
	//
	//return
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
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}
						res, err := utils.GetNamedTasksRepositoryInstance("CommentTasks").FindCommentsTaskById(bsonx.ObjectID(id))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, responses.NewSuccessfulCommentsTaskJsonResponse(*res))
					})
					comments.POST("/", func(c *gin.Context) {

						var req requests.CommentTaskJsonRequest
						err := c.BindJSON(&req)

						if err != nil {
							c.JSON(400, responses.NewValidationErrorJsonResponse())
							return
						}
						task := data.CommentsTask{}
						task.ShortCode = req.ShortCode
						task.Status = "in_progress"
						task.Id = primitive.NewObjectID()
						err = utils.GetNamedTasksRepositoryInstance("CommentTasks").Save(&task)

						if err != nil {
							panic(err)
						}

						go execComments(&task, req.Rules)
						c.JSON(200, responses.NewSuccessfulCommentsTaskJsonResponse(task))
					})
				}
				posts := tasks.Group("/posts")
				{
					posts.GET("/:id", func(c *gin.Context) {
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						res, err := utils.GetNamedTasksRepositoryInstance("HashTagTasks").FindHashTagTaskById(bsonx.ObjectID(id))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, responses.NewSuccessfulHashTagTaskJsonResponse(*res))
					})
					posts.POST("/", func(c *gin.Context) {

						var req requests.HashTagTaskJsonRequest
						err := c.BindJSON(&req)

						if err != nil {
							c.JSON(400, responses.NewValidationErrorJsonResponse())
							return
						}
						task := data.HashTagTask{}
						task.HashTag = req.HashTag
						task.Status = "in_progress"
						task.Id = primitive.NewObjectID()
						err = utils.GetNamedTasksRepositoryInstance("HashTagTasks").Save(&task)
						if err != nil {
							panic(err)
						}

						go execPosts(&task, req.Rules)
						c.JSON(200, responses.NewSuccessfulHashTagTaskJsonResponse(task))
					})
				}
			}
		}
	}
	app.Run("0.0.0.0:80")
}
