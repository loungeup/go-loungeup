package jsonutil

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDocumentChild(t *testing.T) {
	tests := map[string]struct {
		document Document
		path     string
		want     Document
	}{
		"simple child": {
			document: Document(`{"foo": "bar"}`),
			path:     "foo",
			want:     Document("bar"),
		},
		"unknown child": {
			document: Document(`{"foo": "bar"}`),
			path:     "baz",
			want:     Document(""),
		},
		"array child": {
			document: Document(`{"foo": ["bar", "baz"]}`),
			path:     "foo.1",
			want:     Document("baz"),
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.document.Child(tt.path))
		})
	}
}

func TestDocumentChildren(t *testing.T) {
	tests := map[string]struct {
		document Document
		path     string
		want     []Document
	}{
		"simple children": {
			document: Document(`{"foo": ["bar", "baz"]}`),
			path:     "foo",
			want:     []Document{Document("bar"), Document("baz")},
		},
		"invalid children": {
			document: Document(`{"foo": "bar"}`),
			path:     "foo",
			want:     []Document{},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.document.Children(tt.path))
		})
	}
}

func TestDocumentUUIDs(t *testing.T) {
	tests := map[string]struct {
		document Document
		path     string
		want     []uuid.UUID
	}{
		"root": {
			document: Document(`["e4eaaaf2-d142-11e1-b3e4-080027620cdd","e4eaaaf2-d142-11e1-b3e4-080027620cde"]`),
			path:     "@this",
			want: []uuid.UUID{
				uuid.MustParse("e4eaaaf2-d142-11e1-b3e4-080027620cdd"),
				uuid.MustParse("e4eaaaf2-d142-11e1-b3e4-080027620cde"),
			},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.document.UUIDs(tt.path))
		})
	}
}
