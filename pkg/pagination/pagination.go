package pagination

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/loungeup/go-loungeup/pkg/errors"
)

type Pager[S ~[]E, E any] struct {
	reader pageReader[S, E]

	lastErr  error // The last error that happened while reading a page.
	lastPage S     // The most recent page read.

	size int
}

// NewPager creates a pager with the given function to read pages of type S.
func NewPager[S ~[]E, E any](reader pageReader[S, E], options ...pagerOption) *Pager[S, E] {
	const defaultSize = 25

	configuration := &pagerConfiguration{
		size: defaultSize,
	}
	for _, option := range options {
		option(configuration)
	}

	return &Pager[S, E]{
		reader: reader,
		size:   configuration.size,
	}
}

// WithPageSize sets the size of the pages to be read by the pager.
func WithPageSize(size int) pagerOption { return func(c *pagerConfiguration) { c.size = size } }

// Err returns the error, if any, that was encountered during pagination.
func (p *Pager[S, E]) Err() error { return p.lastErr }

// Next prepares the next page for reading with the [Pager.Page] method. It returns true on success, or false if there
// is no next page or an error happened while preparing it. [Pager.Err] should be called to distinguish between the two
// cases.
func (p *Pager[S, E]) Next() bool {
	// If the last page is shorter than the size, there are no more pages.
	if len(p.lastPage) != 0 && len(p.lastPage) < p.size {
		return false
	}

	page, err := p.reader.readPage(p.size)
	if err != nil {
		p.lastErr = err
		return false
	}

	if len(page) == 0 {
		return false
	}

	p.lastPage = page

	return true
}

// Page returns the last page read by the [Pager.Next] method.
func (p *Pager[S, E]) Page() S { return p.lastPage }

type pagerConfiguration struct {
	size int
}

type pagerOption func(*pagerConfiguration)

type pageReader[S ~[]E, E any] interface {
	readPage(size int) (S, error)
}

type keysetPageReader[S ~[]E, E, K any] struct {
	readPageFunc func(size int, lastKey K) (S, K, error)
	lastKey      K
}

func NewKeysetPageReader[S ~[]E, E, K any](
	readPageFunc func(size int, lastKey K) (S, K, error),
) *keysetPageReader[S, E, K] {
	return &keysetPageReader[S, E, K]{readPageFunc: readPageFunc}
}

var _ pageReader[[]any, any] = (*keysetPageReader[[]any, any, any])(nil)

func (r *keysetPageReader[S, E, K]) readPage(size int) (S, error) {
	result, lastKey, err := r.readPageFunc(size, r.lastKey)
	if err != nil {
		return nil, err
	}

	r.lastKey = lastKey

	return result, nil
}

type offsetPagerReader[S ~[]E, E any] struct {
	readPageFunc func(size, offset int) (S, error)
	offset       int
}

func NewOffsetPageReader[S ~[]E, E any](readPageFunc func(size, offset int) (S, error)) *offsetPagerReader[S, E] {
	return &offsetPagerReader[S, E]{readPageFunc: readPageFunc}
}

var _ pageReader[[]any, any] = (*offsetPagerReader[[]any, any])(nil)

func (r *offsetPagerReader[S, E]) readPage(size int) (S, error) {
	result, err := r.readPageFunc(size, r.offset)
	if err != nil {
		return nil, err
	}

	r.offset += size

	return result, nil
}

const (
	keysetSelectorLastKeyQuery = "lastKey"
	keysetSelectorSizeQuery    = "size"
)

type KeysetSelector[T any] struct {
	LastKey T
	Size    int
}

func (s *KeysetSelector[T]) Query() url.Values {
	result := url.Values{}
	result.Add(keysetSelectorLastKeyQuery, fmt.Sprint(s.LastKey))
	result.Add(keysetSelectorSizeQuery, strconv.Itoa(s.Size))

	return result
}

func ParseKeysetSelector[T any](
	query url.Values,
	parseLastKeyFunc func(key string) (T, error),
) (*KeysetSelector[T], error) {
	result := &KeysetSelector[T]{}

	if lastKeyQuery := query.Get(keysetSelectorLastKeyQuery); lastKeyQuery != "" {
		lastKey, err := parseLastKeyFunc(lastKeyQuery)
		if err != nil {
			return nil, &errors.Error{
				Code:            errors.CodeInvalid,
				Message:         "Invalid '" + keysetSelectorLastKeyQuery + "' query parameter",
				UnderlyingError: err,
			}
		}

		result.LastKey = lastKey
	}

	if sizeQuery := query.Get(keysetSelectorSizeQuery); sizeQuery != "" {
		size, err := strconv.Atoi(sizeQuery)
		if err != nil {
			return nil, &errors.Error{
				Code:            errors.CodeInvalid,
				Message:         "Invalid '" + keysetSelectorSizeQuery + "' query parameter",
				UnderlyingError: err,
			}
		}

		result.Size = size
	}

	return result, nil
}

const (
	minLimit  = 0
	minOffset = 0
)

// Limit of a page.
type Limit int

// NewLimit creates a new limit from the provided value.
func NewLimit(limit int) Limit { return Limit(limit) }

// Bound the limit to the minimum value (zero) and the provided maximum value.
func (l Limit) Bound(maxLimit int) int {
	if l > minLimit && int(l) < maxLimit {
		return int(l)
	}

	return maxLimit
}

// Offset of a page.
type Offset int

// NewOffset creates a new offset from the provided value.
func NewOffset(offset int) Offset { return Offset(offset) }

// Bound the offset to the minimum value (zero).
func (o Offset) Bound() int {
	if o > minOffset {
		return int(o)
	}

	return minOffset
}
