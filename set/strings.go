package set

type Strings map[string]struct{}

func NewStringSet() Strings {
	return make(map[string]struct{})
}

func (s Strings) Insert(str string) {
	s[str] = struct{}{}
}

func (s Strings) Delete(str string) {
	delete(s, str)
}

func (s Strings) Contains(str string) bool {
	_, found := s[str]
	return found
}

func (s Strings) ToSlice() []string {
	r := make([]string, 0, len(s))
	for k := range s {
		r = append(r, k)
	}
	return r
}
