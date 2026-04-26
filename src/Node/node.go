package Node

type Node struct {
	Id       int
	Parent   *Node
	children []*Node
	// Physics properties
	PositionX float64
	PositionY float64
	VelocityX float64
	VelocityY float64
	ForceX    float64
	ForceY    float64
	Mass      float64
}

func (n *Node) GetChildren() []*Node {
	return n.children
}

func (n *Node) AddChild(child *Node) error {
	if child == nil {
		return nil
	}
	n.children = append(n.children, child)
	return nil
}

func (n *Node) GetId() int {
	return n.Id
}

func (n *Node) GetParent() *Node {
	return n.Parent
}

func (n *Node) PopChildren(id int) {
	for i, child := range n.children {
		if child.GetId() == id {
			n.children = append(n.children[:i], n.children[i+1:]...)
			return
		}
		continue
	}
}

func NewNode(id int, parent *Node) *Node {
	return &Node{
		Id:        id,
		Parent:    parent,
		children:  make([]*Node, 0),
		PositionX: 0,
		PositionY: 0,
		VelocityX: 0,
		VelocityY: 0,
		ForceX:    0,
		ForceY:    0,
		Mass:      1.0,
	}
}
