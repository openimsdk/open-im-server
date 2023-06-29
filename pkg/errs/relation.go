package errs

var Relation = &relation{m: make(map[int]map[int]struct{})}

func init() {
	Relation.Add(RecordNotFoundError, UserIDNotFoundError)
	Relation.Add(RecordNotFoundError, GroupIDNotFoundError)
	Relation.Add(DuplicateKeyError, UserIDExisted)
	Relation.Add(DuplicateKeyError, GroupIDExisted)
}

type relation struct {
	m map[int]map[int]struct{}
}

func (r *relation) Add(codes ...int) {
	if len(codes) < 2 {
		panic("codes length must be greater than 2")
	}
	for i := 1; i < len(codes); i++ {
		parent := codes[i-1]
		s, ok := r.m[parent]
		if !ok {
			s = make(map[int]struct{})
			r.m[parent] = s
		}
		for _, code := range codes[i:] {
			s[code] = struct{}{}
		}
	}
}

func (r *relation) Is(parent, child int) bool {
	if parent == child {
		return true
	}
	s, ok := r.m[parent]
	if !ok {
		return false
	}
	_, ok = s[child]
	return ok
}
