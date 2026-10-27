package mongokit

import (
	"context"
	"errors"
	"github.com/dinson/mongokit/v2/querybuilder"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r repositoryImpl[T]) FindOne(ctx context.Context, filter *querybuilder.Query) (*T, error) {
	newCtx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	var filters any
	if filter.RawQuery != nil {
		filters = filter.RawQuery
	} else {
		if filter.BatchFilters != nil {
			filters = filter.BatchFilters
		} else {
			if len(filter.Filters) > 0 {
				filters = bson.D{{"$and", filter.Filters}}
			} else {
				filters = bson.D{}
			}
		}
	}

	findOneOptions := options.FindOne()
	if filter.Sort != nil {
		findOneOptions.SetSort(filter.Sort)
	}

	result := r.collection.FindOne(newCtx, filters, findOneOptions)

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, nil
	}

	var resp *T

	if err := result.Decode(&resp); err != nil {
		return nil, err
	}

	return resp, nil
}
