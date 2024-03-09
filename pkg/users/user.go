package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type User struct {
	Id       *primitive.ObjectID `json:"id,omitempty"  bson:"_id,omitempty"`
	Username string              `json:"username,omitempty"`
	Password string              `json:"password,omitempty"`
}

func (u *User) Validate() string {

	if u.Username == "" {
		return "missing username"
	}

	if u.Password == "" {
		return "missing password"
	}

	return ""
}

func (u *User) Prepare() {
	u.Username = strings.TrimSpace(u.Username)
}
