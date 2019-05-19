package balancer_algos

type Balancer interface {
	GetNodes() []string
	Balance() string
}