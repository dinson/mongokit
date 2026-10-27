package mongokit

import (
	"context"
	"github.com/dinson/mongokit/v2/querybuilder"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
)

func (r repositoryImpl[T]) FindAll(ctx context.Context, filter *querybuilder.Query) ([]*T, error) {
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

	var findOpts []options.Lister[options.FindOptions]
	if filter.Options != nil {
		findOpts = append(findOpts, filter.Options)
	}

	cursor, err := r.collection.Find(newCtx, filters, findOpts...)
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		if cursor != nil {
			err = cursor.Close(ctx)
			if err != nil {
				log.Println(err)
			}
		}
	}(cursor, newCtx)

	if err != nil {
		return nil, err
	}

	var resp []*T

	for cursor.Next(newCtx) {
		var item *T

		err = cursor.Decode(&item)
		if err != nil {
			return nil, err
		}

		resp = append(resp, item)
	}

	return resp, nil
}
