package Tree

import (
	"errors"
	"fmt"

	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Node"
)

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

	t.Nodes[nodeId] = node

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

func (t *Tree) RemoveNode(id int) error {
	if id == 0 {
		return errors.New("no puedes eliminar la raiz (id 0)")
	}

	node, exists := t.Nodes[id]
	if !exists {
		return fmt.Errorf("nodo con id %d no encontrado", id)
	}

	// 1. Desvincularlo del Padre (para que ya no exista en la jerarquía)
	if node.Parent != nil {
		node.Parent.PopChildren(id)
	}

	// 2. Función recursiva para borrar el nodo y toda su descendencia del mapa maestro
	var removeRecursive func(n *Node.Node)
	removeRecursive = func(n *Node.Node) {
		for _, child := range n.GetChildren() {
			removeRecursive(child)
		}
		// Lo eliminamos del diccionario
		delete(t.Nodes, n.Id)
	}

	removeRecursive(node)
	return nil
}
