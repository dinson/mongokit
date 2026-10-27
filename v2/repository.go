package mongokit

import (
	"context"
	"github.com/dinson/mongokit/v2/querybuilder"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

const (
	connectionTimeout = 15 * time.Second
)

type Repository[T any] interface {
	// Save a json document into the collection
	//
	// Can be used for both "create" and "update" operations.
	//
	// To "Create" a new record, pass "ID" as nil
	//
	// To "Update" an existing record, pass the primary key objectID hex as "ID"
	//
	// param: entity represents the model of the collection
	Save(ctx context.Context, entity *T, ID *string) (*bson.ObjectID, error)

	// InsertMany can be used to insert multiple records into a collection.
	InsertMany(ctx context.Context, docs []*T) ([]*bson.ObjectID, error)

	// FindAll returns all the matching documents.
	FindAll(ctx context.Context, query *querybuilder.Query) ([]*T, error)

	// FindOne returns the first matching document.
	FindOne(ctx context.Context, query *querybuilder.Query) (*T, error)

	// DeleteOne deletes the first document that matches the query.
	//
	// The query must carry at least one filter. An unfiltered query returns
	// querybuilder.ErrEmptyFilter and nothing is sent to the server.
	DeleteOne(ctx context.Context, query *querybuilder.Query) error

	// DeleteMany deletes all documents matching the query.
	//
	// The query must carry at least one filter. An unfiltered query returns
	// querybuilder.ErrEmptyFilter and nothing is sent to the server, so this
	// can never empty the whole collection by accident.
	DeleteMany(ctx context.Context, query *querybuilder.Query) error
}

type repositoryImpl[T any] struct {
	collection *mongo.Collection
}

/*
		NewRepository initiates crud methods for the given collection.

	 	T represents the model of the collection.

	 	Example usage:

		model := &Users{}

		usersRepo := NewRepository[model](mongoCollectionObject)
*/
func NewRepository[T any](collection *mongo.Collection) Repository[T] {
	return &repositoryImpl[T]{
		collection: collection,
	}
}
