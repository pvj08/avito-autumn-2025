package user

type User struct {
	IsActive bool
	TeamName string
	UserID   string
	Username string
}

type SetIsActiveInput struct {
	IsActive bool
	UserID   string
}

type SetIsActiveOutput struct {
	User
}
