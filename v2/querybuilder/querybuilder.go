package querybuilder

import (
	"fmt"
	"github.com/dinson/mongokit/v2/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"reflect"
	"strings"
)

type QueryBuilder struct {
	filters               []bson.D
	batchFilters          bson.M
	rawQuery              any
	aggregate             bson.A
	fullTextSearchKeyword string
	resultCount           int64
	isSetLimit            bool
	skipCount             int64
	sort                  bson.D
	error                 error
}

type Query struct {
	Filters []bson.D
	// RawQuery lets developer pass in a raw mongodb query.
	// if set, Filters and BatchFilters will be ignored,
	// and only raw query will be taken into consideration for a Find operation.
	// Cannot be used in combination with any other filter fields
	RawQuery any
	// BatchFilters helps to fetch documents against a list of objectIDs.
	// If set, `Filters` will be ignored during a Find operation.
	// Cannot be used in combination with Filters.
	BatchFilters bson.M
	Aggregate    bson.A
	// Sort is the raw sort document built by the query builder.
	// Kept alongside Options because driver v2 option builders are
	// write-only: once a setter is recorded it cannot be read back.
	Sort              any
	Options           *options.FindOptionsBuilder
	CountOptions      *options.CountOptionsBuilder
	DeleteOneOptions  *options.DeleteOneOptionsBuilder
	UpdateOneOptions  *options.UpdateOneOptionsBuilder
	DeleteManyOptions *options.DeleteManyOptionsBuilder
	UpdateManyOptions *options.UpdateManyOptionsBuilder
}

func New() *QueryBuilder {
	var filters []bson.D
	var aggregate bson.A
	return &QueryBuilder{
		filters:   filters,
		aggregate: aggregate,
		sort:      nil,
		error:     nil,
	}
}

func (b *QueryBuilder) RawQuery(q any) *QueryBuilder {
	if q != nil {
		b.rawQuery = q
	}
	return b
}

// Equals ... generic key value matching
// IMPORTANT: Will PANIC if used for matching string values.
// To filter strings, use `EqualString()` instead.
func (b *QueryBuilder) Equals(key KeyMongoDB, value any) *QueryBuilder {
	if reflect.ValueOf(value).IsNil() {
		return b
	}
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

// EqualStringArray match a string value in an array field
func (b *QueryBuilder) EqualStringArray(key KeyMongoDB, role string) *QueryBuilder {
	if len(strings.TrimSpace(role)) == 0 {
		return b
	}
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.M{"$in": bson.A{role}}}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualString(key KeyMongoDB, value string) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualsIDHex(key KeyMongoDB, idHex string) *QueryBuilder {
	if len(strings.TrimSpace(idHex)) == 0 {
		return b
	}

	oID, err := utils.StringToObjectID(idHex)
	if err != nil {
		b.error = err
		return b
	}

	filters := b.filters
	filters = append(filters, bson.D{{key.String(), oID}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualNumber(key KeyMongoDB, value float64) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualInt(key KeyMongoDB, value int) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualInt8(key KeyMongoDB, value int8) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualInt16(key KeyMongoDB, value int16) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualInt32(key KeyMongoDB, value int32) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualInt64(key KeyMongoDB, value int64) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualUint(key KeyMongoDB, value uint) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualUint8(key KeyMongoDB, value uint8) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualUint16(key KeyMongoDB, value uint16) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualUint32(key KeyMongoDB, value uint32) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualUint64(key KeyMongoDB, value uint64) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), value}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) NotEquals(key KeyMongoDB, value any) *QueryBuilder {
	if reflect.ValueOf(value).IsNil() {
		return b
	}
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$ne", value}}}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) EqualsBool(key KeyMongoDB, value bool) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$eq", value}}}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) NotEqualsBool(key KeyMongoDB, value bool) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$ne", value}}}})
	b.filters = filters
	return b
}

// IsNull ... check if value of a key is null
func (b *QueryBuilder) IsNull(key KeyMongoDB) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), nil}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) GreaterThanOrEqualTo(key KeyMongoDB, value int64) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$gte", value}}}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) LessThanOrEqualTo(key KeyMongoDB, value int64) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$lte", value}}}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) GreaterThan(key KeyMongoDB, value int64) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$gt", value}}}})
	b.filters = filters
	return b
}

func (b *QueryBuilder) LessThan(key KeyMongoDB, value int64) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$lt", value}}}})
	b.filters = filters
	return b
}

// InArray ... check if a value exist in an array field
func (b *QueryBuilder) InArray(key KeyMongoDB, value any) *QueryBuilder {
	if reflect.ValueOf(value).IsNil() {
		return b
	}
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$all", bson.A{value}}}}})
	b.filters = filters
	return b
}

// BatchGet ... retrieve multiple documents against a list of object IDs.
// If set, any other filters set will be ignored during a Find operation.
func (b *QueryBuilder) BatchGet(key KeyMongoDB, idHexList []string) *QueryBuilder {
	// Optimize for empty input
	if len(idHexList) == 0 {
		return b
	}
	// Convert hex strings to ObjectIDs with error handling
	objectIDs := make([]bson.ObjectID, 0, len(idHexList)) // Preallocate for efficiency
	for _, idHex := range idHexList {
		oID, err := utils.StringToObjectID(idHex)
		if err != nil {
			b.error = err
			return b
		}
		objectIDs = append(objectIDs, *oID) // Dereference to append the ObjectId value
	}

	filter := bson.M{key.String(): bson.M{"$in": objectIDs}}
	b.batchFilters = filter

	b.filters = []bson.D{} // resetting the existing filters if any

	return b
}

// Exists ... check if a key exists in a document
func (b *QueryBuilder) Exists(key KeyMongoDB) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$exists", true}}}})
	b.filters = filters
	return b
}

// NotExists ... check if a key does not exist in a document
func (b *QueryBuilder) NotExists(key KeyMongoDB) *QueryBuilder {
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$exists", false}}}})
	b.filters = filters
	return b
}

// StartsWith ... query for matching records that starts with value
// Eg: create the query: bson.D{{"type", bson.D{{"$regex", "^E"}}}} to match all "type" that starts with "E"
func (b *QueryBuilder) StartsWith(key KeyMongoDB, value string) *QueryBuilder {
	filters := b.filters
	s := fmt.Sprintf("^%v", value)
	filters = append(filters, bson.D{{key.String(), bson.D{{"$regex", s}}}})
	b.filters = filters
	return b
}

// SetFullTextSearch
// Perform full text search against full text index using the keyword
func (b *QueryBuilder) SetFullTextSearch(keyword string) *QueryBuilder {
	b.fullTextSearchKeyword = keyword
	filters := b.filters
	filters = append(filters, bson.D{{"$text", bson.D{{"$search", keyword}}}})
	b.filters = filters
	return b
}

// SortBySearchScore
// sorts the results by textScore in full text search field
// Note: use only when performing full text search
func (b *QueryBuilder) SortBySearchScore() *QueryBuilder {
	if len(b.fullTextSearchKeyword) == 0 {
		return b
	}
	sort := bson.D{{"score", bson.D{{"$meta", "textScore"}}}}
	b.sort = sort
	return b
}

// SortAsc ... sort ascending by field "key"
// Only allows any 1 sort
// Multiple sort chains will be ignored
func (b *QueryBuilder) SortAsc(key KeyMongoDB) *QueryBuilder {
	if b.sort != nil {
		return b // if sort is already set, skip this step
	}
	b.sort = bson.D{{key.String(), 1}}
	return b
}

// SortDesc ... sort descending by field "key"
// Only allows any 1 sort
// Multiple sort chains will be ignored
func (b *QueryBuilder) SortDesc(key KeyMongoDB) *QueryBuilder {
	if b.sort != nil {
		return b
	}
	b.sort = bson.D{{key.String(), -1}}
	return b
}

// Limit number of results returned by db query
// If limit is not set or count is 0, all records are fetched by default.
// Can be used in usual operation as well as aggregate operations
func (b *QueryBuilder) Limit(count int64) *QueryBuilder {
	b.resultCount = count
	b.isSetLimit = true
	return b
}

// Skip number of results returned by db query
// Can be used in usual operation as well as aggregate operations
func (b *QueryBuilder) Skip(count int64) *QueryBuilder {
	b.skipCount = count
	return b
}

// AfterID paginate results greater than a value for the _id.
// Note: mainly used when we are sorting results in ascending order
func (b *QueryBuilder) AfterID(idHex string) *QueryBuilder {
	if len(strings.TrimSpace(idHex)) == 0 {
		return b
	}

	afterID, err := utils.StringToObjectID(idHex)
	if err != nil {
		return b
	}

	filters := b.filters
	filters = append(filters, bson.D{{"_id", bson.D{{"$gt", afterID}}}})
	b.filters = filters

	return b
}

// BeforeID paginate results lesser than a value for the ID.
// Note: mainly used when we are sorting results in descending order
func (b *QueryBuilder) BeforeID(idHex string) *QueryBuilder {
	if len(strings.TrimSpace(idHex)) == 0 {
		return b
	}

	afterID, err := utils.StringToObjectID(idHex)
	if err != nil {
		return b
	}

	filters := b.filters
	filters = append(filters, bson.D{{"_id", bson.D{{"$lt", afterID}}}})
	b.filters = filters

	return b
}

// Build returns the built query after all the chains are complete
func (b *QueryBuilder) Build() (*Query, error) {
	opts := options.Find()

	if b.isSetLimit {
		opts.SetLimit(b.resultCount)
	}

	opts.SetSkip(b.skipCount)

	if len(b.fullTextSearchKeyword) > 0 {
		opts.SetProjection(bson.D{{"score", bson.D{{"$meta", "textScore"}}}})
	}

	q := &Query{
		Filters:          b.filters,
		RawQuery:         b.rawQuery,
		BatchFilters:     b.batchFilters,
		Options:          opts,
		DeleteOneOptions: options.DeleteOne(),
		UpdateOneOptions: options.UpdateOne(),
	}

	// guard against storing a typed-nil bson.D in the `any` field
	if b.sort != nil {
		opts.SetSort(b.sort)
		q.Sort = b.sort
	}

	if q.Filters == nil {
		q.Filters = []bson.D{}
	}

	return q, b.error
}

// Aggregate returns the aggregate query after all the chains are complete
func (b *QueryBuilder) Aggregate() (*Query, error) {
	q := &Query{
		Aggregate: b.aggregate,
	}

	if q.Aggregate == nil {
		q.Aggregate = bson.A{}
	}

	if q.Aggregate != nil && b.isSetLimit {
		q.Aggregate = append(q.Aggregate, bson.M{"$skip": b.skipCount})
		q.Aggregate = append(q.Aggregate, bson.M{"$limit": b.resultCount})
	}

	return q, b.error
}

// Match filters the documents in the aggregation pipeline based on specified criteria
// Note: It is often used early in the aggregation pipeline to reduce the number of documents processed in subsequent stages.
func (b *QueryBuilder) Match(key KeyMongoDB, value any) *QueryBuilder {
	if reflect.ValueOf(value).IsNil() {
		return b
	}
	aggregate := b.aggregate
	aggregate = append(aggregate, bson.D{{"$match", bson.D{{key.String(), value}}}})
	b.aggregate = aggregate
	return b
}

// Lookup used to perform a left outer join between documents from two different collections.
func (b *QueryBuilder) Lookup(request *LookupModel) *QueryBuilder {
	if request == nil {
		return b
	}

	pipeline := bson.A{}
	if request.Pipeline != nil {
		pipeline = request.Pipeline
	}

	aggregate := b.aggregate
	aggregate = append(
		aggregate, bson.D{
			{"$lookup", bson.D{
				{"from", request.From},
				{"localField", request.LocalField},
				{"foreignField", request.ForeignField},
				{"pipeline", pipeline},
				{"as", request.As},
			}},
		})
	b.aggregate = aggregate
	return b
}

// SortDescStage ... used in aggregation pipeline sort descending by field "key"
func (b *QueryBuilder) SortDescStage(key KeyMongoDB) *QueryBuilder {
	if len(key.String()) == 0 {
		return b
	}

	aggregate := b.aggregate
	aggregate = append(aggregate, bson.D{{"$sort", bson.D{{key.String(), -1}}}})
	b.aggregate = aggregate
	return b
}

// NotEqualStage ... used in aggregation pipeline returns not equal value of "key"
// Note: key should be a pointer
func (b *QueryBuilder) NotEqualStage(key KeyMongoDB, value any) *QueryBuilder {
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		b.error = errInvalidPointer
		return b
	}

	if reflect.ValueOf(value).IsNil() {
		return b
	}
	aggregate := b.aggregate
	aggregate = append(aggregate, bson.D{{"$match", bson.D{{key.String(), bson.D{{"$ne", value}}}}}})
	b.aggregate = aggregate
	return b
}

// RegexSearch is used to case-insensitive regex search keyword in multiple fields
func (b *QueryBuilder) RegexSearch(fields []KeyMongoDB, keyword string) *QueryBuilder {
	if len(strings.TrimSpace(keyword)) == 0 || len(fields) == 0 {
		return b
	}

	regex := bson.D{
		{"$regex", keyword},
		{"$options", "i"},
	}

	var orConditions []bson.D
	for _, field := range fields {
		orConditions = append(orConditions, bson.D{{field.String(), regex}})
	}

	filters := b.filters
	filters = append(filters, bson.D{{"$or", orConditions}})
	b.filters = filters

	return b
}

func (q *Query) GetFilter() any {
	var filter any
	if q.RawQuery != nil {
		filter = q.RawQuery
	} else {
		if q.BatchFilters != nil {
			filter = q.BatchFilters
		} else {
			if len(q.Filters) > 0 {
				filter = bson.D{{"$and", q.Filters}}
			} else {
				filter = bson.D{}
			}
		}
	}

	return filter
}

// GetDeleteFilter returns the filter to use for a delete operation.
//
// Unlike GetFilter it never falls back to an empty document: a query with no
// filters returns ErrEmptyFilter, so a delete can never match the whole
// collection by accident.
func (q *Query) GetDeleteFilter() (any, error) {
	if q == nil {
		return nil, ErrEmptyFilter
	}

	if !isMatchAll(q.RawQuery) {
		return q.RawQuery, nil
	}

	if len(q.BatchFilters) > 0 {
		return q.BatchFilters, nil
	}

	if len(q.Filters) > 0 {
		return bson.D{{"$and", q.Filters}}, nil
	}

	return nil, ErrEmptyFilter
}

// isMatchAll reports whether a filter document would match every document in
// the collection, i.e. it is nil or an empty document/array.
func isMatchAll(filter any) bool {
	if filter == nil {
		return true
	}

	v := reflect.ValueOf(filter)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		return v.Len() == 0
	default:
		return false
	}
}

// MatchAny filters documents where the specified key's value is part of the given array.
func (b *QueryBuilder) MatchAny(key KeyMongoDB, values []interface{}) *QueryBuilder {
	// Ensure the values array is not empty
	if len(values) == 0 {
		return b
	}

	// Add the $in filter to the filters slice
	filters := b.filters
	filters = append(filters, bson.D{{key.String(), bson.D{{"$in", values}}}})
	b.filters = filters

	return b
}
