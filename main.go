package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"giveaway/client/validation"
	"giveaway/client/validation/rules"
	"giveaway/client/web"
	"giveaway/data"
	"giveaway/data/errors"
	"giveaway/data/сontainers"
	"giveaway/http/requests"
	"giveaway/http/responses"
	"giveaway/instagram/solver"
	"giveaway/utils"
	"giveaway/utils/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"
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

func filterWinnerHashTag(ret *сontainers.EntryContainer, rules []validation.IRule) (*data.TagMedia, error) {
	_, t, e := filterWinner(ret, rules)
	if t != nil {
		return t.(*data.TagMedia), nil
	}
	return nil, e
}

func filterWinnerComment(ret *сontainers.EntryContainer, rules []validation.IRule) (int, *data.Comment, error) {
	i, t, e := filterWinner(ret, rules)
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

func execPosts(task *data.HashTagTask) {
	cl := web.NewWebClient(&utils.UserAgentGenerator{}, "http://localhost:8888")
	cl.Init()

	var err error = nil
	var ruleRes bool
	repo := utils.GetNamedTasksRepositoryInstance("HashTagTasks")
	dLogger := logger.DefaultLogger()
	ruleRes, err = runRules(task.PreconditionRules, cl, func(rule validation.IRule) error {
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
	}
	if !ruleRes {
		return
	}

	ret := сontainers.NewEntryContainer()
	err = cl.QueryTag(task.HashTag, func(media data.TagMedia) (bool, error) {
		var shouldAdd = true
		for _, rule := range task.AppendRules {
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
		dLogger.Errorf("error: %v", err)
	}

	ruleRes, err = runRules(task.PostconditionRules, ret, func(rule validation.IRule) error {
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

	winner, err := filterWinnerHashTag(ret, task.SelectRules)

	if err != nil {
		dLogger.Errorf("error: %v", err)
	}

	if winner != nil {
		task.Post = winner
		task.Status = "complete"
	} else {
		task.Status = "cannot_decide_winner"
	}
	err = repo.Save(task)
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
}

func execComments(task *data.CommentsTask) {
	cl := web.NewWebClient(&utils.UserAgentGenerator{}, "http://localhost:8888")
	cl.Init()

	ret := сontainers.NewEntryContainer()
	repo := utils.GetNamedTasksRepositoryInstance("CommentTasks")
	dLogger := logger.DefaultLogger()
	var err error = nil
	var ruleRes bool
	ruleRes, err = runRules(task.PreconditionRules, cl, func(rule validation.IRule) error {
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
	}
	if !ruleRes {
		return
	}

	err = cl.QueryComments(task.ShortCode, func(comment data.Comment) (bool, error) {
		var shouldAdd = true
		for _, rule := range task.AppendRules {
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
		dLogger.Errorf("error: %v", err)
	}

	ruleRes, err = runRules(task.PostconditionRules, ret, func(rule validation.IRule) error {
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

	winnerId, winner, err := filterWinnerComment(ret, task.SelectRules)

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

		task.Winner = winner
		task.Above = above
		task.Below = below

		task.Position = winnerId
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
			return validation.AppendRule, rule
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
						task := &data.CommentsTask{}
						task.ShortCode = req.ShortCode
						task.Status = "in_progress"
						task.Id = primitive.NewObjectID()
						task.PreconditionRules = req.Rules.PreconditionRules()
						task.AppendRules = req.Rules.AppendRules()
						task.PostconditionRules = req.Rules.PostconditionRules()
						task.SelectRules = req.Rules.SelectRules()
						err = utils.GetNamedTasksRepositoryInstance("CommentTasks").Save(task)

						if err != nil {
							panic(err)
						}

						go execComments(task)
						c.JSON(200, responses.NewSuccessfulCommentsTaskJsonResponse(*task))
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
						task := &data.HashTagTask{}
						task.HashTag = req.HashTag
						task.Status = "in_progress"
						task.Id = primitive.NewObjectID()
						task.PreconditionRules = req.Rules.PreconditionRules()
						task.AppendRules = req.Rules.AppendRules()
						task.PostconditionRules = req.Rules.PostconditionRules()
						task.SelectRules = req.Rules.SelectRules()
						err = utils.GetNamedTasksRepositoryInstance("HashTagTasks").Save(task)
						if err != nil {
							panic(err)
						}

						go execPosts(task)
						c.JSON(200, responses.NewSuccessfulHashTagTaskJsonResponse(*task))
					})
				}
			}
		}
	}
	app.Run("0.0.0.0:80")
}
