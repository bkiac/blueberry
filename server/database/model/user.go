package model

import (
	"context"

	"github.com/bkiac/blueberry/server/database"
	"github.com/bkiac/blueberry/server/database/ent"
)

// User holds the query and mutation definitions for the User model.
type User struct{}

// RegisterArgs holds the arguments needed to create a User.
type RegisterArgs struct {
	Username string
	Password string
}

// Create creates a user in the database.
func (m *User) Create(c context.Context, args RegisterArgs) (*ent.User, error) {
	user, err := database.Client.User.Create().SetUsername(args.Username).SetPassword(args.Password).Save(c)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Get returns a user from the database, if it exists.
func (m *User) Get(c context.Context, id int) (*ent.User, error) {
	user, err := database.Client.User.Get(c, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAll returns all of the users from the database.
func (m *User) GetAll(c context.Context) ([]*ent.User, error) {
	users, err := database.Client.User.Query().All(c)
	if err != nil {
		return nil, err
	}
	return users, nil
}
