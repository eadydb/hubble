package stringslice

// Contains returns true if a string slice contains the given string
func Contains(sl []string, s string) bool {
	return Index(sl, s) >= 0
}

// Index returns the index of a first occurrence of a string within a string slice
func Index(sl []string, s string) int {
	for i, a := range sl {
		if a == s {
			return i
		}
	}
	return -1
}

// Insert inserts a string slice into another string slice at the given index
func Insert(sl []string, index int, insert []string) []string {
	newSlice := make([]string, len(sl)+len(insert))
	copy(newSlice[0:index], sl[0:index])
	copy(newSlice[index:index+len(insert)], insert)
	copy(newSlice[index+len(insert):], sl[index:])
	return newSlice
}

// Remove removes a string from a slice of strings
func Remove(s []string, target string) []string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == target {
			s = append(s[:i], s[i+1:]...)
		}
	}
	return s
}
