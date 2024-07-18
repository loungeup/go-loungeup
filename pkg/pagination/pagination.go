package pagination

const defaultLimit = 25

// PageReaderFunc is a function that reads a page of elements of type S with the given limit and offset.
type PageReaderFunc[S ~[]E, E any] func(limit, offset int) (S, error)

type Pager[S ~[]E, E any] struct {
	reader PageReaderFunc[S, E]

	lastErr  error // The last error that happened while reading a page.
	lastPage S     // The most recent page read.

	limit, offset int
}

// NewPager creates a pager with the given function to read pages of type S.
func NewPager[S ~[]E, E any](reader PageReaderFunc[S, E], options ...option) *Pager[S, E] {
	configuration := &configuration{
		limit: defaultLimit,
	}
	for _, option := range options {
		option(configuration)
	}

	return &Pager[S, E]{
		reader: reader,
		limit:  configuration.limit,
	}
}

// WithLimit sets the maximum number of elements to read per page.
func WithLimit(limit int) option { return func(c *configuration) { c.limit = limit } }

// Err returns the error, if any, that was encountered during pagination.
func (p *Pager[S, E]) Err() error { return p.lastErr }

// Next prepares the next page for reading with the [Pager.Page] method. It returns true on success, or false if there
// is no next page or an error happened while preparing it. [Pager.Err] should be called to distinguish between the two
// cases.
func (p *Pager[S, E]) Next() bool {
	// If the last page is shorter than the limit, there are no more pages.
	if len(p.lastPage) != 0 && len(p.lastPage) < p.limit {
		return false
	}

	page, err := p.reader(p.limit, p.offset)
	if err != nil {
		p.lastErr = err
		return false
	}

	if len(page) == 0 {
		return false
	}

	p.lastPage = page
	p.offset += p.limit

	return true
}

// Page returns the last page read by the [Pager.Next] method.
func (p *Pager[S, E]) Page() S { return p.lastPage }

type configuration struct {
	limit int
}

type option func(*configuration)
