package pagination

import (
	"fmt"
	"net/url"
	"strconv"

	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/loungeup/go-loungeup/errors"
)

type Pager[S ~[]E, E any, R PageReader[S, E]] struct {
	Reader R

	lastErr  error // The last error that happened while reading a page.
	lastPage S     // The most recent page read.

	allowShorterPages bool
	size              int
}

type (
	ESKeysetPager[S ~[]E, E, K any]          = Pager[S, E, *ESKeysetPageReader[S, E, K]]
	ESCompositeKeysetPager[S ~[]E, E, K any] = Pager[S, E, *ESCompositeKeysetPageReader[S, E, K]]
	KeysetPager[S ~[]E, E, K any]            = Pager[S, E, *KeysetPageReader[S, E, K]]
	OffsetPager[S ~[]E, E any]               = Pager[S, E, *OffsetPagerReader[S, E]]
)

// NewPager creates a pager with the given function to read pages of type S.
func NewPager[S ~[]E, E any, R PageReader[S, E]](reader R, options ...PagerOption) *Pager[S, E, R] {
	const defaultSize = 25

	configuration := &PagerConfig{
		allowShorterPages: false,
		size:              defaultSize,
	}
	for _, option := range options {
		option(configuration)
	}

	return &Pager[S, E, R]{
		Reader:            reader,
		allowShorterPages: configuration.allowShorterPages,
		size:              configuration.size,
	}
}

// AllowShorterPages allows the pager to continue reading pages even if the last page is shorter than the size.
func AllowShorterPages() PagerOption {
	return func(config *PagerConfig) { config.allowShorterPages = true }
}

// WithPageSize sets the size of the pages to be read by the pager.
func WithPageSize(size int) PagerOption {
	return func(config *PagerConfig) { config.size = size }
}

// Err returns the error, if any, that was encountered during pagination.
func (p *Pager[S, E, R]) Err() error { return p.lastErr }

// Next prepares the next page for reading with the [Pager.Page] method. It returns true on success, or false if there
// is no next page or an error happened while preparing it. [Pager.Err] should be called to distinguish between the two
// cases.
func (p *Pager[S, E, R]) Next() bool {
	if !p.allowShorterPages {
		// If the last page is shorter than the size, there are no more pages.
		if p.lastPage != nil && len(p.lastPage) < p.size {
			return false
		}
	}

	page, err := p.Reader.ReadPage(p.size)
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
func (p *Pager[S, E, R]) Page() S { return p.lastPage }

// Reset the pager to its initial state.
func (p *Pager[S, E, R]) Reset() {
	p.Reader.Reset()
	p.lastErr = nil
	p.lastPage = nil
}

type PagerConfig struct {
	size              int
	allowShorterPages bool
}

type PagerOption func(config *PagerConfig)

type PageReader[S ~[]E, E any] interface {
	ReadPage(size int) (S, error)
	Reset()
}

type KeysetPageReader[S ~[]E, E, K any] struct {
	LastKey      K
	readPageFunc func(size int, lastKey K) (S, K, error)
}

type KeysetPageReaderConfig[K any] struct {
	lastKey K
}

type KeysetPageReaderOption[K any] func(config *KeysetPageReaderConfig[K])

func NewKeysetPageReader[S ~[]E, E, K any](
	readPageFunc func(size int, lastKey K) (S, K, error),
	options ...KeysetPageReaderOption[K],
) *KeysetPageReader[S, E, K] {
	config := &KeysetPageReaderConfig[K]{}
	for _, option := range options {
		option(config)
	}

	return &KeysetPageReader[S, E, K]{
		readPageFunc: readPageFunc,
		LastKey:      config.lastKey,
	}
}

func WithKeysetPageReaderLastKey[K any](lastKey K) KeysetPageReaderOption[K] {
	return func(config *KeysetPageReaderConfig[K]) { config.lastKey = lastKey }
}

var _ PageReader[[]any, any] = (*KeysetPageReader[[]any, any, any])(nil)

func (r *KeysetPageReader[S, E, K]) ReadPage(size int) (S, error) {
	result, lastKey, err := r.readPageFunc(size, r.LastKey)
	if err != nil {
		return nil, err
	}

	r.LastKey = lastKey

	return result, nil
}

func (r *KeysetPageReader[S, E, K]) Reset() {
	var emptyKey K
	r.LastKey = emptyKey
}

type ESKeysetPageReader[S ~[]E, E, K any] struct {
	LastKey      K
	readPageFunc func(size int, lastKey K, pit *estypes.PointInTimeReference, query *estypes.Query) (S, K, error)
	pit          *estypes.PointInTimeReference
	query        *estypes.Query
}

type ESKeysetPageReaderConfig[K any] struct {
	lastKey K
	pit     *estypes.PointInTimeReference
	query   *estypes.Query
}

type ESKeysetPageReaderOption[K any] func(config *ESKeysetPageReaderConfig[K])

func WithESKeysetPageReaderLastKey[K any](lastKey K) ESKeysetPageReaderOption[K] {
	return func(config *ESKeysetPageReaderConfig[K]) { config.lastKey = lastKey }
}

func WithESKeysetPageReaderPIT[K any](pit *estypes.PointInTimeReference) ESKeysetPageReaderOption[K] {
	return func(config *ESKeysetPageReaderConfig[K]) { config.pit = pit }
}

func WithESKeysetPageReaderQuery[K any](query *estypes.Query) ESKeysetPageReaderOption[K] {
	return func(config *ESKeysetPageReaderConfig[K]) { config.query = query }
}

func NewESKeysetPageReader[S ~[]E, E, K any](
	readPageFunc func(size int, lastKey K, pit *estypes.PointInTimeReference, query *estypes.Query) (S, K, error),
	options ...ESKeysetPageReaderOption[K],
) *ESKeysetPageReader[S, E, K] {
	config := &ESKeysetPageReaderConfig[K]{}

	for _, option := range options {
		option(config)
	}

	return &ESKeysetPageReader[S, E, K]{
		LastKey:      config.lastKey,
		readPageFunc: readPageFunc,
		pit:          config.pit,
		query:        config.query,
	}
}

var _ PageReader[[]any, any] = (*ESKeysetPageReader[[]any, any, any])(nil)

func (r *ESKeysetPageReader[S, E, K]) ReadPage(size int) (S, error) {
	result, lastKey, err := r.readPageFunc(size, r.LastKey, r.pit, r.query)
	if err != nil {
		return nil, err
	}

	r.LastKey = lastKey

	return result, nil
}

func (r *ESKeysetPageReader[S, E, K]) Reset() {
	var emptyKey K
	r.LastKey = emptyKey
	r.pit = nil
}

type ESCompositeKeysetPageReader[S ~[]E, E, K any] struct {
	LastKey       K
	readPageFunc  func(size int, lastKey K, pit *estypes.PointInTimeReference, query *estypes.Query, aggs map[string]estypes.Aggregations) (S, K, error)
	pit           *estypes.PointInTimeReference
	query         *estypes.Query
	compositeAgg  map[string]estypes.Aggregations
	compositeSize int
}

type ESCompositeKeysetPageReaderConfig[K any] struct {
	lastKey       K
	pit           *estypes.PointInTimeReference
	query         *estypes.Query
	compositeAgg  map[string]estypes.Aggregations
	compositeSize int
}

type ESCompositeKeysetPageReaderOption[K any] func(config *ESCompositeKeysetPageReaderConfig[K])

func WithESCompositeKeysetPageReaderLastKey[K any](lastKey K) ESCompositeKeysetPageReaderOption[K] {
	return func(config *ESCompositeKeysetPageReaderConfig[K]) { config.lastKey = lastKey }
}

func WithESCompositeKeysetPageReaderPIT[K any](pit *estypes.PointInTimeReference) ESCompositeKeysetPageReaderOption[K] {
	return func(config *ESCompositeKeysetPageReaderConfig[K]) { config.pit = pit }
}

func WithESCompositeKeysetPageReaderQuery[K any](query *estypes.Query) ESCompositeKeysetPageReaderOption[K] {
	return func(config *ESCompositeKeysetPageReaderConfig[K]) { config.query = query }
}

func WithESCompositeKeysetPageReaderAgg[K any](compositeAgg map[string]estypes.Aggregations) ESCompositeKeysetPageReaderOption[K] {
	return func(config *ESCompositeKeysetPageReaderConfig[K]) { config.compositeAgg = compositeAgg }
}

func WithESCompositeKeysetPageReaderSize[K any](compositeSize int) ESCompositeKeysetPageReaderOption[K] {
	return func(config *ESCompositeKeysetPageReaderConfig[K]) { config.compositeSize = compositeSize }
}

func NewESCompositeKeysetPageReader[S ~[]E, E, K any](
	readPageFunc func(size int, lastKey K, pit *estypes.PointInTimeReference, query *estypes.Query, aggs map[string]estypes.Aggregations) (S, K, error),
	options ...ESCompositeKeysetPageReaderOption[K],
) *ESCompositeKeysetPageReader[S, E, K] {
	config := &ESCompositeKeysetPageReaderConfig[K]{}

	for _, option := range options {
		option(config)
	}

	return &ESCompositeKeysetPageReader[S, E, K]{
		LastKey:       config.lastKey,
		readPageFunc:  readPageFunc,
		pit:           config.pit,
		query:         config.query,
		compositeAgg:  config.compositeAgg,
		compositeSize: config.compositeSize,
	}
}

var _ PageReader[[]any, any] = (*ESCompositeKeysetPageReader[[]any, any, any])(nil)

func (r *ESCompositeKeysetPageReader[S, E, K]) ReadPage(size int) (S, error) {
	result, lastKey, err := r.readPageFunc(r.compositeSize, r.LastKey, r.pit, r.query, r.compositeAgg)
	if err != nil {
		return nil, err
	}

	r.LastKey = lastKey

	return result, nil
}

func (r *ESCompositeKeysetPageReader[S, E, K]) Reset() {
	var emptyKey K
	r.LastKey = emptyKey
	r.pit = nil
}

type OffsetPagerReader[S ~[]E, E any] struct {
	readPageFunc func(size, offset int) (S, error)
	offset       int
}

func NewOffsetPageReader[S ~[]E, E any](readPageFunc func(size, offset int) (S, error)) *OffsetPagerReader[S, E] {
	return &OffsetPagerReader[S, E]{readPageFunc: readPageFunc}
}

var _ PageReader[[]any, any] = (*OffsetPagerReader[[]any, any])(nil)

func (r *OffsetPagerReader[S, E]) ReadPage(size int) (S, error) {
	result, err := r.readPageFunc(size, r.offset)
	if err != nil {
		return nil, err
	}

	r.offset += size

	return result, nil
}

func (r *OffsetPagerReader[S, E]) Reset() { r.offset = 0 }

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
