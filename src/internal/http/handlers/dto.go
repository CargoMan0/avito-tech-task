package handlers

type teamDTO struct {
	Name    string          `json:"team_name"`
	Members []teamMemberDTO `json:"members"`
}

type teamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type pullRequestDTO struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	Status            string   `json:"status"`
	AuthorID          string   `json:"author_id"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	MergedAt          string   `json:"merged_at,omitempty"`
	CreatedAt         string   `json:"created_at,omitempty"`
}

type userDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
type pullRequestShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}
