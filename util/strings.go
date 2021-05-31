package util

type StringSlice []string

func (strs StringSlice) Has(str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}
