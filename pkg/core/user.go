package core

import (
	"fmt"
	"github.com/education-bot/pkg/utils"
)

type User struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Group  string `json:"group"`
	Course string `json:"course"`
}

// Key get key bytes
func (user *User) Key() []byte {
	return []byte(fmt.Sprintf("%d", user.Id))
}

// Value get value bytes
func (user *User) Value() []byte {
	return utils.ToJsonBytes(user)
}
