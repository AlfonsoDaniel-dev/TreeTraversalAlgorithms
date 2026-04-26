package Tree

import (
	"fmt"

	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Node"
)

type BfsState struct {
	CurrentNode *Node.Node
	Visited     []*Node.Node
	queue       []*Node.Node
	Unvisited   []*Node.Node
	PathTaken   []int
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

func (s BfsState) GetPathTaken() []int {
	return s.PathTaken
}

func (t *Tree) TraversalBfsSteps(startNodeId int) ([]TraversalStep, error) {
	// 1. Buscamos el nodo de inicio en el mapa maestro
	startNode, ok := t.Nodes[startNodeId]
	if !ok {
		return nil, fmt.Errorf("no se encontro nodo con ID: %d", startNodeId)
	}

	var history []TraversalStep

	// 2. Iniciamos la cola con el nodo seleccionado
	queue := []*Node.Node{startNode}

	discovered := make(map[int]bool)
	discovered[startNode.Id] = true

	var visited []*Node.Node
	stepId := 0

	var pathSoFar []int

	for len(queue) > 0 {
		actual := queue[0]
		queue = queue[1:]

		pathSoFar = append(pathSoFar, actual.Id)

		for _, child := range actual.GetChildren() {
			if !discovered[child.Id] {
				discovered[child.Id] = true
				queue = append(queue, child)
			}
		}

		padre := actual.GetParent()
		if padre != nil && !discovered[padre.Id] {
			discovered[padre.Id] = true
			queue = append(queue, padre)
		}

		var undiscovered []*Node.Node
		for id, node := range t.Nodes {
			if !discovered[id] {
				undiscovered = append(undiscovered, node)
			}
		}

		pathSoFarSnapShot := make([]int, len(pathSoFar))
		copy(pathSoFarSnapShot, pathSoFar)

		queueSnapshot := make([]*Node.Node, len(queue))
		copy(queueSnapshot, queue)

		visitedSnapshot := make([]*Node.Node, len(visited))
		copy(visitedSnapshot, visited)

		step := BfsState{
			CurrentNode: actual,
			Visited:     visitedSnapshot,
			queue:       queueSnapshot,
			Unvisited:   undiscovered,
			PathTaken:   pathSoFarSnapShot,
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
