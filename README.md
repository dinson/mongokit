# MongoDB Repository Pattern

A generic MongoDB repository implementation in Go that provides a reusable abstraction layer for MongoDB operations.

## Features

- Generic repository interface for any entity type
- Type-safe implementations
- Standard CRUD operations
- Query builder support
- Extensible for entity-specific operations

## Installation

v1
```bash
go get github.com/dinson/mongokit
```

v2 - uses mongo driver v2
```bash
go get github.com/dinson/mongokit/v2@v2.0.0
```

## Usage

### Create a new repository
```
type User struct {
    ID *primitive.ObjectId `bson:"_id,omitempty"` // important to name the primary key as "_id" and include "omitempty" tag
    Name string `bson:"name"`
}

repo := NewRepository[User](mongoCollection)    // mongoCollection is collection object (from official mongoDB driver)
```

### Insert new document
```
user := &User{}
id, err := repo.Save(ctx, user, nil)
```

### Update an existing document
```
user := &User{}
id, err := repo.Save(ctx, user, &userID) // userID is the id hex
```

### Retrieve multiple documents
```
query, _ := queryBuilder.New().EqualString("status", "active").Build()
users, err := repo.Find(ctx, query)
```

### Retrieve single document
```
query, _ := queryBuilder.New().EqualString("email", "user@example.com").Build()
user, err := repo.FindOne(ctx, query)
```

### Delete document
```
// delete one document by id
query := queryBuilder.New().EqualsIDHex("_id", id).Build()
err := repo.DeleteOne(ctx, query)

// delete multiple documents
query := queryBuilder.New().EqualString("name", "Dan").Build()
err := repo.DeleteMany(ctx, query)
```

## Contributing

1. Fork the repository 
2. Create your feature branch 
3. Commit your changes 
4. Push to the branch 
5. Create a new Pull Request
