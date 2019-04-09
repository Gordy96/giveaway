package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"giveaway/client/api"
	"giveaway/client/validation"
	"giveaway/client/validation/rules"
	"giveaway/client/web"
	"giveaway/data"
	"giveaway/data/errors"
	"giveaway/data/tasks"
	"giveaway/data/сontainers"
	"giveaway/http/requests"
	"giveaway/http/responses"
	"giveaway/instagram/account/repository"
	"giveaway/instagram/solver"
	"giveaway/instagram/structures"
	"giveaway/utils"
	"giveaway/utils/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"io/ioutil"
	"net/http"
	"time"
)

func filterWinner(ret *сontainers.EntryContainer, rules []validation.IRule) (int, interface{}, error) {
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

func filterWinnerHashTag(ret *сontainers.EntryContainer, rules []validation.IRule, excluded []int) (int, *data.TagMedia, error) {
	var t interface{}
	var e error
	var i int
	for {
		i, t, e = filterWinner(ret, rules)
		var persisits = false
		for _, x := range excluded {
			if x == i {
				persisits = true
			}
		}
		if !persisits {
			break
		}
	}
	if t != nil {
		return i, t.(*data.TagMedia), nil
	}
	return -1, nil, e
}

func filterWinnerHashTagStories(ret *сontainers.EntryContainer, rules []validation.IRule, excluded []int) (int, *structures.StoryItem, error) {
	var t interface{}
	var e error
	var i int
	for {
		i, t, e = filterWinner(ret, rules)
		var persisits = false
		for _, x := range excluded {
			if x == i {
				persisits = true
			}
		}
		if !persisits {
			break
		}
	}
	if t != nil {
		return i, t.(*structures.StoryItem), nil
	}
	return -1, nil, e
}

func filterWinnerComment(ret *сontainers.EntryContainer, rules []validation.IRule, excluded []int) (int, *data.Comment, error) {
	var t interface{}
	var e error
	var i int
	for {
		i, t, e = filterWinner(ret, rules)
		var persisits = false
		for _, x := range excluded {
			if x == i {
				persisits = true
			}
		}
		if !persisits {
			break
		}
	}
	if t != nil {
		return i, t.(*data.Comment), e
	}
	return -1, nil, e
}

func runRules(rules []validation.IRule, arg interface{}, failCallback func(validation.IRule) error) (bool, error) {
	for _, rule := range rules {
		ruleRes, err := rule.Validate(arg)
		if err != nil {
			logger.DefaultLogger().Errorf("error: %v", err)
			return false, err
		}
		if !ruleRes {
			return false, failCallback(rule)
		}
	}
	return true, nil
}

func execStories(task *tasks.StoriesTask) {
	var err error = nil
	var ruleRes bool

	acc := repository.GetRepositoryInstance().FindByUsername(task.Account)
	cl := api.NewApiClientWithAccount(acc)
	cl.Login()
	repo := utils.GetNamedTasksRepositoryInstance("HashTagStoryTasks")
	dLogger := logger.DefaultLogger()
	clw := web.NewWebClient(&utils.UserAgentGenerator{}, acc.Proxy)
	ruleRes, err = runRules(task.Rules.PreconditionRules, clw, func(rule validation.IRule) error {
		task.Status = "cancelled"
		task.Comment = fmt.Sprintf("failed on precondition rule: %v", rule)
		err = repo.Save(task)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
		task.Status = "failed"
		task.Comment = fmt.Sprintf("got error: %v", err)
		err = repo.Save(task)
	}
	if !ruleRes {
		return
	}

	var ret = сontainers.NewEntryContainer()
	var story *structures.Story = nil
	for {
		watched := make([]structures.WatchedStoryEntry, 0)
		story, err = cl.QueryHashTagStories(task.HashTag, func(item structures.StoryItem) (bool, error) {
			var shouldAdd = true
			start := time.Now()
			for _, rule := range task.Rules.AppendingRules {
				shouldAdd, err = rule.Validate(&item)
				if err != nil {
					switch /*e := */ err.(type) {
					case errors.BeforeMinimumDate:
						return true, err
					case errors.AfterMaximumDate:
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
				ret.Add(&item)
			}
			time.Sleep(time.Duration(1 * float64(time.Second)))
			watched = append(watched, structures.WatchedStoryEntry{
				Item: item,
				Duration: [2]int64{
					start.Unix(),
					start.Unix() + int64(item.VideoDuration),
				},
			})
			go func(t *tasks.StoriesTask, i structures.StoryItem) {
				uri := i.ImageVersions2.Candidates[0].URL
				r, err := http.Get(uri)
				if err != nil {
					dLogger.Errorf("something bad happened while trying to get %s: %s", uri, err)
					return
				}
				content, _ := ioutil.ReadAll(r.Body)
				entry := data.StorySource{
					Id:        primitive.NewObjectID(),
					CreatedAt: time.Now().UnixNano(),
					TaskId:    task.Id,
					StoryId:   i.Pk,
					Data:      content,
				}
				utils.Database().Collection("StorySources").InsertOne(nil, entry)
			}(task, item)
			return true, nil
		})
		if len(watched) > 0 {
			cl.MarkHashTagStoriesAsSeen(story, watched...)
		}
		if story == nil {
			break
		}
		if err != nil {
			switch /*e := */ err.(type) {
			case errors.BeforeMinimumDate:
			case errors.AfterMaximumDate:
				break
			default:
				dLogger.Errorf("something happened: %s", err)
			}
		}
	}

	ruleRes, err = runRules(task.Rules.PostconditionRules, ret, func(rule validation.IRule) error {
		task.Status = "failed"
		task.Comment = fmt.Sprintf("failed on postcondition rule: %v", rule)
		err = repo.Save(task)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
	if !ruleRes {
		return
	}

	task.Winners = make([]structures.StoryItem, 0)
	taken := make([]int, 0)
	for i := 0; i < task.NumWinners; i++ {
		id, winner, err := filterWinnerHashTagStories(ret, task.Rules.SelectRules, taken)
		taken = append(taken, id)
		if err != nil {
			dLogger.Errorf("error: %v", err)
		}
		if winner != nil {
			task.Winners = append(task.Winners, *winner)
		} else {
			break
		}
	}

	if len(task.Winners) == task.NumWinners {
		task.Status = "complete"
	} else {
		task.Status = "cannot_decide_winner"
	}
	err = repo.Save(task)
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
	repository.GetRepositoryInstance().Release(acc)
}

func execPosts(task *tasks.HashTagTask) {
	cl := web.NewWebClient(&utils.UserAgentGenerator{}, "http://localhost:8888")
	cl.Init()

	var err error = nil
	var ruleRes bool
	repo := utils.GetNamedTasksRepositoryInstance("HashTagTasks")
	dLogger := logger.DefaultLogger()
	ruleRes, err = runRules(task.Rules.PreconditionRules, cl, func(rule validation.IRule) error {
		task.Status = "cancelled"
		task.Comment = fmt.Sprintf("failed on precondition rule: %v", rule)
		err = repo.Save(task)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
		task.Status = "failed"
		task.Comment = fmt.Sprintf("got error: %v", err)
		err = repo.Save(task)
	}
	if !ruleRes {
		return
	}

	ret := сontainers.NewEntryContainer()
	err = cl.QueryTag(task.HashTag, func(media data.TagMedia) (bool, error) {
		var shouldAdd = true
		for _, rule := range task.Rules.AppendingRules {
			shouldAdd, err = rule.Validate(&media)
			if err != nil {
				switch /*e := */ err.(type) {
				case errors.BeforeMinimumDate:
					return false, err
				case errors.AfterMaximumDate:
					return true, err
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
		dLogger.Errorf("error: %v", err)
	}

	ruleRes, err = runRules(task.Rules.PostconditionRules, ret, func(rule validation.IRule) error {
		task.Status = "failed"
		task.Comment = fmt.Sprintf("failed on postcondition rule: %v", rule)
		err = repo.Save(task)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
	if !ruleRes {
		return
	}

	task.Winners = make([]*data.TagMedia, 0)
	taken := make([]int, 0)
	for i := 0; i < task.NumWinners; i++ {

		id, winner, err := filterWinnerHashTag(ret, task.Rules.SelectRules, taken)
		taken = append(taken, id)
		if err != nil {
			dLogger.Errorf("error: %v", err)
		}
		if winner != nil {
			task.Winners = append(task.Winners, winner)
		} else {
			break
		}
	}

	if len(task.Winners) == task.NumWinners {
		task.Status = "complete"
	} else {
		task.Status = "cannot_decide_winner"
	}
	err = repo.Save(task)
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
}

func execComments(task *tasks.CommentsTask) {
	cl := web.NewWebClient(&utils.UserAgentGenerator{}, "http://localhost:8888")
	cl.Init()

	ret := сontainers.NewEntryContainer()
	repo := utils.GetNamedTasksRepositoryInstance("CommentTasks")
	dLogger := logger.DefaultLogger()
	var err error = nil
	var ruleRes bool
	ruleRes, err = runRules(task.Rules.PreconditionRules, cl, func(rule validation.IRule) error {
		task.Status = "cancelled"
		task.Comment = fmt.Sprintf("failed on precondition rule: %v", rule)
		err = repo.Save(task)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
		task.Status = "failed"
		task.Comment = fmt.Sprintf("got error: %v", err)
		err = repo.Save(task)
	}
	if !ruleRes {
		return
	}

	err = cl.QueryComments(task.ShortCode, func(comment data.Comment) (bool, error) {
		var shouldAdd = true
		for _, rule := range task.Rules.AppendingRules {
			shouldAdd, err = rule.Validate(&comment)
			if err != nil {
				switch /*e := */ err.(type) {
				case errors.BeforeMinimumDate:
					return false, err
				case errors.AfterMaximumDate:
					return true, err
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
		dLogger.Errorf("error: %v", err)
	}

	ruleRes, err = runRules(task.Rules.PostconditionRules, ret, func(rule validation.IRule) error {
		task.Status = "failed"
		task.Comment = fmt.Sprintf("failed on postcondition rule: %v", rule)
		err = repo.Save(task)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
	if !ruleRes {
		return
	}

	task.Winners = make([]data.CommentWinner, 0)
	taken := make([]int, 0)
	for i := 0; i < task.NumWinners; i++ {
		winnerId, winner, err := filterWinnerComment(ret, task.Rules.SelectRules, taken)
		taken = append(taken, winnerId)
		var wSet data.CommentWinner
		if err != nil {
			dLogger.Errorf("error: %v", err)
		}

		if winner != nil {
			above := make([]*data.Comment, 0)
			below := make([]*data.Comment, 0)

			for i := winnerId - 1; i >= 0 && i >= winnerId-2; i-- {
				above = append([]*data.Comment{ret.Get(i).Value.(*data.Comment)}, above...)
			}
			for i := winnerId + 1; i < ret.Length() && i <= winnerId+2; i++ {
				below = append(below, ret.Get(i).Value.(*data.Comment))
			}

			wSet.Winner = winner
			wSet.Above = above
			wSet.Below = below

			wSet.Position = winnerId
			task.Winners = append(task.Winners, wSet)
		} else {
			break
		}
	}

	if len(task.Winners) == task.NumWinners {
		task.Status = "complete"
	} else {
		task.Status = "cannot_decide_winner"
	}

	err = repo.Save(task)
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
}

func main() {
	validation.RegisterRuleConstructor(validation.RuleConstructorMap{
		"DateRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			tArr := i.(map[string]interface{})["limits"].([]interface{})
			limits := [2]int64{}
			if len(tArr) == 1 {
				if t := tArr[0]; t != nil {
					limits[0] = int64(t.(float64))
				}
			} else {
				if t := tArr[0]; t != nil {
					limits[0] = int64(t.(float64))
				}
				if t := tArr[1]; t != nil {
					limits[1] = int64(t.(float64))
				}
			}
			rule := &rules.DateRule{Name: "DateRule", Limits: limits}
			return validation.AppendingRule, rule
		},
		"FollowsRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			idi := i.(map[string]interface{})["id"].(interface{})
			rule := &rules.FollowsRule{Name: "FollowsRule", Id: idi.(string)}
			return validation.SelectRule, rule
		},
		"FollowersRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			username := i.(map[string]interface{})["username"].(interface{})
			amount := i.(map[string]interface{})["amount"].(interface{})
			condition := i.(map[string]interface{})["condition"].(interface{})
			rule := &rules.FollowersRule{Name: "FollowersRule", Amount: int64(amount.(float64)), Username: username.(string), Condition: condition.(string)}
			return validation.PreconditionRule, rule
		},
		"PostLikesRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			shortcode := i.(map[string]interface{})["shortcode"].(interface{})
			amount := i.(map[string]interface{})["amount"].(interface{})
			condition := i.(map[string]interface{})["condition"].(interface{})
			rule := &rules.PostLikesRule{Name: "PostLikesRule", Amount: int64(amount.(float64)), ShortCode: shortcode.(string), Condition: condition.(string)}
			return validation.PreconditionRule, rule
		},
		"PostCommentsRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			shortcode := i.(map[string]interface{})["shortcode"].(interface{})
			amount := i.(map[string]interface{})["amount"].(interface{})
			condition := i.(map[string]interface{})["condition"].(interface{})
			rule := &rules.PostCommentsRule{Name: "PostCommentsRule", Amount: int64(amount.(float64)), ShortCode: shortcode.(string), Condition: condition.(string)}
			return validation.PreconditionRule, rule
		},
		"ParticipantsRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			amount := i.(map[string]interface{})["amount"].(interface{})
			condition := i.(map[string]interface{})["condition"].(interface{})
			rule := &rules.ParticipantsRule{Name: "ParticipantsRule", Amount: int64(amount.(float64)), Condition: condition.(string)}
			return validation.PostconditionRule, rule
		},
		"StoryHasHashTagRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			hashTag := i.(map[string]interface{})["hashtag"].(interface{})
			rule := &rules.StoryHasHashTagRule{Name: "StoryHasHashTagRule", HashTag: hashTag.(string)}
			return validation.SelectRule, rule
		},
		"StoryHasMentionRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			var username string
			var id string
			if u, p := i.(map[string]interface{})["username"].(interface{}); p {
				username = u.(string)
			}
			if u, p := i.(map[string]interface{})["id"].(interface{}); p {
				id = u.(string)
			}
			rule := &rules.StoryHasMentionRule{Name: "StoryHasMentionRule", Username: username, ID: id}
			return validation.SelectRule, rule
		},
		"StoryHasPostRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			shortCode := i.(map[string]interface{})["shortcode"].(interface{})
			rule := &rules.StoryHasPostRule{Name: "StoryHasPostRule", ShortCode: shortCode.(string)}
			return validation.SelectRule, rule
		},
		"StoryHasExternalLinkRule": func(i interface{}) (validation.RuleType, validation.IRule) {
			link := i.(map[string]interface{})["link"].(interface{})
			rule := &rules.StoryHasExternalLinkRule{Name: "StoryHasExternalLinkRule", Link: link.(string)}
			return validation.SelectRule, rule
		},
	})

	solv := solver.GetInstance()
	solv.Run()
	app := gin.Default()
	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			tasksGroup := v1.Group("/tasks")
			{
				comments := tasksGroup.Group("/comments")
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
						req.NumWinners = 1
						err := c.BindJSON(&req)

						if err != nil {
							c.JSON(400, responses.NewValidationErrorJsonResponse())
							return
						}
						task := &tasks.CommentsTask{}
						task.ShortCode = req.ShortCode
						task.Status = "in_progress"
						task.Id = primitive.NewObjectID()
						task.Rules = req.Rules
						task.NumWinners = req.NumWinners
						err = utils.GetNamedTasksRepositoryInstance("CommentTasks").Save(task)

						if err != nil {
							panic(err)
						}

						go execComments(task)
						c.JSON(200, responses.NewSuccessfulCommentsTaskJsonResponse(*task))
					})
				}
				posts := tasksGroup.Group("/posts")
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
						req.NumWinners = 1
						err := c.BindJSON(&req)

						if err != nil {
							c.JSON(400, responses.NewValidationErrorJsonResponse())
							return
						}
						task := &tasks.HashTagTask{}
						task.HashTag = req.HashTag
						task.Status = "in_progress"
						task.Id = primitive.NewObjectID()
						task.Rules = req.Rules
						task.NumWinners = req.NumWinners
						err = utils.GetNamedTasksRepositoryInstance("HashTagTasks").Save(task)
						if err != nil {
							panic(err)
						}

						go execPosts(task)
						c.JSON(200, responses.NewSuccessfulHashTagTaskJsonResponse(*task))
					})
				}
				stories := tasksGroup.Group("/stories")
				{
					stories.GET("/:id", func(c *gin.Context) {
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						res, err := utils.GetNamedTasksRepositoryInstance("HashTagStoryTasks").FindHashTagStoryTaskById(bsonx.ObjectID(id))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, responses.NewSuccessfulHashTagStoryTaskJsonResponse(*res))
					})
					stories.POST("/", func(c *gin.Context) {

						var req requests.HashTagStoryTaskJsonRequest
						req.NumWinners = 1
						err := c.BindJSON(&req)

						if err != nil {
							c.JSON(400, responses.NewValidationErrorJsonResponse())
							return
						}
						task := &tasks.StoriesTask{}
						task.HashTag = req.HashTag
						task.Status = "in_progress"
						task.Id = primitive.NewObjectID()
						task.Rules = req.Rules
						task.NumWinners = req.NumWinners

						acc := repository.GetRepositoryInstance().GetOldestUsed()
						task.Account = acc.Username

						err = utils.GetNamedTasksRepositoryInstance("HashTagStoryTasks").Save(task)
						if err != nil {
							panic(err)
						}

						go execStories(task)
						c.JSON(200, responses.NewSuccessfulHashTagStoryTaskJsonResponse(*task))
					})
				}
			}
		}
	}
	app.Run("0.0.0.0:80")
}
