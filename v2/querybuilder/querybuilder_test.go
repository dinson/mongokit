package querybuilder

import (
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestGetDeleteFilterRejectsUnfilteredQueries(t *testing.T) {
	tests := []struct {
		name  string
		query *Query
	}{
		{"nil query", nil},
		{"zero value query", &Query{}},
		{"empty filter slice", &Query{Filters: []bson.D{}}},
		{"empty raw query document", &Query{RawQuery: bson.D{}}},
		{"empty raw query map", &Query{RawQuery: bson.M{}}},
		{"nil batch filters", &Query{BatchFilters: bson.M{}}},
		{"builder with no filters", mustBuild(t, New())},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := tt.query.GetDeleteFilter()
			if !errors.Is(err, ErrEmptyFilter) {
				t.Fatalf("got err %v, want ErrEmptyFilter", err)
			}
			if filter != nil {
				t.Fatalf("got filter %v, want nil", filter)
			}
		})
	}
}

func TestGetDeleteFilterAcceptsFilteredQueries(t *testing.T) {
	oID := bson.NewObjectID()

	tests := []struct {
		name  string
		query *Query
	}{
		{"builder filters", mustBuild(t, New().EqualString("name", "Dan"))},
		{"batch get", mustBuild(t, New().BatchGet("_id", []string{oID.Hex()}))},
		{"raw query", &Query{RawQuery: bson.D{{"name", "Dan"}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := tt.query.GetDeleteFilter()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if filter == nil {
				t.Fatal("got nil filter, want a filter document")
			}
		})
	}
}

// BatchGet resets any previously chained filters, so a batch delete must not
// silently widen to the whole collection when the ID list is empty.
func TestGetDeleteFilterEmptyBatchGetKeepsExistingFilters(t *testing.T) {
	q := mustBuild(t, New().EqualString("name", "Dan").BatchGet("_id", nil))

	filter, err := q.GetDeleteFilter()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filter == nil {
		t.Fatal("got nil filter, want the chained name filter")
	}
}

func mustBuild(t *testing.T, b *QueryBuilder) *Query {
	t.Helper()
	q, err := b.Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}
	return q
}
