package entity

type ErrorResponse struct {
	ErrorCode int       `json:"error_code"`
	Message   ErrorCode `json:"message"`
}

type ErrorCode string

const (
	TeamExists     ErrorCode = "TEAM_EXISTS"
	PRExists       ErrorCode = "PR_EXISTS"
	PRMerged       ErrorCode = "PR_MERGED"
	NotAssigned    ErrorCode = "NOT_ASSIGNED"
	NoCandidate    ErrorCode = "NO_CANDIDATE"
	NotFound       ErrorCode = "NOT_FOUND"
	InvalidRequest ErrorCode = "INVALID_REQUEST"
	InternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
)
