package querybuilder

import "go.mongodb.org/mongo-driver/v2/bson"

type LookupModel struct {
	From         string // The foreign collection - specifies the target collection from which to retrieve the documents.
	LocalField   string // specifies the field from the input documents
	ForeignField string // Field from the documents of the "from" collection - specifies the field from the target collection
	Pipeline     bson.A // additional stages for filtering the joined documents.
	As           string // Output array field - specifies the name of the new array field that will be added to the input documents
}
