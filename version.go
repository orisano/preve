package preve

type Version struct {
	EventID       string `json:"event_id" validate:"required"`
	PullRequestID string `json:"pull_request_id" validate:"required"`
}

func IsOlderVersion(current, version *Version) bool {
	if len(current.EventID) != len(version.EventID) {
		return len(current.EventID) > len(version.EventID)
	}
	return current.EventID > version.EventID
}
