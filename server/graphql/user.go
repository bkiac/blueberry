package graphql

import (
	"context"
	"strconv"

	"github.com/bkiac/blueberry/server/database"
	"github.com/graph-gophers/graphql-go"
)

type IDArgs struct {
	ID graphql.ID
}

type rootResolver struct{}

func (r *rootResolver) Users(c context.Context) ([]*userResolver, error) {
	userEntities, e := database.GetUsers(c)
	if e != nil {
		return nil, e
	}

	var userResolvers []*userResolver
	for _, userEntity := range userEntities {
		userResolvers = append(userResolvers, &userResolver{&UserGraphQL{
			ID:       userEntity.ID,
			Username: userEntity.Username,
		}})
	}

	return userResolvers, nil
}

func (r *rootResolver) User(c context.Context, args IDArgs) (*userResolver, error) {
	id64, e := strconv.ParseInt(string(args.ID), 0, 64)
	id := int(id64)
	if e != nil {
		return nil, e
	}

	userEntity, e := database.GetUser(c, id)
	if e != nil {
		return nil, e
	}

	userResolver := &userResolver{&UserGraphQL{
		ID:       userEntity.ID,
		Username: userEntity.Username,
	}}

	return userResolver, nil
}

type RegisterArgs struct {
	Username string
	Password string
}

func (r *rootResolver) Register(c context.Context, args RegisterArgs) (*userResolver, error) {
	userEntity, e := database.CreateUser(c, args)
	if e != nil {
		return nil, e
	}

	userResolver := &userResolver{&UserGraphQL{
		ID:       userEntity.ID,
		Username: userEntity.Username,
	}}

	return userResolver, nil
}

type UserGraphQL struct {
	ID       int
	Username string
}

type userResolver struct{ user *UserGraphQL }

func (r *userResolver) ID() graphql.ID {
	return graphql.ID(strconv.Itoa(r.user.ID))
}

func (r *userResolver) Username() string {
	return r.user.Username
}
