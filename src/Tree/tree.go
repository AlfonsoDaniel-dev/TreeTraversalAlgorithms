package Tree

import (
	"errors"
	"fmt"

	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Node"
)

type TraversalStep struct {
	Id          int
	CurrentNode *Node.Node
	Visited     []*Node.Node
	Unvisited   []*Node.Node
}

type Tree struct {
	LastNodeId int
	Root       *Node.Node
	Nodes      map[int]*Node.Node
}

func NewTree() *Tree {

	root := Node.NewNode(0, nil)

	t := &Tree{
		Root:       root,
		LastNodeId: 0,
		Nodes:      make(map[int]*Node.Node),
	}

	t.Nodes[t.LastNodeId] = root

	return t
}

func (t *Tree) AddNodeFromRoot() error {

	if t.Root == nil {
		errStr := fmt.Sprintf("not root on the tree")
		return errors.New(errStr)
	}

	nodeId := t.LastNodeId + 1
	node := Node.NewNode(nodeId, t.Root)

	t.Root.AddChild(node)

	t.LastNodeId++

	return nil
}
func (t *Tree) AddNode(NodeId int) error {

	node, ok := t.Nodes[NodeId]
	if !ok {
		errStr := fmt.Sprintf("error: no Node with the id: %d found", NodeId)
		return errors.New(errStr)
	}

	nodeId := t.LastNodeId + 1
	newNode := Node.NewNode(nodeId, node)

	node.AddChild(newNode)

	t.Nodes[nodeId] = newNode

	t.LastNodeId++

	return nil
}

func (t *Tree) TraversalBfsSteps() ([]TraversalStep, error) {
	if t.Root == nil {
		return nil, errors.New("root is nil")
	}

	var history []TraversalStep

	queue := []*Node.Node{t.Root}
	var visited []*Node.Node

	stepId := 0

	for len(queue) > 0 {

		actual := queue[0]

		queueSnapshot := make([]*Node.Node, len(queue))
		copy(queueSnapshot, queue)

		visitedSnapshot := make([]*Node.Node, len(visited))
		copy(visitedSnapshot, visited)

		step := TraversalStep{
			Id:          stepId,
			CurrentNode: actual,
			Visited:     visitedSnapshot,
			Unvisited:   queueSnapshot,
		}

		history = append(history, step)

		queue = queue[1:]
		visited = append(visited, actual)

		for _, children := range actual.GetChildren() {
			queue = append(queue, children)
		}

		stepId++
	}

	return history, nil
}

func (t *Tree) TraversalDfsSteps() ([]TraversalStep, error) {
	
}
