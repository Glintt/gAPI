package utils

func ArrayContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveStringFromArray(s []string, r string) []string {
	var s2 []string
	for _, v := range s {
		if v != r {
			s2 = append(s2, v)
		}
	}
	return s2
}
