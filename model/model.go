package model

type User struct {
	ID       string `json:"user_id"`
	Nickname string `json:"nickname"`
	Comment  string `json:"comment,omitempty"`
}

func (u *User) FillNickname() {
	if u.Nickname == "" {
		u.Nickname = u.ID
	}
}

type CreateUserRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}
