package authproxy

func inArray(val string, array []string) (exists bool) {
	exists = false

	for _, v := range array {
		if val == v {
			exists = true
			return
		}
	}

	return
}
