package Tree

import (
	"errors"

	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Node"
)

type BfsState struct {
	CurrentNode *Node.Node
	Visited     []*Node.Node
	queue       []*Node.Node
	Unvisited   []*Node.Node
}

func (s BfsState) GetCurrent() *Node.Node {
	return s.CurrentNode
}

func (s BfsState) GetVisited() []*Node.Node {
	return s.Visited
}

func (s BfsState) GetFrontier() []*Node.Node {
	return s.queue
}

func (s BfsState) GetUnseen() []*Node.Node {
	return s.Unvisited
}
func (t *Tree) TraversalBfsSteps() ([]TraversalStep, error) {
	if t.Root == nil {
		return nil, errors.New("root is nil")
	}

	var history []TraversalStep
	queue := []*Node.Node{t.Root}
	discovered := make(map[int]bool)
	discovered[t.Root.Id] = true

	var visited []*Node.Node
	stepId := 0

	for len(queue) > 0 {
		actual := queue[0]
		queue = queue[1:]

		for _, child := range actual.GetChildren() {
			if !discovered[child.Id] {
				discovered[child.Id] = true
				queue = append(queue, child)
			}
		}

		var undiscovered []*Node.Node
		for id, node := range t.Nodes {
			if !discovered[id] {
				undiscovered = append(undiscovered, node)
			}
		}

		queueSnapshot := make([]*Node.Node, len(queue))
		copy(queueSnapshot, queue)

		visitedSnapshot := make([]*Node.Node, len(visited))
		copy(visitedSnapshot, visited)

		step := BfsState{
			CurrentNode: actual,
			Visited:     visitedSnapshot,
			queue:       queueSnapshot,
			Unvisited:   undiscovered,
		}

		history = append(history, TraversalStep{
			Id:    stepId,
			State: step,
		})

		visited = append(visited, actual)
		stepId++
	}

	return history, nil
}
