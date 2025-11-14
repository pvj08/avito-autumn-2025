package domain

type UserID string
type TeamID string
type PullRequestID string

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type User struct {
	ID       UserID `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`

	// один пользователь — одна команда
	TeamID TeamID `json:"team_id"`
}

type Team struct {
	ID   TeamID `json:"id"`
	Name string `json:"name"`

	MemberIDs []UserID `json:"member_ids"`
}

type PullRequest struct {
	ID    PullRequestID `json:"id"`
	Title string        `json:"title"`

	AuthorID UserID   `json:"author_id"`
	Status   PRStatus `json:"status"`

	// Список назначенных ревьюверов (0..2).
	// Инварианты:
	//  - не более двух элементов;
	//  - не должен содержать автора;
	//  - изменяется только когда Status == OPEN.
	ReviewerIDs []UserID `json:"reviewer_ids"`
}
