package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bkiac/blueberry/server/ent"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var db *ent.Client
var schema *graphql.Schema

// CreateUser creates a user in the database.
func CreateUser(c context.Context, args RegisterArgs) (*ent.User, error) {
	u, e := db.User.Create().SetUsername(args.Username).SetPassword(args.Password).Save(c)
	if e != nil {
		return nil, e
	}
	return u, nil
}

// GetUsers returns all of the users from the database.
func GetUsers(c context.Context) ([]*ent.User, error) {
	us, e := db.User.Query().All(c)
	if e != nil {
		return nil, e
	}
	return us, nil
}

// GetUser returns a user from the database, if it exists.
func GetUser(c context.Context, id int) (*ent.User, error) {
	u, e := db.User.Get(c, id)
	if e != nil {
		return nil, e
	}
	return u, nil
}

type IDArgs struct {
	ID string
}

type RootResolver struct{}

func (r *RootResolver) Users(c context.Context) ([]*UserResolver, error) {
	userEntities, e := GetUsers(c)
	if e != nil {
		return nil, e
	}

	var userResolvers []*UserResolver
	for _, userEntity := range userEntities {
		userResolvers = append(userResolvers, &UserResolver{&UserGraphQL{
			ID:       userEntity.ID,
			Username: userEntity.Username,
		}})
	}

	return userResolvers, nil
}

func (r *RootResolver) User(c context.Context, args IDArgs) (*UserResolver, error) {
	id64, e := strconv.ParseInt(string(args.ID), 0, 64)
	id := int(id64)
	if e != nil {
		return nil, e
	}

	userEntity, e := GetUser(c, id)
	if e != nil {
		return nil, e
	}

	userResolver := &UserResolver{&UserGraphQL{
		ID:       userEntity.ID,
		Username: userEntity.Username,
	}}

	return userResolver, nil
}

type RegisterArgs struct {
	Username string
	Password string
}

func (r *RootResolver) Register(c context.Context, args RegisterArgs) (*UserResolver, error) {
	userEntity, e := CreateUser(c, args)
	if e != nil {
		return nil, e
	}

	userResolver := &UserResolver{&UserGraphQL{
		ID:       userEntity.ID,
		Username: userEntity.Username,
	}}

	return userResolver, nil
}

type UserGraphQL struct {
	ID       int
	Username string
}

type UserResolver struct{ user *UserGraphQL }

func (r *UserResolver) ID() graphql.ID {
	return graphql.ID(strconv.Itoa(r.user.ID))
}

func (r *UserResolver) Username() string {
	return r.user.Username
}

func main() {
	ctx := context.Background()
	var e error

	// Sets up database connection.
	dsn := os.Getenv("DATABASE")
	db, e = ent.Open("postgres", dsn)
	if e != nil {
		log.Fatalf("Failed to open database: %v", e)
		panic(e)
	}
	log.Println("💽 Database opened successfully!")
	defer db.Close()
	if e := db.Schema.Create(ctx); e != nil {
		log.Fatalf("Failed to create schema resources: %v", e)
	}

	// Reads and parses the GraphQL schema.
	bs, e := ioutil.ReadFile("./schema.graphql")
	if e != nil {
		log.Fatalf("Failed to read GraphQL schema: %v", e)
	}
	ss := string(bs)
	schema, e = graphql.ParseSchema(ss, &RootResolver{})
	if e != nil {
		log.Fatalf("Failed to parse GraphQL schema: %v", e)
	}

	// Starts GraphQL server.
	http.Handle("/", &relay.Handler{Schema: schema})
	log.Println("📡 GraphQL server is listening!")
	e = http.ListenAndServe(":8080", nil)
	if e != nil {
		log.Fatalf("Failed to start GraphQL server: %v", e)
	}
}
