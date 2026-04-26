package Tree

import (
	"errors"

	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Node"
)

type DfsState struct {
	CurrentNode *Node.Node
	Visited     []*Node.Node
	stack       []*Node.Node
	Unvisited   []*Node.Node
}

func (s DfsState) GetCurrent() *Node.Node    { return s.CurrentNode }
func (s DfsState) GetVisited() []*Node.Node  { return s.Visited }
func (s DfsState) GetFrontier() []*Node.Node { return s.stack }
func (s DfsState) GetUnseen() []*Node.Node   { return s.Unvisited }

func (t *Tree) TraversalDfsSteps() ([]TraversalStep, error) {
	if t.Root == nil {
		return nil, errors.New("no root in the tree")
	}

	var history []TraversalStep
	stack := []*Node.Node{t.Root}

	discovered := make(map[int]bool)
	discovered[t.Root.Id] = true

	var visited []*Node.Node
	stepId := 0

	for len(stack) > 0 {

		n := len(stack) - 1

		actual := stack[n]
		stack = stack[:n]

		children := actual.GetChildren()

		for i := len(children) - 1; i >= 0; i-- {
			child := children[i]
			if !discovered[child.Id] {
				discovered[child.Id] = true
				stack = append(stack, child)
			}
		}

		var undiscoverd []*Node.Node
		for id, node := range t.Nodes {
			if !discovered[id] {
				undiscoverd = append(undiscoverd, node)
			}
		}

		stackSnapshot := make([]*Node.Node, len(stack))
		copy(stackSnapshot, stack)

		visitedSnapshot := make([]*Node.Node, len(visited))
		copy(visitedSnapshot, visited)

		state := DfsState{
			CurrentNode: actual,
			Visited:     visitedSnapshot,
			stack:       stackSnapshot,
			Unvisited:   undiscoverd,
		}

		history = append(history, TraversalStep{
			Id:    stepId,
			State: state,
		})

		visited = append(visited, actual)
		stepId++
	}

	return history, nil
}
