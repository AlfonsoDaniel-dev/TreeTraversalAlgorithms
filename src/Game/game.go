package Game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"strconv"

	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Node"
	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"

	// Ajusta esta ruta a tu módulo
	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Tree"
)

const sampleRate = 44100

var audioCtx *audio.Context

func init() {
	// Inicializa el hardware de sonido de la computadora a 44.1 kHz
	audioCtx = audio.NewContext(sampleRate)
}

// generaBlip crea una onda de sonido corta (50 milisegundos) a 880Hz (Nota La/A5)
func generaBlip(freq float64) []byte {
	const duration = 0.05 // 50 milisegundos
	length := int(sampleRate * duration)

	b := make([]byte, length*4)

	for i := 0; i < length; i++ {
		volumen := math.Sin(2 * math.Pi * freq * float64(i) / sampleRate)
		v16 := int16(volumen * (math.MaxInt16 / 6)) // /6 para un volumen agradable

		b[4*i] = byte(v16)
		b[4*i+1] = byte(v16 >> 8)
		b[4*i+2] = byte(v16)
		b[4*i+3] = byte(v16 >> 8)
	}
	return b
}

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
	beepPlayers map[int]*audio.Player
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
		TicksPerFrame:  30,
	}

	// Cargamos el sonido al reproductor en la memoria de la tarjeta de sonido
	g.beepPlayers = make(map[int]*audio.Player)

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

	for id := range g.Nodes {
		if _, exists := g.Tree.Nodes[id]; !exists {
			delete(g.Nodes, id)
		}
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

// playBeepForNode calcula la profundidad y toca una nota musical
func (g *Game) playBeepForNode(currNode interface{}) {
	// 1. Convertimos la interfaz al tipo correcto (Node.Node de tu Tree)
	// Ajusta "Tree.Node" si tu struct se llama diferente (ej: *Tree.Node)
	node, ok := currNode.(*Node.Node)
	if !ok || node == nil {
		return
	}

	// 2. Calcular la profundidad (viajando hacia el padre)
	depth := 0
	padre := node.GetParent()
	for padre != nil {
		depth++
		padre = padre.GetParent()
	}

	// 3. Generar el sonido solo si no existe en caché
	if _, exists := g.beepPlayers[depth]; !exists {
		// FÓRMULA MUSICAL: Escala de Tonos Enteros (Whole Tone Scale)
		// Base: 440Hz (Nota La/A4). Sube un tono completo por cada nivel.
		freq := 440.0 * math.Pow(2.0, float64(depth*2)/12.0)
		g.beepPlayers[depth] = audioCtx.NewPlayerFromBytes(generaBlip(freq))
	}

	// 4. Disparar el sonido
	g.beepPlayers[depth].Rewind()
	g.beepPlayers[depth].Play()
}

// drawTelemetry calcula y dibuja las estadísticas matemáticas del árbol en tiempo real
func (g *Game) drawTelemetry(screen *ebiten.Image) {
	if g.Tree == nil || len(g.Nodes) == 0 {
		return
	}

	totalNodes := len(g.Nodes)
	totalEdges := len(g.Edges)
	leafNodes := 0
	nonLeafNodes := 0
	maxDepth := 0

	// 1. Cálculos Matemáticos en tiempo real
	for _, n := range g.Tree.Nodes {
		hijos := n.GetChildren()
		if len(hijos) == 0 {
			leafNodes++ // Si no tiene hijos, es una Hoja
		} else {
			nonLeafNodes++ // Si tiene hijos, es un nodo interno
		}

		// Calcular la profundidad de este nodo subiendo hasta la raíz
		profundidad := 0
		padre := n.GetParent()
		for padre != nil {
			profundidad++
			padre = padre.GetParent()
		}
		if profundidad > maxDepth {
			maxDepth = profundidad
		}
	}

	// Factor de ramificación (¿Cuántos hijos tiene un nodo en promedio?)
	avgBranching := 0.0
	if nonLeafNodes > 0 {
		avgBranching = float64(totalEdges) / float64(nonLeafNodes)
	}

	// 2. Medidas y Posición del Panel (Esquina Inferior Derecha)
	panelW, panelH := float32(260), float32(130)
	panelX := float32(g.ScreenWidth) - panelW - 20
	panelY := float32(g.ScreenHeight) - panelH - 20

	// 3. Dibujar el Panel Translúcido
	// Fondo principal
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{18, 18, 24, 220}, true)
	// Cabecera decorativa
	vector.DrawFilledRect(screen, panelX, panelY, panelW, 25, color.RGBA{45, 45, 58, 255}, true)

	// 4. Imprimir la información
	xText := int(panelX) + 15
	yText := int(panelY) + 17

	// Título
	text.Draw(screen, "INFO DEL GRAFO", basicfont.Face7x13, xText, yText, color.RGBA{78, 205, 196, 255})

	// Datos
	yText += 25
	text.Draw(screen, fmt.Sprintf("Total Nodos: %d", totalNodes), basicfont.Face7x13, xText, yText, color.White)

	yText += 20
	text.Draw(screen, fmt.Sprintf("Total Aristas: %d", totalEdges), basicfont.Face7x13, xText, yText, color.White)

	yText += 20
	text.Draw(screen, fmt.Sprintf("Nodos Hoja: %d", leafNodes), basicfont.Face7x13, xText, yText, color.RGBA{255, 107, 107, 255})

	yText += 20
	text.Draw(screen, fmt.Sprintf("Distancia max hasta la raiz: %d", maxDepth), basicfont.Face7x13, xText, yText, color.RGBA{149, 117, 205, 255})

	yText += 20
	text.Draw(screen, fmt.Sprintf("Factor Ramificacion: %.2f", avgBranching), basicfont.Face7x13, xText, yText, color.White)
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
	if inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if g.SelectedNodeID != 0 {
			g.Tree.RemoveNode(g.SelectedNodeID)
			g.SelectedNodeID = 0 // Regresamos la selección a la raíz por seguridad
			g.syncVisuals()
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

	// Paso Manual -> Adelante
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.CurrentStep < len(g.TraversalSteps)-1 {
		g.CurrentStep++
		// Extraemos el nodo actual del estado y lo sonificamos
		currNode := g.TraversalSteps[g.CurrentStep].State.GetCurrent()
		g.playBeepForNode(currNode)
	}

	// Paso Manual -> Atrás
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.CurrentStep > 0 {
		g.CurrentStep--
		currNode := g.TraversalSteps[g.CurrentStep].State.GetCurrent()
		g.playBeepForNode(currNode)
	}

	// Avance Automático (Play)
	if g.IsPlaying && len(g.TraversalSteps) > 0 {
		g.tickCounter++
		if g.tickCounter >= g.TicksPerFrame {
			g.tickCounter = 0
			if g.CurrentStep < len(g.TraversalSteps)-1 {
				g.CurrentStep++

				// ¡Sonido Dinámico!
				currNode := g.TraversalSteps[g.CurrentStep].State.GetCurrent()
				g.playBeepForNode(currNode)

			} else {
				g.IsPlaying = false
			}
		}
	}
}

/* OLD Physics engine do not uncomment
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
} */

func (g *Game) updatePhysics() {
	// Constantes calibradas para el árbol colgante
	const repulsion = 4000.0
	const springLen = 80.0 // Resortes más cortos para jerarquía
	const springK = 0.15   // Mayor tensión para cargar el peso
	const gravityY = 0.3   // Gravedad de cascada (+Y hacia abajo)
	const friction = 0.80
	const maxSpeed = 25.0

	cx := float64(g.ScreenWidth / 2)

	// 1. Repulsión Asimétrica
	for id1, n1 := range g.Nodes {
		for id2, n2 := range g.Nodes {
			if id1 == id2 {
				continue
			}

			dx, dy := n2.X-n1.X, n2.Y-n1.Y
			if math.Abs(dx) > 350 || math.Abs(dy) > 350 {
				continue
			}

			if math.Abs(dx) < 0.1 && math.Abs(dy) < 0.1 {
				dx, dy = rand.Float64()*2-1, rand.Float64()*2-1
			}
			distSq := dx*dx + dy*dy
			if distSq < 10.0 {
				distSq = 10.0
			}
			f := repulsion / distSq
			dist := math.Sqrt(distSq)

			// Repulsión Fuerte en X (Abre el árbol) y Débil en Y (Mantiene niveles)
			n1.VX -= (dx / dist) * f * 1.5
			n1.VY -= (dy / dist) * f * 0.2
		}
	}

	// 2. Atracción de Resortes
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

	// 3. Gravedad y Anclaje
	padding := float64(30.0)
	for _, n := range g.Nodes {
		if n == g.draggedNode {
			continue
		}

		// EL CLAVO: Anclaje estricto de la raíz (Nodo 0)
		if n.ID == 0 {
			n.X += (cx - n.X) * 0.1
			n.Y += (100.0 - n.Y) * 0.1
			n.VX *= 0.5
			n.VY *= 0.5
		} else {
			// Gravedad direccional y atracción leve al centro horizontal
			n.VY += gravityY
			n.VX += (cx - n.X) * 0.005
		}

		n.VX *= friction
		n.VY *= friction

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

		n.X += n.VX
		n.Y += n.VY

		if math.IsNaN(n.X) || math.IsNaN(n.Y) {
			n.X, n.Y, n.VX, n.VY = cx, 150, 0, 0
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
	// Panel Superior Transparente
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

		// Instrucciones actualizadas con [X] Borrar
		text.Draw(screen, fmt.Sprintf("Sel: %d | [C] Hijo | [X] Borrar | [N] Limpiar | [1-5] Catalogo", g.SelectedNodeID), basicfont.Face7x13, 20, 45, color.White)

		// Menú de Algoritmos (Dinámico a la derecha)
		menuDerechaX := g.ScreenWidth - 600
		text.Draw(screen, "LANZAR DESDE NODO SELECCIONADO:", basicfont.Face7x13, menuDerechaX, 25, color.RGBA{84, 110, 122, 255})
		text.Draw(screen, "[B] Ver BFS   |   [D] Ver DFS", basicfont.Face7x13, menuDerechaX, 45, color.RGBA{255, 107, 107, 255})

	} else if g.Mode == ModePlayback {
		text.Draw(screen, "=== MODO REPRODUCCION ===", basicfont.Face7x13, 20, 25, color.RGBA{149, 117, 205, 255})
		status := "Pausado"
		if g.IsPlaying {
			status = "Corriendo"
		}
		text.Draw(screen, fmt.Sprintf("Paso: %d/%d (%s)", g.CurrentStep, len(g.TraversalSteps)-1, status), basicfont.Face7x13, 20, 45, color.White)
		text.Draw(screen, "[Espacio] Play/Pausa | [<-] [->] Manual | [ESC] Salir a Edicion", basicfont.Face7x13, 350, 45, color.RGBA{255, 200, 0, 255})

		//VISUALIZADOR DE MEMORIA (Esquina inferior izquierda) ---
		if len(g.TraversalSteps) > 0 && g.CurrentStep < len(g.TraversalSteps) {
			state := g.TraversalSteps[g.CurrentStep].State
			currNode := state.GetCurrent()

			// 1. Mostrar Rastro de Movimiento
			if currNode != nil {
				textoMov := fmt.Sprintf("MOVIMIENTO: Nodo %d (Punto de Inicio)", currNode.Id)
				if currNode.Parent != nil {
					textoMov = fmt.Sprintf("MOVIMIENTO: De %d  ->  A %d", currNode.Parent.Id, currNode.Id)
				}
				text.Draw(screen, textoMov, basicfont.Face7x13, 20, g.ScreenHeight-60, color.RGBA{255, 107, 107, 255})
			}

			// 2. Mostrar Estructura (Frontera)
			text.Draw(screen, "Queue/Pila", basicfont.Face7x13, 20, g.ScreenHeight-40, color.RGBA{78, 205, 196, 255})
			strEstructura := "[ "
			for _, fNode := range state.GetFrontier() {
				strEstructura += strconv.Itoa(fNode.Id) + " "
			}
			strEstructura += "]"
			text.Draw(screen, strEstructura, basicfont.Face7x13, 20, g.ScreenHeight-20, color.White)

			// 4. MOSTRAR RUTA COMPLETA (Historial de pasos)
			text.Draw(screen, "RUTA RECORRIDA:", basicfont.Face7x13, 20, g.ScreenHeight-100, color.RGBA{149, 117, 205, 255})

			rutaStr := ""
			for i, nodeID := range state.GetPathTaken() {
				rutaStr += strconv.Itoa(nodeID)
				if i < len(state.GetPathTaken())-1 {
					rutaStr += " -> "
				}
			}

			text.Draw(screen, rutaStr, basicfont.Face7x13, 20, g.ScreenHeight-80, color.White)
		}
	}

	// Indicador de Velocidad (Fijo esquina superior derecha)
	text.Draw(screen, fmt.Sprintf("VELOCIDAD: %d%%", speedPercent), basicfont.Face7x13, g.ScreenWidth-150, 25, color.RGBA{78, 205, 196, 255})
	text.Draw(screen, "Teclas [ARRIBA / ABAJO]", basicfont.Face7x13, g.ScreenWidth-185, 45, color.RGBA{84, 110, 122, 255})
	g.drawTelemetry(screen)
}

func (g *Game) Layout(w, h int) (int, int) { return g.ScreenWidth, g.ScreenHeight }
