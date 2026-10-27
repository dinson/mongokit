# mongokit v2

A generic MongoDB repository for Go, built on the **official MongoDB Go driver v2**
(`go.mongodb.org/mongo-driver/v2`).

v2 is a separate module from v1 and depends on the v2 driver only. The two can
coexist in one build, but a single collection object cannot be shared between
them — the driver types are unrelated.

## Features

- Generic repository interface for any entity type
- Type-safe CRUD operations
- Chainable query builder
- Deletes are refused unless the query carries a filter

## Installation

```bash
go get github.com/dinson/mongokit/v2
```

Requires Go 1.25+ and `go.mongodb.org/mongo-driver/v2`.

## Migrating from v1

| v1 | v2 |
| --- | --- |
| `go.mongodb.org/mongo-driver` | `go.mongodb.org/mongo-driver/v2` |
| `primitive.ObjectID` | `bson.ObjectID` |
| `primitive.ObjectIDFromHex` | `bson.ObjectIDFromHex` |
| `primitive.NewObjectID` | `bson.NewObjectID` |
| `Query.Options` is `*options.FindOptions` | `Query.Options` is `*options.FindOptionsBuilder` |
| `DeleteOne`/`DeleteMany` accept an empty query | empty query returns `querybuilder.ErrEmptyFilter` |

The `primitive` package is gone in driver v2; everything moved into `bson`.
Driver v2 also split each options type into an arguments struct (`options.FindOptions`)
and a builder (`options.FindOptionsBuilder`). Operations take the builder, and the
builder's recorded setters cannot be read back — if you need the sort document
after `Build()`, read `Query.Sort` rather than the options.

## Usage

### Define a model

```go
import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
    // The primary key must be named "_id" and carry "omitempty",
    // otherwise a create would write a zero ObjectID.
    ID     *bson.ObjectID `bson:"_id,omitempty"`
    Name   string         `bson:"name"`
    Email  string         `bson:"email"`
    Status string         `bson:"status"`
}
```

### Create a repository

```go
import mongokit "github.com/dinson/mongokit/v2"

// collection is a *mongo.Collection from the v2 driver
repo := mongokit.NewRepository[User](collection)
```

### Insert a document

```go
user := &User{Name: "Dan", Email: "dan@example.com", Status: "active"}

id, err := repo.Save(ctx, user, nil) // nil ID means "create"
```

### Update a document

```go
id, err := repo.Save(ctx, user, &userIDHex)
```

`Save` is an upsert keyed on `_id`. Passing a nil `ID` generates a new
`bson.ObjectID`; passing a hex string updates that document, creating it if it
does not exist.

### Insert many

```go
ids, err := repo.InsertMany(ctx, []*User{userA, userB})
```

### Retrieve multiple documents

```go
import "github.com/dinson/mongokit/v2/querybuilder"

query, err := querybuilder.New().
    EqualString("status", "active").
    SortDesc("created_at").
    Limit(20).
    Build()
if err != nil {
    return err
}

users, err := repo.FindAll(ctx, query)
```

### Retrieve a single document

```go
query, err := querybuilder.New().EqualString("email", "dan@example.com").Build()
if err != nil {
    return err
}

user, err := repo.FindOne(ctx, query)
```

`FindOne` returns `(nil, nil)` when no document matches — it does not return
`mongo.ErrNoDocuments`.

### Delete documents

```go
// delete one document by id
query, err := querybuilder.New().EqualsIDHex("_id", idHex).Build()
if err != nil {
    return err
}
err = repo.DeleteOne(ctx, query)

// delete multiple documents
query, err = querybuilder.New().EqualString("name", "Dan").Build()
if err != nil {
    return err
}
err = repo.DeleteMany(ctx, query)
```

**Filters are mandatory for deletes.** A query with no filter — including a
raw query that is an empty document — returns `querybuilder.ErrEmptyFilter`
and nothing is sent to the server, so a delete can never match the whole
collection by accident:

```go
query, _ := querybuilder.New().Build() // no filters chained

err := repo.DeleteMany(ctx, query)
if errors.Is(err, querybuilder.ErrEmptyFilter) {
    // refused before reaching MongoDB
}
```

To clear a collection deliberately, use the driver directly
(`collection.DeleteMany(ctx, bson.D{})` or `collection.Drop(ctx)`).

## Query builder

Every method returns the builder, so calls chain. `Build()` returns the
`*Query` consumed by the repository, plus the first error encountered while
chaining (an invalid ObjectID hex, for example).

### Matching

| Method | Description |
| --- | --- |
| `Equals(key, value)` | Generic equality. Panics on non-nilable values — use the typed variants for those |
| `EqualString(key, value)` | String equality |
| `EqualsIDHex(key, idHex)` | Equality against an ObjectID given as hex |
| `EqualNumber`, `EqualInt`, `EqualInt8/16/32/64` | Numeric equality |
| `EqualUint`, `EqualUint8/16/32/64` | Unsigned equality |
| `EqualsBool(key, value)` | Boolean equality |
| `EqualStringArray(key, value)` | Match a string inside an array field |
| `NotEquals(key, value)`, `NotEqualsBool(key, value)` | `$ne` |
| `IsNull(key)` | Field is null |
| `Exists(key)`, `NotExists(key)` | `$exists` |
| `GreaterThan`, `GreaterThanOrEqualTo` | `$gt`, `$gte` |
| `LessThan`, `LessThanOrEqualTo` | `$lt`, `$lte` |
| `InArray(key, value)` | `$all` against an array field |
| `MatchAny(key, values)` | `$in` against a list of values |
| `StartsWith(key, value)` | Anchored regex |
| `RegexSearch(fields, keyword)` | Case-insensitive regex across several fields |
| `BatchGet(key, idHexList)` | Fetch by a list of ObjectID hexes. **Resets any filters chained before it** |
| `RawQuery(q)` | Use a raw MongoDB query. Overrides `Filters` and `BatchFilters` |

### Sorting, paging and search

| Method | Description |
| --- | --- |
| `SortAsc(key)`, `SortDesc(key)` | Only the first sort in a chain is applied |
| `SetFullTextSearch(keyword)` | `$text` search against a text index |
| `SortBySearchScore()` | Sort by `textScore`; only after `SetFullTextSearch` |
| `Limit(count)`, `Skip(count)` | Paging |
| `AfterID(idHex)`, `BeforeID(idHex)` | Cursor paging on `_id` |

### Aggregation

`Match`, `Lookup`, `SortDescStage` and `NotEqualStage` build pipeline stages.
Finish those chains with `Aggregate()` instead of `Build()`.

```go
query, err := querybuilder.New().
    Match("status", &status).
    Lookup(&querybuilder.LookupModel{
        From:         "orders",
        LocalField:   "_id",
        ForeignField: "user_id",
        As:           "orders",
    }).
    Aggregate()
```

The repository does not expose an aggregate method yet; run `query.Aggregate`
through the driver's `collection.Aggregate`.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
