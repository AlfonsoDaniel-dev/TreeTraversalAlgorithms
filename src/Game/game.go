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

	// Ajusta esta ruta a tu módulo
	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Tree"
)

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

	CatalogTrees     []*Tree.Tree
	CatalogNames     []string
	CurrentTreeIndex int

	Tree  *Tree.Tree
	Nodes map[int]*VisualNode
	Edges []*VisualEdge

	Mode           int
	SelectedNodeID int

	TraversalSteps []Tree.TraversalStep
	CurrentStep    int
	IsPlaying      bool
	TicksPerFrame  int // Control de velocidad
	tickCounter    int

	draggedNode *VisualNode
}

func NewGame(bgColor color.Color, width, height int, initialTree *Tree.Tree) *Game {
	g := &Game{
		BgColor:        bgColor,
		ScreenWidth:    width,
		ScreenHeight:   height,
		Tree:           initialTree,
		Nodes:          make(map[int]*VisualNode),
		Edges:          make([]*VisualEdge, 0),
		Mode:           ModeEdit,
		SelectedNodeID: 0,
		TicksPerFrame:  30, // Velocidad inicial (0.5s por paso)
	}
	g.syncVisuals()
	return g
}

func (g *Game) LoadTreeFromCatalog(index int) {
	if index < 0 || index >= len(g.CatalogTrees) {
		return
	}
	g.CurrentTreeIndex = index
	g.Tree = g.CatalogTrees[index]
	g.Nodes = make(map[int]*VisualNode)
	g.Edges = make([]*VisualEdge, 0)
	g.TraversalSteps = nil
	g.CurrentStep, g.Mode, g.IsPlaying, g.SelectedNodeID = 0, ModeEdit, false, 0
	g.syncVisuals()
}

func (g *Game) syncVisuals() {
	if g.Tree == nil {
		return
	}
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
	g.Edges = make([]*VisualEdge, 0)
	for _, parent := range g.Tree.Nodes {
		for _, child := range parent.GetChildren() {
			if vP, ok1 := g.Nodes[parent.Id]; ok1 {
				if vC, ok2 := g.Nodes[child.Id]; ok2 {
					g.Edges = append(g.Edges, &VisualEdge{From: vP, To: vC})
				}
			}
		}
	}
}

func (g *Game) Update() error {
	g.handleMouse()
	// Selector de velocidad (Funciona en ambos modos)
	if ebiten.IsKeyPressed(ebiten.KeyUp) && g.TicksPerFrame > 2 {
		g.TicksPerFrame -= 1 // Más rápido (menos ticks de espera)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && g.TicksPerFrame < 120 {
		g.TicksPerFrame += 1 // Más lento
	}

	if g.Mode == ModeEdit {
		g.handleEditMode()
	} else {
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
				if g.Mode == ModeEdit {
					g.SelectedNodeID = id
				}
				break
			}
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.draggedNode != nil {
		g.draggedNode.X, g.draggedNode.Y = fx, fy
		g.draggedNode.VX, g.draggedNode.VY = 0, 0
	} else {
		g.draggedNode = nil
	}
}

func (g *Game) handleEditMode() {
	// [N] NUEVO ARBOL: Borra todo y deja solo la raíz
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.Tree = Tree.NewTree() // Crea un árbol fresco (solo nodo 0)
		g.Nodes = make(map[int]*VisualNode)
		g.Edges = make([]*VisualEdge, 0)
		g.SelectedNodeID = 0
		g.syncVisuals()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		if err := g.Tree.AddNode(g.SelectedNodeID); err == nil {
			g.syncVisuals()
		}
	}

	// Carga de catálogo
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.LoadTreeFromCatalog(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.LoadTreeFromCatalog(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.LoadTreeFromCatalog(2)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.LoadTreeFromCatalog(3)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.LoadTreeFromCatalog(4)
	}

	// Ejecutar Algoritmos pasando el nodo seleccionado
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		if steps, err := g.Tree.TraversalBfsSteps(g.SelectedNodeID); err == nil {
			g.TraversalSteps, g.CurrentStep, g.IsPlaying, g.Mode = steps, 0, true, ModePlayback
		} else {
			fmt.Println("Error BFS:", err)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if steps, err := g.Tree.TraversalDfsSteps(g.SelectedNodeID); err == nil {
			g.TraversalSteps, g.CurrentStep, g.IsPlaying, g.Mode = steps, 0, true, ModePlayback
		} else {
			fmt.Println("Error DFS:", err)
		}
	}
}

func (g *Game) handlePlaybackMode() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.Mode, g.IsPlaying = ModeEdit, false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.IsPlaying = !g.IsPlaying
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.CurrentStep < len(g.TraversalSteps)-1 {
		g.CurrentStep++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.CurrentStep > 0 {
		g.CurrentStep--
	}

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

func (g *Game) updatePhysics() {
	const repulsion, springLen, springK, gravity, friction, maxSpeed = 4000.0, 100.0, 0.06, 0.015, 0.82, 20.0
	for id1, n1 := range g.Nodes {
		for id2, n2 := range g.Nodes {
			if id1 == id2 {
				continue
			}
			dx, dy := n2.X-n1.X, n2.Y-n1.Y
			if math.Abs(dx) > 350 || math.Abs(dy) > 350 {
				continue
			}
			distSq := dx*dx + dy*dy + 0.1
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
	for _, n := range g.Nodes {
		if n == g.draggedNode {
			continue
		}
		n.VX += (cx - n.X) * gravity
		n.VY += (cy - n.Y) * gravity
		n.VX, n.VY = n.VX*friction, n.VY*friction
		if n.VX > maxSpeed {
			n.VX = maxSpeed
		} else if n.VX < -maxSpeed {
			n.VX = -maxSpeed
		}
		if n.VY > maxSpeed {
			n.VY = maxSpeed
		} else if n.VY < -maxSpeed {
			n.VY = -maxSpeed
		}
		n.X, n.Y = n.X+n.VX, n.Y+n.VY
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.BgColor)

	// Aristas con color modernizado (Slate Blue)
	edgeColor := color.RGBA{45, 45, 58, 255}
	for _, edge := range g.Edges {
		vector.StrokeLine(screen, float32(edge.From.X), float32(edge.From.Y), float32(edge.To.X), float32(edge.To.Y), 2, edgeColor, true)
	}

	var currentState Tree.TraversalState
	if g.Mode == ModePlayback && len(g.TraversalSteps) > 0 {
		currentState = g.TraversalSteps[g.CurrentStep].State
	}

	// Paleta de Colores Moderna
	colorUnseen := color.RGBA{84, 110, 122, 255}   // Slate Gray
	colorFrontier := color.RGBA{78, 205, 196, 255} // Neon Cyan
	colorVisited := color.RGBA{149, 117, 205, 255} // Soft Purple
	colorCurrent := color.RGBA{255, 107, 107, 255} // Coral Red

	for id, n := range g.Nodes {
		c := colorUnseen
		if currentState != nil {
			if curr := currentState.GetCurrent(); curr != nil && curr.Id == id {
				c = colorCurrent
			} else {
				isF := false
				for _, f := range currentState.GetFrontier() {
					if f.Id == id {
						isF = true
						break
					}
				}
				if isF {
					c = colorFrontier
				} else {
					for _, v := range currentState.GetVisited() {
						if v.Id == id {
							c = colorVisited
							break
						}
					}
				}
			}
		}

		// Efecto de Selección
		if g.Mode == ModeEdit && g.SelectedNodeID == id {
			vector.DrawFilledCircle(screen, float32(n.X), float32(n.Y), n.Radius+4, color.White, true)
		}

		vector.DrawFilledCircle(screen, float32(n.X), float32(n.Y), n.Radius, c, true)
		idStr := strconv.Itoa(n.ID)
		text.Draw(screen, idStr, basicfont.Face7x13, int(n.X)-(len(idStr)*7/2), int(n.Y)+4, color.Black)
	}

	g.drawUI(screen)
}

func (g *Game) drawUI(screen *ebiten.Image) {
	// Panel Superior (Hicimos el panel más alto: 70px)
	vector.DrawFilledRect(screen, 0, 0, float32(g.ScreenWidth), 70, color.RGBA{18, 18, 24, 220}, true)

	// Barra de velocidad (UI Visual)
	speedPercent := 100 - (g.TicksPerFrame * 100 / 120)
	vector.DrawFilledRect(screen, 10, 60, float32(g.ScreenWidth-20), 4, color.RGBA{45, 45, 58, 255}, true)
	vector.DrawFilledRect(screen, 10, 60, float32((g.ScreenWidth-20)*speedPercent/100), 4, color.RGBA{78, 205, 196, 255}, true)

	if g.Mode == ModeEdit {
		titulo := "=== MODO EDICION ==="
		if len(g.CatalogNames) > 0 {
			titulo = fmt.Sprintf("=== MODO EDICION | Arbol: %s ===", g.CatalogNames[g.CurrentTreeIndex])
		}
		text.Draw(screen, titulo, basicfont.Face7x13, 20, 25, color.RGBA{78, 205, 196, 255})

		// Línea de edición
		text.Draw(screen, fmt.Sprintf("Seleccionado: %d  |  [C] Agregar Hijo  |  [N] Nuevo Arbol  |  [1-5] Catalogo", g.SelectedNodeID), basicfont.Face7x13, 20, 45, color.White)

		// Línea de algoritmos (Separada a la derecha)
		text.Draw(screen, "LANZAR DESDE NODO SELECCIONADO:", basicfont.Face7x13, 650, 25, color.RGBA{84, 110, 122, 255})
		text.Draw(screen, "[B] Ver BFS   |   [D] Ver DFS", basicfont.Face7x13, 650, 45, color.RGBA{255, 107, 107, 255})

	} else if g.Mode == ModePlayback {
		text.Draw(screen, "=== MODO REPRODUCCION ===", basicfont.Face7x13, 20, 25, color.RGBA{149, 117, 205, 255})
		status := "Pausado"
		if g.IsPlaying {
			status = "Corriendo"
		}
		text.Draw(screen, fmt.Sprintf("Paso: %d/%d (%s)", g.CurrentStep, len(g.TraversalSteps)-1, status), basicfont.Face7x13, 20, 45, color.White)
		text.Draw(screen, "[Espacio] Play/Pausa  |  [<-] [->] Paso manual  |  [ESC] Salir a Edicion", basicfont.Face7x13, 400, 45, color.RGBA{255, 200, 0, 255})
	}

	// Indicador de Velocidad (Esquina superior derecha)
	text.Draw(screen, fmt.Sprintf("VELOCIDAD: %d%%", speedPercent), basicfont.Face7x13, g.ScreenWidth-150, 25, color.RGBA{78, 205, 196, 255})
	text.Draw(screen, "Teclas [ARRIBA / ABAJO]", basicfont.Face7x13, g.ScreenWidth-185, 45, color.RGBA{84, 110, 122, 255})
}

func (g *Game) Layout(w, h int) (int, int) { return g.ScreenWidth, g.ScreenHeight }
