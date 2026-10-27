package mongokit

import (
	"context"
	"github.com/dinson/mongokit/v2/querybuilder"
)

func (r repositoryImpl[T]) DeleteOne(ctx context.Context, query *querybuilder.Query) error {
	newCtx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	filters, err := query.GetDeleteFilter()
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(newCtx, filters)
	if err != nil {
		return err
	}

	return nil
}
