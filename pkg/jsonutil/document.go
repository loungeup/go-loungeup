package jsonutil

import (
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

// Document represents a JSON document.
type Document string

const NilDocument = Document("")

// ReadDocument reads a document from the given reader.
func ReadDocument(reader io.Reader) (Document, error) {
	documentAsString, err := readString(reader)
	if err != nil {
		return "", err
	}

	return Document(documentAsString), nil
}

// readString reads a string from the given reader.
func readString(reader io.Reader) (string, error) {
	builder := &strings.Builder{}
	if _, err := io.Copy(builder, reader); err != nil {
		return "", err
	}

	// Eventually close the reader.
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	return builder.String(), nil
}

// Child returns a new Document representing the child element at the specified path.
func (d Document) Child(path string) Document { return Document(gjson.Get(d.String(), path).String()) }

// Children returns an array of child documents that match the given path.
func (d Document) Children(path string) []Document {
	result := []Document{}

	element := gjson.Get(d.String(), path)
	if !element.IsArray() {
		return result
	}

	for _, child := range element.Array() {
		result = append(result, Document(child.String()))
	}

	return result
}

// Float64 returns the float64 value at the specified JSON path.
func (d Document) Float64(path string) float64 { return gjson.Get(d.String(), path).Float() }

// Int returns the integer value at the specified JSON path.
func (d Document) Int(path string) int { return int(gjson.Get(d.String(), path).Int()) }

// String returns the string representation of the document.
func (d Document) String() string { return string(d) }

// Time returns the time value at the specified JSON path.
func (d Document) Time(path string) time.Time { return gjson.Get(d.String(), path).Time() }

// UUID returns the UUID value at the specified JSON path.
func (d Document) UUID(path string) uuid.UUID {
	result, _ := uuid.Parse(gjson.Get(d.String(), path).String())

	return result
}

// UUIDs returns an array of UUIDs that match the given path.
func (d Document) UUIDs(path string) []uuid.UUID {
	result := []uuid.UUID{}

	for _, child := range d.Children(path) {
		parsedID, err := uuid.Parse(child.String())
		if err != nil {
			continue
		}

		result = append(result, parsedID)
	}

	return result
}
