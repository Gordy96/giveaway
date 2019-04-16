package tasks

import (
	"fmt"
	"giveaway/client/api"
	"giveaway/client/validation"
	"giveaway/client/web"
	"giveaway/data"
	"giveaway/data/errors"
	"giveaway/data/сontainers"
	"giveaway/http/proxies"
	"giveaway/instagram/account/repository"
	"giveaway/instagram/structures"
	"giveaway/instagram/structures/stories"
	"giveaway/utils"
	"giveaway/utils/logger"
	dbRepo "giveaway/utils/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

type TaskStatus string

const (
	New        TaskStatus = "new"
	InProgress TaskStatus = "in_progress"
	Complete   TaskStatus = "complete"
	Cancelled  TaskStatus = "cancelled"
	Failed     TaskStatus = "failed"
	NoWinner   TaskStatus = "no_winner"
)

func runRules(rules []validation.IRule, arg interface{}, failCallback func(validation.IRule, error) (bool, error)) (bool, error) {
	for _, rule := range rules {
		ruleRes, err := rule.Validate(arg)
		if !ruleRes {
			return failCallback(rule, err)
		}
	}
	return true, nil
}

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
func filterWinnerHashTagStories(ret *сontainers.EntryContainer, rules []validation.IRule, excluded []int) (int, *stories.StoryItem, error) {
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
		return i, t.(*stories.StoryItem), nil
	}
	return -1, nil, e
}

type BaseTaskModel struct {
	Id         primitive.ObjectID        `json:"_id,omitempty" bson:"_id"`
	SourceUrl  string                    `json:"source_url" bson:"source_url"`
	CreatedAt  int64                     `json:"created_at" bson:"created_at"`
	FinishedAt int64                     `json:"finished_at" bson:"finished_at"`
	Status     TaskStatus                `json:"status" bson:"status"`
	Comment    string                    `json:"comment" bson:"comment"`
	Rules      validation.RuleCollection `json:"rules" bson:"rules"`
	NumWinners int                       `json:"num_winners"`
}

type CommentsTask struct {
	BaseTaskModel `json:",inline" bson:",inline"`
	ShortCode     string               `json:"shortcode" bson:"shortcode"`
	Winners       []data.CommentWinner `json:"winners" bson:"winners"`
}

func (c *CommentsTask) GetKey() interface{} {
	return c.Id
}

func (c *CommentsTask) FetchData() {
	cl := web.NewWebClient(&utils.UserAgentGenerator{}, proxies.GetGlobalInstance().GetNext())

	cl.Init()

	repo := dbRepo.GetNamedRepositoryInstance("CommentTasks")
	dLogger := logger.DefaultLogger()
	var err error = nil
	var ruleRes bool
	ruleRes, err = runRules(c.Rules.PreconditionRules, cl, func(rule validation.IRule, e error) (bool, error) {
		c.Status = Cancelled
		c.Comment = fmt.Sprintf("failed on precondition rule: %v", rule)
		err = repo.Save(c)
		if err != nil {
			return false, err
		}
		return false, nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
		c.Status = Failed
		c.Comment = fmt.Sprintf("got error: %v", err)
		err = repo.Save(c)
	}
	summary, err := cl.GetPostSummary(c.ShortCode)
	c.SourceUrl = summary.DisplayURL
	if err != nil {
		dLogger.Errorf("error: %v", err)
		c.Status = Failed
		c.Comment = fmt.Sprintf("got error: %v", err)
		err = repo.Save(c)
	}
	if !ruleRes {
		return
	}

	store := dbRepo.GetNamedRepositoryInstance("Comments")

	err = cl.QueryComments(c.ShortCode, func(comment data.Comment) (bool, error) {
		temp := data.CommentContainer{
			TaskId:  c.Id,
			Id:      primitive.NewObjectID(),
			Comment: comment,
		}
		var shouldAdd bool
		shouldAdd, err = runRules(c.Rules.AppendingRules, &comment, func(rule validation.IRule, e error) (bool, error) {
			if e != nil {
				switch /*e := */ e.(type) {
				case errors.BeforeMinimumDate:
					return false, e
				case errors.AfterMaximumDate:
					return false, nil
				}
			}
			return false, e
		})
		if shouldAdd {
			store.Save(temp)
		}
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		dLogger.Errorf("error: %v", err)
	} else {
		c.Status = Complete
		repo.Save(c)
		c.DecideWinner()
	}
}

//Will not use it
func (c *CommentsTask) DropData() {
	store := dbRepo.GetNamedRepositoryInstance("Comments")
	store.DeleteMany(bson.M{
		"task_id": c.Id,
	})
}

func (c *CommentsTask) DecideWinner() {
	store := dbRepo.GetNamedRepositoryInstance("Comments")
	cursor, _ := store.FindAll(bson.M{
		"task_id": c.Id,
	})
	ret := сontainers.NewEntryContainer()
	var err error = nil
	var container data.CommentContainer
	for cursor.Next(nil) {
		err = cursor.Decode(&container)
		if err != nil {
			panic(err)
		}
		var shouldAdd bool
		var comment = container.Comment
		shouldAdd, err = runRules(c.Rules.AppendingRules, &comment, func(rule validation.IRule, e error) (bool, error) {
			if e != nil {
				switch /*e := */ e.(type) {
				case errors.BeforeMinimumDate:
					return false, e
				case errors.AfterMaximumDate:
					return false, nil
				}
			}
			return false, e
		})
		if shouldAdd {
			ret.Add(&comment)
		}
		if err != nil {
			break
		}
	}
	repo := dbRepo.GetNamedRepositoryInstance("CommentTasks")
	dLogger := logger.DefaultLogger()
	ruleRes, err := runRules(c.Rules.PostconditionRules, ret, func(rule validation.IRule, e error) (bool, error) {
		c.Status = Failed
		c.Comment = fmt.Sprintf("failed on postcondition rule: %v", rule)
		err = repo.Save(c)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
	if !ruleRes {
		return
	}

	c.Winners = make([]data.CommentWinner, 0)
	if ret.LengthNoDuplicates() >= c.NumWinners {
		taken := make([]int, 0)
		for i := 0; i < c.NumWinners; i++ {
			winnerId, winner, err := filterWinnerComment(ret, c.Rules.SelectRules, taken)
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
				c.Winners = append(c.Winners, wSet)
			} else {
				break
			}
		}
	}

	if len(c.Winners) == c.NumWinners {
		c.Status = Complete
	} else {
		c.Status = NoWinner
	}

	err = repo.Save(c)
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
}

type HashTagTask struct {
	BaseTaskModel `json:",inline" bson:",inline"`
	HashTag       string          `json:"hashtag" bson:"hashtag"`
	Winners       []data.TagMedia `json:"winners" bson:"winners"`
}

func (h *HashTagTask) GetKey() interface{} {
	return h.Id
}

func (h *HashTagTask) FetchData() {
	cl := web.NewWebClient(&utils.UserAgentGenerator{}, proxies.GetGlobalInstance().GetNext())
	cl.Init()

	repo := dbRepo.GetNamedRepositoryInstance("HashTagTasks")
	dLogger := logger.DefaultLogger()

	var err error = nil
	var ruleRes bool
	ruleRes, err = runRules(h.Rules.PreconditionRules, cl, func(rule validation.IRule, e error) (bool, error) {
		h.Status = Cancelled
		h.Comment = fmt.Sprintf("failed on precondition rule: %v", rule)
		err = repo.Save(h)
		if err != nil {
			return false, err
		}
		return false, nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
		h.Status = Failed
		h.Comment = fmt.Sprintf("got error: %v", err)
		err = repo.Save(h)
	}
	if !ruleRes {
		return
	}

	store := dbRepo.GetNamedRepositoryInstance("Posts")
	var dupes map[string]bool = make(map[string]bool)
	err = cl.QueryTag(h.HashTag, func(media data.TagMedia, counter *int) (bool, error) {
		temp := data.PostContainer{
			TaskId: h.Id,
			Id:     primitive.NewObjectID(),
			Post:   media,
		}
		var shouldAdd bool
		shouldAdd, err = runRules(h.Rules.AppendingRules, &media, func(rule validation.IRule, e error) (bool, error) {
			if e != nil {
				switch /*e := */ e.(type) {
				case errors.BeforeMinimumDate:
					return false, nil
				case errors.AfterMaximumDate:
					if _, p := dupes[media.Id]; !p {
						dupes[media.Id] = true
						*counter++
					}
					return false, nil
				}
			}
			return false, e
		})
		if shouldAdd {
			if _, p := dupes[media.Id]; !p {
				store.Save(temp)
				dupes[media.Id] = true
				*counter++
			}
		}
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		dLogger.Errorf("error: %v", err)
	} else {
		h.Status = Complete
		repo.Save(h)
		h.DecideWinner()
	}
}

//Will not use it
func (h *HashTagTask) DropData() {
	store := dbRepo.GetNamedRepositoryInstance("Posts")
	store.DeleteMany(bson.M{
		"task_id": h.Id,
	})
}

func (h *HashTagTask) DecideWinner() {
	store := dbRepo.GetNamedRepositoryInstance("Posts")
	cursor, _ := store.FindAll(bson.M{
		"task_id": h.Id,
	})
	ret := сontainers.NewEntryContainer()
	var err error = nil
	var container data.PostContainer
	for cursor.Next(nil) {
		err = cursor.Decode(&container)
		if err != nil {
			panic(err)
		}
		var shouldAdd bool
		var post = container.Post
		shouldAdd, err = runRules(h.Rules.AppendingRules, &post, func(rule validation.IRule, e error) (bool, error) {
			if e != nil {
				switch /*e := */ e.(type) {
				case errors.BeforeMinimumDate:
					return false, e
				case errors.AfterMaximumDate:
					return false, nil
				}
			}
			return false, e
		})
		if shouldAdd {
			ret.Add(&post)
		}
		if err != nil {
			break
		}
	}
	repo := dbRepo.GetNamedRepositoryInstance("HashTagTasks")
	dLogger := logger.DefaultLogger()
	ruleRes, err := runRules(h.Rules.PostconditionRules, ret, func(rule validation.IRule, e error) (bool, error) {
		h.Status = Failed
		h.Comment = fmt.Sprintf("failed on postcondition rule: %v", rule)
		err = repo.Save(h)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
	if !ruleRes {
		return
	}

	h.Winners = make([]data.TagMedia, 0)
	if ret.LengthNoDuplicates() >= h.NumWinners {
		taken := make([]int, 0)
		for i := 0; i < h.NumWinners; i++ {

			id, winner, err := filterWinnerHashTag(ret, h.Rules.SelectRules, taken)
			taken = append(taken, id)
			if err != nil {
				dLogger.Errorf("error: %v", err)
			}
			if winner != nil {
				h.Winners = append(h.Winners, *winner)
			} else {
				break
			}
		}
	}

	if len(h.Winners) == h.NumWinners {
		h.Status = Complete
	} else {
		h.Status = NoWinner
	}

	err = repo.Save(h)
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
}

type StoryWinner struct {
	Story stories.StoryItem `json:"story" bson:"story"`
	Data  []byte            `json:"data" bson:"data"`
}

type StoriesTask struct {
	BaseTaskModel `json:",inline" bson:",inline"`
	HashTag       string        `json:"hashtag" bson:"hashtag"`
	Winners       []StoryWinner `json:"winners" bson:"winners"`
	Account       string        `json:"account" bson:"account"`
}

func (s *StoriesTask) GetKey() interface{} {
	return s.Id
}

func (s *StoriesTask) FetchData() {
	var err error = nil
	var ruleRes bool

	acc := repository.GetRepositoryInstance().FindByUsername(s.Account)
	cl := api.NewApiClientWithAccount(acc)
	cl.Login()
	repo := dbRepo.GetNamedRepositoryInstance("HashTagStoryTasks")
	dLogger := logger.DefaultLogger()
	clw := web.NewWebClient(&utils.UserAgentGenerator{}, proxies.GetGlobalInstance().GetNext())
	ruleRes, err = runRules(s.Rules.PreconditionRules, clw, func(rule validation.IRule, e error) (bool, error) {
		s.Status = Cancelled
		s.Comment = fmt.Sprintf("failed on precondition rule: %v", rule)
		err = repo.Save(s)
		if err != nil {
			return false, err
		}
		return false, nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
		s.Status = Failed
		s.Comment = fmt.Sprintf("got error: %v", err)
		err = repo.Save(s)
	}
	if !ruleRes {
		return
	}

	var story *structures.Story = nil
	store := dbRepo.GetNamedRepositoryInstance("Stories")
	for {
		watched := make([]structures.WatchedStoryEntry, 0)
		var upserted int64 = 0
		story, err = cl.QueryHashTagStories(s.HashTag, func(item stories.StoryItem) (bool, error) {
			start := time.Now()
			time.Sleep(time.Duration(math.Min(5, item.VideoDuration) * float64(time.Second)))
			watched = append(watched, structures.WatchedStoryEntry{
				Item: item,
				Duration: [2]int64{
					start.Unix(),
					time.Now().Unix(),
				},
			})
			uri := item.ImageVersions2.Candidates[0].URL
			r, err := http.Get(uri)
			if err != nil {
				dLogger.Errorf("something bad happened while trying to get %s: %s", uri, err)
				return false, err
			}
			content, _ := ioutil.ReadAll(r.Body)
			entry := data.StoryContainer{
				Id:        item.Pk,
				CreatedAt: time.Now().UnixNano(),
				TaskId:    s.Id,
				Data:      content,
				Story:     item,
			}
			ins, err := store.Insert(entry)
			if err != nil {
				return false, err
			}
			upserted += ins
			return true, nil
		})
		if len(watched) > 0 {
			cl.MarkHashTagStoriesAsSeen(story, watched...)
		}
		if err != nil {
			switch /*e := */ err.(type) {
			case errors.AfterMaximumDate:
				break
			default:
				dLogger.Errorf("something happened: %s", err)
			}
		}
		if story == nil || upserted == 0 {
			break
		}
	}
	s.Status = Complete
	repo.Save(s)
	s.DecideWinner()
}

//Will not use it
func (s *StoriesTask) DropData() {}

func (s *StoriesTask) DecideWinner() {
	store := dbRepo.GetNamedRepositoryInstance("Stories")
	cursor, _ := store.FindAll(bson.M{
		"task_id": s.Id,
	})
	ret := сontainers.NewEntryContainer()
	var err error = nil
	var container data.StoryContainer
	for cursor.Next(nil) {
		err = cursor.Decode(&container)
		if err != nil {
			panic(err)
		}
		var shouldAdd bool
		var stry = container.Story
		shouldAdd, err = runRules(s.Rules.AppendingRules, &stry, func(rule validation.IRule, e error) (bool, error) {
			if e != nil {
				switch /*e := */ e.(type) {
				case errors.BeforeMinimumDate:
					return false, nil
				case errors.AfterMaximumDate:
					return false, e
				}
			}
			return false, e
		})
		if err != nil {
			break
		}
		if shouldAdd {
			ret.Add(&stry)
		}
	}
	repo := dbRepo.GetNamedRepositoryInstance("HashTagStoryTasks")
	dLogger := logger.DefaultLogger()
	ruleRes, err := runRules(s.Rules.PostconditionRules, ret, func(rule validation.IRule, e error) (bool, error) {
		s.Status = Failed
		s.Comment = fmt.Sprintf("failed on postcondition rule: %v", rule)
		err = repo.Save(s)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
	if !ruleRes {
		return
	}

	s.Winners = make([]StoryWinner, 0)
	taken := make([]int, 0)

	for i := 0; i < s.NumWinners; i++ {
		id, winner, err := filterWinnerHashTagStories(ret, s.Rules.SelectRules, taken)
		taken = append(taken, id)
		if err != nil {
			dLogger.Errorf("error: %v", err)
		}
		if winner != nil {
			sc, _ := store.FindStory(winner.Pk)
			s.Winners = append(s.Winners, StoryWinner{
				Story: *winner,
				Data:  sc.Data,
			})
		} else {
			break
		}
	}
	if len(s.Winners) == s.NumWinners {
		s.Status = Complete
	} else {
		s.Status = NoWinner
	}

	err = repo.Save(s)
	if err != nil {
		dLogger.Errorf("error: %v", err)
	}
	acrep := repository.GetRepositoryInstance()
	acc := acrep.FindByUsername(s.Account)
	acrep.Release(acc)
}
