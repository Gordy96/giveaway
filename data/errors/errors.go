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

type BeforeMinimumDate struct {
}

func (e BeforeMinimumDate) Error() string {
	return "examined date is before allowed"
}

type AfterMaximumDate struct {
}

func (e AfterMaximumDate) Error() string {
	return "examined date is after allowed"
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

type EndOfListError struct {
}

func (e EndOfListError) Error() string {
	return "reached end of a list"
}
