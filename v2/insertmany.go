package mongokit

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r repositoryImpl[T]) InsertMany(ctx context.Context, docs []*T) ([]*bson.ObjectID, error) {
	newCtx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	var docInterfaceList []any

	for _, d := range docs {
		docInterfaceList = append(docInterfaceList, d)
	}

	res, err := r.collection.InsertMany(newCtx, docInterfaceList)
	if err != nil {
		return nil, err
	}

	insertedIDs := res.InsertedIDs

	var primitiveIDs []*bson.ObjectID

	for _, i := range insertedIDs {
		if i != nil {
			primitiveOID := i.(bson.ObjectID)
			primitiveIDs = append(primitiveIDs, &primitiveOID)
		}
	}

	return primitiveIDs, nil
}
