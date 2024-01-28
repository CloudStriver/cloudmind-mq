package heap

type Pair struct {
	First  int
	Second int
}
type PairHeap []Pair

func (h PairHeap) Len() int { return len(h) }
func (h PairHeap) Less(i, j int) bool {
	if h[i].First != h[j].First {
		return h[i].First < h[j].First
	} else {
		return h[i].Second < h[j].Second
	}
}
func (h PairHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *PairHeap) Push(x Pair) {
	*h = append(*h, x)
}
func (h *PairHeap) Pop() Pair {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
