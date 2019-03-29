package errors

type HttpForbidden struct {
}

func (e HttpForbidden) Error() string {
	return "forbidden"
}

type HttpTooManyRequests struct {
}

func (e HttpTooManyRequests) Error() string {
	return "too many requests"
}

type ShouldStopIterationError struct {
}

func (e ShouldStopIterationError) Error() string {
	return "should break loop"
}

type UnknownRuleError struct {
}

func (e UnknownRuleError) Error() string {
	return "unknown rule occasion"
}

type LoginRequired struct {
}

func (e LoginRequired) Error() string {
	return "login required"
}

type ValidationCriticalFailure struct {
}

func (e ValidationCriticalFailure) Error() string {
	return "critical error while processing"
}
