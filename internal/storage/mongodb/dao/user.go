package dao

import "github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"

type User struct {
	ID        int64  `json:"id" bson:"_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	AvatarURL string `json:"avatar_url" bson:"avatar_url"`
}

func (u *User) ToDomain() domain.User {
	if u == nil {
		return domain.User{}
	}

	return domain.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		AvatarURL: u.AvatarURL,
	}
}

func UserFromDomain(d domain.User) User {
	return User{
		ID:        d.ID,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		AvatarURL: d.AvatarURL,
	}
}
