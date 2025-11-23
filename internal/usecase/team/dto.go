package team

type Team struct {
	Members  []TeamMember
	TeamName string
}

type TeamMember struct {
	IsActive bool
	UserID   string
	Username string
}

type AddInput struct {
	Team
}

type AddOutput struct {
	Team
}

type GetInput struct {
	TeamName string
}

type GetOutput struct {
	Team
}
