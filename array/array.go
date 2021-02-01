package array

func RemoveDuplicates(arr []string) []string {
	m := make(map[string]bool)
	var res []string
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			res = append(res, v)
			m[v] = true
		}
	}

	return res
}
