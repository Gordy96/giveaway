package responses

import "giveaway/data"

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
	Result                data.CommentsTask `json:"result"`
}

func NewSuccessfulCommentsTaskJsonResponse(task data.CommentsTask) SuccessfulCommentsTaskJsonResponse {
	r := SuccessfulCommentsTaskJsonResponse{}
	r.Status = true
	r.Result = task
	return r
}

type SuccessfulHashTagTaskJsonResponse struct {
	HasStatusJsonResponse `json:",inline"`
	Result                data.HashTagTask `json:"result"`
}

func NewSuccessfulHashTagTaskJsonResponse(task data.HashTagTask) SuccessfulHashTagTaskJsonResponse {
	r := SuccessfulHashTagTaskJsonResponse{}
	r.Status = true
	r.Result = task
	return r
}
