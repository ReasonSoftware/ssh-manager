package app

// SliceToString is a helper function to format a slice of strings
// into a comma separated string.
func SliceToString(slice []string) string {
	var o string
	for i, item := range slice {
		if i == 0 {
			o = item
		} else {
			o = o + ", " + item
		}
	}

	if o == "" {
		return "none"
	}

	return o
}
