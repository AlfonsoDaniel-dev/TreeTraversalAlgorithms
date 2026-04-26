package Game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"

	// IMPORTANTE: Ajusta estas rutas a tu módulo
	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Tree"
)

// --- ESTADOS DEL JUEGO ---
const (
	ModeEdit = iota
	ModePlayback
)

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

type Game struct {
	BgColor      color.Color
	ScreenWidth  int
	ScreenHeight int

	Tree  *Tree.Tree // Tu árbol lógico original
	Nodes map[int]*VisualNode
	Edges []*VisualEdge

	// Máquina de Estados
	Mode           int
	SelectedNodeID int // El nodo que el maestro tiene seleccionado para agregarle hijos

	// Reproductor
	TraversalSteps []Tree.TraversalStep
	CurrentStep    int
	IsPlaying      bool
	TicksPerFrame  int
	tickCounter    int

	draggedNode *VisualNode
}

func NewGame(bgColor color.Color, width, height int, tree *Tree.Tree) *Game {
	g := &Game{
		BgColor:        bgColor,
		ScreenWidth:    width,
		ScreenHeight:   height,
		Tree:           tree,
		Nodes:          make(map[int]*VisualNode),
		Edges:          make([]*VisualEdge, 0),
		Mode:           ModeEdit,
		SelectedNodeID: 0, // Por defecto seleccionamos la raíz
		TicksPerFrame:  40,
	}
	g.syncVisuals() // Sincronización inicial
	return g
}

// syncVisuals lee tu árbol lógico y crea los elementos físicos que falten
func (g *Game) syncVisuals() {
	if g.Tree == nil {
		return
	}

	// 1. Crear nodos visuales faltantes
	for id := range g.Tree.Nodes {
		if _, exists := g.Nodes[id]; !exists {
			g.Nodes[id] = &VisualNode{
				ID:     id,
				X:      float64(g.ScreenWidth/2) + (rand.Float64()*100 - 50),
				Y:      float64(g.ScreenHeight/2) + (rand.Float64()*100 - 50),
				Radius: 16.0,
			}
		}
	}

	// 2. Reconstruir aristas
	g.Edges = make([]*VisualEdge, 0)
	for _, parent := range g.Tree.Nodes {
		for _, child := range parent.GetChildren() {
			if vParent, ok1 := g.Nodes[parent.Id]; ok1 {
				if vChild, ok2 := g.Nodes[child.Id]; ok2 {
					g.Edges = append(g.Edges, &VisualEdge{From: vParent, To: vChild})
				}
			}
		}
	}
}

func (g *Game) Update() error {
	g.handleMouse() // El ratón funciona en todos los modos (para arrastrar)

	if g.Mode == ModeEdit {
		g.handleEditMode()
	} else if g.Mode == ModePlayback {
		g.handlePlaybackMode()
	}

	g.updatePhysics()
	return nil
}

func (g *Game) handleMouse() {
	mx, my := ebiten.CursorPosition()
	fx, fy := float64(mx), float64(my)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for id, n := range g.Nodes {
			dx, dy := n.X-fx, n.Y-fy
			if math.Sqrt(dx*dx+dy*dy) <= float64(n.Radius) {
				g.draggedNode = n

				// Si estamos en modo edición, hacer clic también selecciona el nodo
				if g.Mode == ModeEdit {
					g.SelectedNodeID = id
				}
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

func (g *Game) handleEditMode() {
	// CREAR HIJO: Presionar 'C' agrega un hijo al nodo seleccionado
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		err := g.Tree.AddNode(g.SelectedNodeID)
		if err == nil {
			g.syncVisuals() // Actualizar pantalla con el nuevo nodo
		} else {
			fmt.Println("Error creando nodo:", err)
		}
	}

	// INICIAR BFS: Presionar 'B'
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		steps, err := g.Tree.TraversalBfsSteps()
		if err == nil {
			g.TraversalSteps = steps
			g.CurrentStep = 0
			g.IsPlaying = true
			g.Mode = ModePlayback
		}
	}

	// INICIAR DFS: Presionar 'D' (Asegúrate de que TraversalDfsSteps exista y devuelva lo mismo)
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		steps, err := g.Tree.TraversalDfsSteps()
		if err == nil {
			g.TraversalSteps = steps
			g.CurrentStep = 0
			g.IsPlaying = true
			g.Mode = ModePlayback
		}
	}
}

func (g *Game) handlePlaybackMode() {
	// Presionar ESC para salir del reproductor y volver a editar
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Mode = ModeEdit
		g.IsPlaying = false
		g.TraversalSteps = nil
	}

	// Controles del reproductor
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.IsPlaying = !g.IsPlaying
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.CurrentStep < len(g.TraversalSteps)-1 {
		g.CurrentStep++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.CurrentStep > 0 {
		g.CurrentStep--
	}

	// Auto-avance
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

// (La función updatePhysics() se queda EXACTAMENTE igual a la última versión "Blindada")
func (g *Game) updatePhysics() {
	const repulsion = 4000.0
	const springLen = 100.0
	const springK = 0.05
	const gravity = 0.01
	const friction = 0.80
	const maxSpeed = 20.0

	for id1, n1 := range g.Nodes {
		for id2, n2 := range g.Nodes {
			if id1 == id2 {
				continue
			}
			dx, dy := n2.X-n1.X, n2.Y-n1.Y
			if math.Abs(dx) < 0.1 && math.Abs(dy) < 0.1 {
				dx, dy = rand.Float64()*2-1, rand.Float64()*2-1
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

	cx, cy := float64(g.ScreenWidth/2), float64(g.ScreenHeight/2)
	padding := float64(30.0)

	for _, n := range g.Nodes {
		if n == g.draggedNode {
			continue
		}

		n.VX += (cx - n.X) * gravity
		n.VY += (cy - n.Y) * gravity
		n.VX *= friction
		n.VY *= friction

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

		n.X += n.VX
		n.Y += n.VY

		if math.IsNaN(n.X) || math.IsNaN(n.Y) {
			n.X, n.Y, n.VX, n.VY = cx, cy, 0, 0
		}

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

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.BgColor)

	// Dibujar Aristas
	for _, edge := range g.Edges {
		vector.StrokeLine(screen, float32(edge.From.X), float32(edge.From.Y), float32(edge.To.X), float32(edge.To.Y), 2, color.RGBA{100, 100, 100, 255}, true)
	}

	// Evaluar estado para colores
	var currentState Tree.TraversalState
	if g.Mode == ModePlayback && len(g.TraversalSteps) > 0 {
		currentState = g.TraversalSteps[g.CurrentStep].State
	}

	// Dibujar Nodos e IDs
	// ... (dentro de la función Draw)

	// Dibujar Nodos e IDs
	for id, n := range g.Nodes {
		c := color.RGBA{150, 150, 150, 255} // Unseen por defecto

		if currentState != nil {
			if curr := currentState.GetCurrent(); curr != nil && curr.Id == id {
				c = color.RGBA{255, 50, 50, 255} // Rojo (Actual)
			} else {
				isFrontier := false
				for _, f := range currentState.GetFrontier() {
					if f.Id == id {
						isFrontier = true
						break
					}
				}

				// CORRECCIÓN DE SINTAXIS AQUÍ:
				if isFrontier {
					c = color.RGBA{255, 255, 0, 255} // Amarillo (Frontera)
				} else {
					for _, v := range currentState.GetVisited() {
						if v.Id == id {
							c = color.RGBA{0, 100, 200, 255} // Azul (Visitado)
							break
						}
					}
				}
			}
		}

		// Efecto de Selección y Arrastre
		// ...

		// Efecto de Selección y Arrastre
		if g.Mode == ModeEdit && g.SelectedNodeID == id {
			// Borde blanco si está seleccionado
			vector.DrawFilledCircle(screen, float32(n.X), float32(n.Y), n.Radius+3, color.White, true)
		}
		if n == g.draggedNode {
			c = color.RGBA{100, 255, 100, 255}
		}

		vector.DrawFilledCircle(screen, float32(n.X), float32(n.Y), n.Radius, c, true)

		// Dibujar el ID del nodo
		idStr := strconv.Itoa(n.ID)
		textX := int(n.X) - (len(idStr) * 7 / 2)
		textY := int(n.Y) + 4
		text.Draw(screen, idStr, basicfont.Face7x13, textX, textY, color.Black)
	}

	// DIBUJAR LA INTERFAZ DE USUARIO (GUI)
	g.drawUI(screen)
}

func (g *Game) drawUI(screen *ebiten.Image) {
	// Fondo translúcido para el panel superior
	vector.DrawFilledRect(screen, 0, 0, float32(g.ScreenWidth), 40, color.RGBA{0, 0, 0, 180}, true)

	if g.Mode == ModeEdit {
		text.Draw(screen, "=== MODO EDICION ===", basicfont.Face7x13, 10, 25, color.RGBA{100, 255, 100, 255})
		text.Draw(screen, "Clic: Seleccionar Nodo  |  [C]: Agregar Hijo al nodo seleccionado", basicfont.Face7x13, 180, 25, color.White)
		text.Draw(screen, "[B]: Ver BFS  |  [D]: Ver DFS", basicfont.Face7x13, 700, 25, color.RGBA{255, 200, 0, 255})
	} else if g.Mode == ModePlayback {
		text.Draw(screen, "=== MODO REPRODUCCION ===", basicfont.Face7x13, 10, 25, color.RGBA{100, 200, 255, 255})

		status := "Pausado"
		if g.IsPlaying {
			status = "Reproduciendo"
		}

		text.Draw(screen, fmt.Sprintf("Paso: %d/%d (%s)", g.CurrentStep, len(g.TraversalSteps)-1, status), basicfont.Face7x13, 220, 25, color.White)
		text.Draw(screen, "[Espacio]: Play/Pausa  |  [<-] [->]: Paso manual  |  [ESC]: Volver a Editar", basicfont.Face7x13, 450, 25, color.RGBA{255, 200, 0, 255})
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}
