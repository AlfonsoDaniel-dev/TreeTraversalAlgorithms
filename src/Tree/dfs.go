package Tree

import (
	"fmt"

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

func (t *Tree) TraversalDfsSteps(startNodeId int) ([]TraversalStep, error) {
	// 1. Buscamos el nodo de inicio
	startNode, ok := t.Nodes[startNodeId]
	if !ok {
		return nil, fmt.Errorf("no se encontro nodo con ID: %d", startNodeId)
	}

	var history []TraversalStep

	// 2. Iniciamos la pila con el nodo seleccionado
	stack := []*Node.Node{startNode}

	discovered := make(map[int]bool)
	discovered[startNode.Id] = true

	var visited []*Node.Node
	stepId := 0

	for len(stack) > 0 {

		n := len(stack) - 1
		actual := stack[n]
		stack = stack[:n]

		// 1. Explorar hacia ARRIBA (Apilamos al padre) <-- LA NUEVA MAGIA
		padre := actual.GetParent()
		if padre != nil && !discovered[padre.Id] {
			discovered[padre.Id] = true
			stack = append(stack, padre)
		}

		// 2. Explorar hacia ABAJO (Apilamos a los hijos en reversa)
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
