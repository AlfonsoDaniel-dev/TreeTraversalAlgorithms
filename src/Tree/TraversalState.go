package Tree

import "github.com/AlfonsoDaniel-dev/TreeTraversal/src/Node"

type TraversalState interface {
	GetCurrent() *Node.Node
	GetVisited() []*Node.Node
	GetFrontier() []*Node.Node
	GetUnseen() []*Node.Node
	GetPathTaken() []int
}
