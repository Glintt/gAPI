package servicediscovery


func pageFromTo(page int, arrayLen int) (int , int) {
	from := ((page-1) * PAGE_LENGTH)
	to := from + PAGE_LENGTH
	if from >= arrayLen {
		from = arrayLen - PAGE_LENGTH
		if from < 0 {
			from = 0
		}
	}

	if to >= arrayLen {
		to = arrayLen
	}

	return from, to
}