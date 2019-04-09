package responses

import (
	"giveaway/data/tasks"
)

type HasStatusJsonResponse struct {
	Status bool `json:"status" bson:"status"`
}

type NotFoundJsonResponse struct {
	HasStatusJsonResponse `json:",inline"`
	Error                 string `json:"error"`
}

func NewNotFoundJsonResponse() NotFoundJsonResponse {
	r := NotFoundJsonResponse{}
	r.Status = false
	r.Error = "err_not_found"
	return r
}

type ValidationErrorJsonResponse struct {
	HasStatusJsonResponse `json:",inline"`
	Error                 string `json:"error"`
}

func NewValidationErrorJsonResponse() ValidationErrorJsonResponse {
	r := ValidationErrorJsonResponse{}
	r.Status = false
	r.Error = "request_validation_error"
	return r
}

type SuccessfulCommentsTaskJsonResponse struct {
	HasStatusJsonResponse `json:",inline"`
	Result                tasks.CommentsTask `json:"result"`
}

func NewSuccessfulCommentsTaskJsonResponse(task tasks.CommentsTask) SuccessfulCommentsTaskJsonResponse {
	r := SuccessfulCommentsTaskJsonResponse{}
	r.Status = true
	r.Result = task
	return r
}

type SuccessfulHashTagTaskJsonResponse struct {
	HasStatusJsonResponse `json:",inline"`
	Result                tasks.HashTagTask `json:"result"`
}

func NewSuccessfulHashTagTaskJsonResponse(task tasks.HashTagTask) SuccessfulHashTagTaskJsonResponse {
	r := SuccessfulHashTagTaskJsonResponse{}
	r.Status = true
	r.Result = task
	return r
}

type SuccessfulHashTagStoryTaskJsonResponse struct {
	HasStatusJsonResponse `json:",inline"`
	Result                tasks.StoriesTask `json:"result"`
}

func NewSuccessfulHashTagStoryTaskJsonResponse(task tasks.StoriesTask) SuccessfulHashTagStoryTaskJsonResponse {
	r := SuccessfulHashTagStoryTaskJsonResponse{}
	r.Status = true
	r.Result = task
	return r
}
