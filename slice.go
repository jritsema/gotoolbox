package gotoolbox

import "strings"

//SliceContains returns true if a slice contains a string
func SliceContains(s *[]string, e string) bool {
	for _, str := range *s {
		if str == e {
			return true
		}
	}
	return false
}

//SliceContainsLike returns true if a slice contains a string
//that contains another string
func SliceContainsLike(s *[]string, e string) bool {
	for _, str := range *s {
		if strings.Contains(str, e) {
			return true
		}
	}
	return false
}
