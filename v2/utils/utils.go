package utils

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// StringToObjectID ... convert a string to mongo ObjectID
var StringToObjectID = func(objectIDStr string) (*bson.ObjectID, error) {
	objectId, err := bson.ObjectIDFromHex(objectIDStr)
	if err != nil {
		return nil, err
	}
	return &objectId, nil
}
