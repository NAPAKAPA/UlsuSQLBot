package models

import (
	"encoding/json"
	"github.com/education-bot/pkg/core"
	"github.com/education-bot/pkg/database"
)

const (
	UsersTable = "users"
)

type UserRequest struct {
	Id   int64      `json:"userId"`
	User *core.User `json:"user"`
}

// Add by id
func (user *UserRequest) Add() (item *core.User, err error) {
	db, err := database.OpenDb(UsersTable)
	if err != nil {
		return
	}
	item = user.User
	err = db.BitAdd(item.Key(), item.Value())
	return
}

// Get by id
func (user *UserRequest) Get(id []byte) (item *core.User, err error) {
	db, err := database.OpenDb(UsersTable)
	if err != nil {
		return
	}
	err, bytes := db.BitGet(id)
	if err != nil {
		return
	}
	if bytes == nil {
		return
	}
	err = json.Unmarshal(bytes, &item)
	if err != nil {
		return
	}
	return
}

// Update by id
func (user *UserRequest) Update(id []byte) (item *core.User, err error) {
	db, err := database.OpenDb(UsersTable)
	if err != nil {
		return
	}
	item = user.User
	err = db.BitAdd(id, item.Value())
	return
}

// Delete by id
func (user *UserRequest) Delete(id []byte) (err error) {
	db, err := database.OpenDb(UsersTable)
	if err != nil {
		return
	}
	return db.Delete(id)
}
