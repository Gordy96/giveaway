package main

import (
	"github.com/gin-gonic/gin"
	"giveaway/client/validation"
	"giveaway/client/validation/rules"
	"giveaway/client/web"
	"giveaway/data/tasks"
	"giveaway/http/proxies"
	"giveaway/http/requests"
	"giveaway/http/responses"
	"giveaway/instagram/account/repository"
	"giveaway/instagram/solver"
	"giveaway/queue"
	"giveaway/utils"
	repository2 "giveaway/utils/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func main() {
	//Run queue instance
	queue.GetGlobalInstance(10).Run()
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
			rule := &rules.FollowsRule{Name: "FollowsRule"}
			if x, ok := i.(map[string]interface{})["id"]; ok {
				rule.Id = x.(interface{}).(string)
			}
			if x, ok := i.(map[string]interface{})["username"]; rule.Id == "" && ok {
				rule.Username = x.(interface{}).(string)
				cl := web.NewWebClient(&utils.UserAgentGenerator{}, proxies.GetGlobalInstance().GetNext())
				cl.Init()
				r, e := cl.GetUserInfo(rule.Username)
				if e != nil {
					panic(e)
				}
				rule.Id = r.Id
			}
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
	//Run solver instance
	solver.GetInstance().Run()
	app := gin.Default()
	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			util := v1.Group("/util/info")
			{
				util.GET("/post/:code", func(c *gin.Context) {
					code := c.Param("code")
					cl := web.NewWebClient(&utils.UserAgentGenerator{}, proxies.GetGlobalInstance().GetNext())
					r, err := cl.GetPostSummary(code)
					if err != nil {
						c.JSON(404, responses.NewNotFoundJsonResponse())
						return
					}
					c.JSON(200, responses.NewSuccessfulPostInfoResponse(*r))
				})
				util.GET("/tag/:tag", func(c *gin.Context) {
					tag := c.Param("tag")
					cl := web.NewWebClient(&utils.UserAgentGenerator{}, proxies.GetGlobalInstance().GetNext())
					r, err := cl.GetHashTagSummary(tag)
					if err != nil {
						c.JSON(404, responses.NewNotFoundJsonResponse())
						return
					}
					c.JSON(200, responses.NewSuccessfulTagInfoResponse(*r))
				})
			}
			tasksGroup := v1.Group("/tasks")
			{
				comments := tasksGroup.Group("/comments")
				{
					comments.GET("/:id/decide", func(c *gin.Context) {
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}
						var task tasks.CommentsTask
						err = repository2.GetNamedRepositoryInstance("CommentTasks").FindTaskById(bsonx.ObjectID(id), &task)

						queue.GetGlobalInstance().Enqueue(&tasks.DecideCommentsWinnerCommand{Task: task})
						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, responses.NewSuccessfulCommentsTaskJsonResponse(task))
					})
					comments.GET("/:id", func(c *gin.Context) {
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}
						var task tasks.CommentsTask
						err = repository2.GetNamedRepositoryInstance("CommentTasks").FindTaskById(bsonx.ObjectID(id), &task)
						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, responses.NewSuccessfulCommentsTaskJsonResponse(task))
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
						task.Status = tasks.InProgress
						task.Id = primitive.NewObjectID()
						task.Rules = req.Rules
						task.NumWinners = req.NumWinners

						err = repository2.GetNamedRepositoryInstance("CommentTasks").Save(task)

						if err != nil {
							panic(err)
						}
						queue.GetGlobalInstance().Enqueue(&tasks.FetchCommentDataCommand{Task: *task})
						c.JSON(200, responses.NewSuccessfulCommentsTaskJsonResponse(*task))
					})
				}
				posts := tasksGroup.Group("/posts")
				{
					posts.GET("/:id/decide", func(c *gin.Context) {
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						var task tasks.HashTagTask
						err = repository2.GetNamedRepositoryInstance("HashTagTasks").FindTaskById(bsonx.ObjectID(id), &task)
						queue.GetGlobalInstance().Enqueue(&tasks.DecideHashTagWinnerCommand{Task: task})
						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, responses.NewSuccessfulHashTagTaskJsonResponse(task))
					})
					posts.GET("/:id", func(c *gin.Context) {
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}
						var task tasks.HashTagTask
						err = repository2.GetNamedRepositoryInstance("HashTagTasks").FindTaskById(bsonx.ObjectID(id), &task)
						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, responses.NewSuccessfulHashTagTaskJsonResponse(task))
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
						err = repository2.GetNamedRepositoryInstance("HashTagTasks").Save(task)
						if err != nil {
							panic(err)
						}

						queue.GetGlobalInstance().Enqueue(&tasks.FetchHashTagDataCommand{Task: *task})
						c.JSON(200, responses.NewSuccessfulHashTagTaskJsonResponse(*task))
					})
				}
				stories := tasksGroup.Group("/stories")
				{
					stories.GET("/:id/decide", func(c *gin.Context) {
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}
						var task tasks.StoriesTask
						err = repository2.GetNamedRepositoryInstance("HashTagStoryTasks").FindTaskById(bsonx.ObjectID(id), task)
						queue.GetGlobalInstance().Enqueue(&tasks.DecideStoriesWinnerCommand{Task: task})
						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, responses.NewSuccessfulHashTagStoryTaskJsonResponse(task))
					})
					stories.GET("/:id", func(c *gin.Context) {
						id, err := primitive.ObjectIDFromHex(c.Param("id"))

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}
						var task tasks.StoriesTask
						err = repository2.GetNamedRepositoryInstance("HashTagStoryTasks").FindTaskById(bsonx.ObjectID(id), &task)

						if err != nil {
							c.JSON(404, responses.NewNotFoundJsonResponse())
							return
						}

						c.JSON(200, responses.NewSuccessfulHashTagStoryTaskJsonResponse(task))
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

						err = repository2.GetNamedRepositoryInstance("HashTagStoryTasks").Save(task)
						if err != nil {
							panic(err)
						}

						queue.GetGlobalInstance().Enqueue(&tasks.FetchStoriesDataCommand{Task: *task})
						c.JSON(200, responses.NewSuccessfulHashTagStoryTaskJsonResponse(*task))
					})
				}
			}
		}
	}
	app.Run("0.0.0.0:8080")
}
