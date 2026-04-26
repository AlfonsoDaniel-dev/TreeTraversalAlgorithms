package Game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	// Importaciones exactas de tu proyecto
	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Tree"
)

// --- ESTRUCTURAS VISUALES ---

type VisualNode struct {
	ID     int
	X, Y   float64
	VX, VY float64
	Radius float32
}

type VisualEdge struct {
	From *VisualNode
	To   *VisualNode
}

// --- ESTRUCTURA PRINCIPAL DEL JUEGO ---

type Game struct {
	BgColor      color.Color
	ScreenWidth  int
	ScreenHeight int

	Nodes map[int]*VisualNode
	Edges []*VisualEdge

	TraversalSteps []Tree.TraversalStep
	CurrentStep    int
	IsPlaying      bool
	TicksPerFrame  int
	tickCounter    int

	draggedNode *VisualNode
}

// NewGame construye el motor visual
func NewGame(bgColor color.Color, width, height int, tree *Tree.Tree) *Game {
	g := &Game{
		BgColor:       bgColor,
		ScreenWidth:   width,
		ScreenHeight:  height,
		Nodes:         make(map[int]*VisualNode),
		Edges:         make([]*VisualEdge, 0),
		TicksPerFrame: 30, // Cambia de paso cada medio segundo
	}

	if tree != nil {
		// 1. Mapear nodos lógicos a visuales (Usamos directamente tu map Nodes)
		for id := range tree.Nodes {
			g.Nodes[id] = &VisualNode{
				ID: id,
				// Esparcimos los nodos en un radio aleatorio de 100px desde el centro
				X:      float64(width/2) + (rand.Float64()*200 - 100),
				Y:      float64(height/2) + (rand.Float64()*200 - 100),
				Radius: 15.0,
			}
		}

		// 2. Mapear aristas (Edges) iterando sobre los hijos
		for _, parent := range tree.Nodes {
			for _, child := range parent.GetChildren() {
				if vParent, ok1 := g.Nodes[parent.Id]; ok1 {
					if vChild, ok2 := g.Nodes[child.Id]; ok2 {
						g.Edges = append(g.Edges, &VisualEdge{
							From: vParent,
							To:   vChild,
						})
					}
				}
			}
		}
	}

	return g
}

func (g *Game) Update() error {
	g.handleInput()
	g.updatePhysics()
	g.updatePlayback()
	return nil
}

func (g *Game) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.IsPlaying = !g.IsPlaying
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.CurrentStep < len(g.TraversalSteps)-1 {
		g.CurrentStep++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.CurrentStep > 0 {
		g.CurrentStep--
	}

	mx, my := ebiten.CursorPosition()
	fx, fy := float64(mx), float64(my)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, n := range g.Nodes {
			dx, dy := n.X-fx, n.Y-fy
			if math.Sqrt(dx*dx+dy*dy) <= float64(n.Radius) {
				g.draggedNode = n
				break
			}
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.draggedNode != nil {
		g.draggedNode.X = fx
		g.draggedNode.Y = fy
		g.draggedNode.VX = 0
		g.draggedNode.VY = 0
	} else {
		g.draggedNode = nil
	}
}

func (g *Game) updatePhysics() {
	const repulsion = 4000.0
	const springLen = 100.0
	const springK = 0.05
	const gravity = 0.01
	const friction = 0.80
	const maxSpeed = 20.0

	// 1. Repulsión
	for id1, n1 := range g.Nodes {
		for id2, n2 := range g.Nodes {
			if id1 == id2 {
				continue
			}
			dx, dy := n2.X-n1.X, n2.Y-n1.Y

			// Micro-temblor anti-singularidad
			if math.Abs(dx) < 0.1 && math.Abs(dy) < 0.1 {
				dx = rand.Float64()*2 - 1
				dy = rand.Float64()*2 - 1
			}

			distSq := dx*dx + dy*dy
			if distSq < 10.0 {
				distSq = 10.0
			}

			f := repulsion / distSq
			dist := math.Sqrt(distSq)

			n1.VX -= (dx / dist) * f
			n1.VY -= (dy / dist) * f
		}
	}

	// 2. Atracción (Resortes)
	for _, edge := range g.Edges {
		dx, dy := edge.To.X-edge.From.X, edge.To.Y-edge.From.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 1 {
			f := (dist - springLen) * springK
			edge.From.VX += (dx / dist) * f
			edge.From.VY += (dy / dist) * f
			edge.To.VX -= (dx / dist) * f
			edge.To.VY -= (dy / dist) * f
		}
	}

	// 3. Gravedad Central y Movimiento
	cx, cy := float64(g.ScreenWidth/2), float64(g.ScreenHeight/2)
	padding := float64(30.0)

	for _, n := range g.Nodes {
		if n == g.draggedNode {
			continue
		}

		// Gravedad
		n.VX += (cx - n.X) * gravity
		n.VY += (cy - n.Y) * gravity

		// Fricción
		n.VX *= friction
		n.VY *= friction

		// Clamp de velocidad (Límite máximo)
		if n.VX > maxSpeed {
			n.VX = maxSpeed
		}
		if n.VX < -maxSpeed {
			n.VX = -maxSpeed
		}
		if n.VY > maxSpeed {
			n.VY = maxSpeed
		}
		if n.VY < -maxSpeed {
			n.VY = -maxSpeed
		}

		// Aplicar movimiento
		n.X += n.VX
		n.Y += n.VY

		// ESCUDO ANTI-NAN: Si la matemática colapsa, lo reiniciamos al centro
		if math.IsNaN(n.X) || math.IsNaN(n.Y) {
			n.X = cx
			n.Y = cy
			n.VX, n.VY = 0, 0
		}

		// Paredes invisibles
		if n.X < padding {
			n.X = padding
			n.VX *= -0.5
		}
		if n.X > float64(g.ScreenWidth)-padding {
			n.X = float64(g.ScreenWidth) - padding
			n.VX *= -0.5
		}
		if n.Y < padding {
			n.Y = padding
			n.VY *= -0.5
		}
		if n.Y > float64(g.ScreenHeight)-padding {
			n.Y = float64(g.ScreenHeight) - padding
			n.VY *= -0.5
		}
	}
}

func (g *Game) updatePlayback() {
	if g.IsPlaying && len(g.TraversalSteps) > 0 {
		g.tickCounter++
		if g.tickCounter >= g.TicksPerFrame {
			g.tickCounter = 0
			if g.CurrentStep < len(g.TraversalSteps)-1 {
				g.CurrentStep++
			} else {
				g.IsPlaying = false
			}
		}
	}
}

// --- RENDERIZADO ---

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.BgColor)

	// Dibujar Aristas
	for _, edge := range g.Edges {
		vector.StrokeLine(
			screen,
			float32(edge.From.X), float32(edge.From.Y),
			float32(edge.To.X), float32(edge.To.Y),
			2, color.RGBA{100, 100, 100, 255}, true,
		)
	}

	// Analizar estado actual del algoritmo
	var currentState Tree.TraversalState
	if len(g.TraversalSteps) > 0 {
		currentState = g.TraversalSteps[g.CurrentStep].State
	}

	// Dibujar Nodos
	for id, n := range g.Nodes {
		c := color.RGBA{150, 150, 150, 255} // Gris (Unseen)

		if currentState != nil {
			if curr := currentState.GetCurrent(); curr != nil && curr.Id == id {
				c = color.RGBA{255, 50, 50, 255} // Rojo (Current)
			} else {
				isFrontier := false
				for _, f := range currentState.GetFrontier() {
					if f.Id == id {
						isFrontier = true
						break
					}
				}
				if isFrontier {
					c = color.RGBA{255, 255, 0, 255} // Amarillo (Frontier)
				} else {
					for _, v := range currentState.GetVisited() {
						if v.Id == id {
							c = color.RGBA{0, 100, 200, 255}
							break
						} // Azul (Visited)
					}
				}
			}
		}

		if n == g.draggedNode {
			c = color.RGBA{100, 255, 100, 255} // Verde brillante al arrastrar
		}

		vector.DrawFilledCircle(screen, float32(n.X), float32(n.Y), n.Radius, c, true)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}
