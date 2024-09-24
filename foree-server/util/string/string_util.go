package string_util

func StringSet(s []string) map[string]bool {
	m := make(map[string]bool)
	for _, k := range s {
		if _, ok := m[k]; !ok {
			m[k] = true
		}
	}
	return m
}

func ChunkSlice(s []string, splitSize int) [][]string {
	numberOfSlices := len(s) / splitSize
	remainder := len(s) % splitSize

	ret := make([][]string, 0)
	start := 0
	end := 0

	for i := 0; i < numberOfSlices; i++ {
		end += splitSize
		ret = append(ret, s[start:end])
		start = end
	}

	if remainder > 0 {
		end = start + remainder
		ret = append(ret, s[start:end])
	}
	return ret
}
