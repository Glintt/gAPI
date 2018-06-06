package database


func PageFromTo(page int, page_length int, arrayLen int) (int , int) {
	from := ((page-1) * page_length)
	to := from + page_length
	if from >= arrayLen {
		from = arrayLen - page_length
		if from < 0 {
			from = 0
		}
	}

	if to >= arrayLen {
		to = arrayLen
	}

	return from, to
}