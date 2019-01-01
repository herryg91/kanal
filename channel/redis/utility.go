package chRedis

func MapStringBoolToArrOfStr(i map[string]bool) []string {
	o := []string{}
	for k, _ := range i {
		o = append(o, k)
	}
	return o
}
