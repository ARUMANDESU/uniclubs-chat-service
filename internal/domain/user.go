package domain

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
}
