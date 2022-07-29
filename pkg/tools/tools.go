package tools

type SplitResult struct {
	Item []string
}
type Splitter struct {
	splitCount int
	data       []string
}

func NewSplitter(splitCount int, data []string) *Splitter {
	return &Splitter{splitCount: splitCount, data: data}
}
func (s *Splitter) GetSplitResult() (result []*SplitResult) {
	remain := len(s.data) % s.splitCount
	integer := len(s.data) / s.splitCount
	for i := 0; i < integer; i++ {
		r := new(SplitResult)
		r.Item = s.data[i*s.splitCount : (i+1)*s.splitCount]
		result = append(result, r)
	}
	if remain > 0 {
		r := new(SplitResult)
		r.Item = s.data[integer*s.splitCount:]
		result = append(result, r)
	}
	return result
}
