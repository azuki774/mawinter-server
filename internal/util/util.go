package util

// InsertIfNotExists inserts a string into a slice if it does not already exist
func InsertStringIfNotExists(slice *[]string, value string) {
	for _, v := range *slice {
		if v == value {
			return
		}
	}
	*slice = append(*slice, value)
}
