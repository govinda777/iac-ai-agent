package models

import "time"

// ReviewRequest representa uma requisição de review de PR
type ReviewRequest struct {
	Repository     string `json:"repository"`
	PRNumber       int    `json:"pr_number"`
	Owner          string `json:"owner"`
	InstallationID int64  `json:"installation_id,omitempty"`
}

// ReviewResponse representa o resultado de um review
type ReviewResponse struct {
	ID               string          `json:"id"`
	Repository       string          `json:"repository"`
	PRNumber         int             `json:"pr_number"`
	Score            int             `json:"score"`
	Status           string          `json:"status"` // approved, changes_requested, commented
	Summary          string          `json:"summary"`
	FilesAnalyzed    int             `json:"files_analyzed"`
	TotalSuggestions int             `json:"total_suggestions"`
	Analysis         AnalysisDetails `json:"analysis"`
	FileReviews      []FileReview    `json:"file_reviews"`
	Timestamp        time.Time       `json:"timestamp"`
}

// FileReview representa o review de um arquivo específico
type FileReview struct {
	Filename    string       `json:"filename"`
	Status      string       `json:"status"` // added, modified, deleted
	Additions   int          `json:"additions"`
	Deletions   int          `json:"deletions"`
	Changes     int          `json:"changes"`
	Suggestions []Suggestion `json:"suggestions"`
	Score       int          `json:"score"`
	Comments    []Comment    `json:"comments"`
}

// Comment representa um comentário no código
type Comment struct {
	Path     string `json:"path"`
	Position int    `json:"position,omitempty"`
	Line     int    `json:"line"`
	Body     string `json:"body"`
	Side     string `json:"side,omitempty"` // LEFT, RIGHT
}

// GitHubWebhookPayload representa o payload de um webhook do GitHub
type GitHubWebhookPayload struct {
	Action       string              `json:"action"`
	Number       int                 `json:"number,omitempty"`
	PullRequest  *GitHubPullRequest  `json:"pull_request,omitempty"`
	Repository   *GitHubRepository   `json:"repository"`
	Sender       *GitHubUser         `json:"sender"`
	Installation *GitHubInstallation `json:"installation,omitempty"`
}

// GitHubPullRequest representa um PR do GitHub
type GitHubPullRequest struct {
	ID        int64         `json:"id"`
	Number    int           `json:"number"`
	Title     string        `json:"title"`
	Body      string        `json:"body"`
	State     string        `json:"state"`
	Head      *GitHubCommit `json:"head"`
	Base      *GitHubCommit `json:"base"`
	User      *GitHubUser   `json:"user"`
	HTMLURL   string        `json:"html_url"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// GitHubCommit representa um commit do GitHub
type GitHubCommit struct {
	Label string            `json:"label"`
	Ref   string            `json:"ref"`
	SHA   string            `json:"sha"`
	Repo  *GitHubRepository `json:"repo"`
}

// GitHubRepository representa um repositório do GitHub
type GitHubRepository struct {
	ID            int64       `json:"id"`
	Name          string      `json:"name"`
	FullName      string      `json:"full_name"`
	Owner         *GitHubUser `json:"owner"`
	Private       bool        `json:"private"`
	HTMLURL       string      `json:"html_url"`
	DefaultBranch string      `json:"default_branch"`
}

// GitHubUser representa um usuário do GitHub
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Type      string `json:"type"`
	AvatarURL string `json:"avatar_url"`
}

// GitHubInstallation representa uma instalação do GitHub App
type GitHubInstallation struct {
	ID int64 `json:"id"`
}

// PRScore representa o score calculado para um PR
type PRScore struct {
	Total           int            `json:"total"`
	Security        int            `json:"security"`
	BestPractices   int            `json:"best_practices"`
	Performance     int            `json:"performance"`
	Maintainability int            `json:"maintainability"`
	Documentation   int            `json:"documentation"`
	Breakdown       map[string]int `json:"breakdown"`
}
