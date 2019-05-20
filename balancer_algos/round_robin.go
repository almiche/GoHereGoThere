package balancer_algos

type roundRobin struct {
	next_node int // Previous response returned
	nodes     []string
}

func NewRoundRobin() *roundRobin {
	return &roundRobin{
		next_node: 0,
		nodes:     nil,
	}
}

func (r roundRobin) GetNodes() []string {
	return r.nodes
}

func (r *roundRobin) SetNodes(nodes []string) {
	r.nodes = nodes
}

func (r *roundRobin) Balance() string {
	response := r.nodes[r.next_node]
	r.next_node = 1 + r.next_node
	if r.next_node == len(r.nodes) {
		r.next_node = 0
	}
	return response
}
