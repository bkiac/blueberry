package graphql

import (
	"io/ioutil"
	"log"

	"github.com/graph-gophers/graphql-go"
)

var Schema *graphql.Schema

func SetupSchema() {
	// Reads and parses the GraphQL schema.
	bytestr, err := ioutil.ReadFile("./schema.graphql")
	if err != nil {
		log.Fatalf("Failed to read GraphQL schema: %v", err)
	}
	schemastr := string(bytestr)
	Schema, err = graphql.ParseSchema(schemastr, &RootResolver{})
	if err != nil {
		log.Fatalf("Failed to parse GraphQL schema: %v", err)
	}
}
