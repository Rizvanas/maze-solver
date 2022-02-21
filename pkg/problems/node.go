package problems

import "fmt"

type Node struct {
	X, Y        int
	Connections []*Node
}

func (n Node) Describe() string {
	return fmt.Sprintf("%d;%d", n.X, n.Y)
}
