package mongokit

import (
	"context"
	"github.com/dinson/mongokit/v2/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r repositoryImpl[T]) Save(ctx context.Context, entity *T, ID *string) (*bson.ObjectID, error) {
	newCtx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	objectID := bson.NewObjectID()
	if ID != nil {
		oID, err := utils.StringToObjectID(*ID)
		if err != nil {
			return nil, err
		}
		objectID = *oID
	}

	opts := options.UpdateOne().SetUpsert(true)
	filter := bson.D{{"_id", objectID}}

	res, err := r.collection.UpdateOne(newCtx, filter, bson.D{{"$set", entity}}, opts)
	if err != nil {
		return nil, err
	}

	oid := objectID
	if res.UpsertedID != nil {
		oid = res.UpsertedID.(bson.ObjectID)
	}

	return &oid, nil
}
