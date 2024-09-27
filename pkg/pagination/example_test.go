package pagination

func Example() {
	// Create a pager. The pager will call the readIDsPage function to read pages. The readIDsPage function uses the
	// size and offset parameters to read a page of IDs.
	pager := NewPager(NewOffsetPageReader(readIDsPage))

	for pager.Next() { // Iterate over pages.
		page := pager.Page() // Get the current page.

		_ = page // Process the page.
	}

	if err := pager.Err(); err != nil {
		// Check for errors.
	}
}

func readIDsPage(size, offset int) ([]int, error) {
	return []int{}, nil // Read a page with the given size and offset.
}
