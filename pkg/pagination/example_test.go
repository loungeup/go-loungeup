package pagination

func Example() {
	// Create a pager. The pager will call the readIDsPage function to read pages. The readIDsPage function uses the
	// limit and offset parameters to read a page of IDs.
	pager := NewPager(readIDsPage)

	for pager.Next() { // Iterate over pages.
		page := pager.Page() // Get the current page.

		_ = page // Process the page.
	}

	if err := pager.Err(); err != nil {
		// Check for errors.
	}
}

func readIDsPage(limit, offset int) ([]int, error) {
	return []int{}, nil // Read a page with the given limit and offset.
}
