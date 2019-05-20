package balancer_algos

type Balancer interface {
	GetNodes() []string
	SetNodes([]string)
	Balance() *string
}

func MapOfAlgos() map[string]Balancer {
	return map[string]Balancer{
		"RoundRobin": NewRoundRobin(),
	}
}
