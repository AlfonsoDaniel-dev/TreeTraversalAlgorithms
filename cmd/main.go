package main

import (
	"image/color"
	"log"

	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Game"
	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Tree"
	"github.com/hajimehoshi/ebiten/v2"
	// IMPORTANTE: Cambia "tu_proyecto" por el nombre real de tu módulo
)

func main() {
	// 1. Configuración de la Ventana
	ebiten.SetWindowTitle("Visualizador BFS/DFS - Motor de Físicas")
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// 2. Creación de la Estructura Lógica (El Árbol)
	// Tu estructura inicializa el Root con ID 0 automáticamente
	tree := Tree.NewTree()

	// Nivel 1 (Hijos de la Raíz 0)
	tree.AddNodeFromRoot() // Crea ID 1
	tree.AddNodeFromRoot() // Crea ID 2
	tree.AddNodeFromRoot() // Crea ID 3

	// Nivel 2
	tree.AddNode(1) // Crea ID 4 (Hijo de 1)
	tree.AddNode(1) // Crea ID 5 (Hijo de 1)
	tree.AddNode(2) // Crea ID 6 (Hijo de 2)
	tree.AddNode(2) // Crea ID 7 (Hijo de 2)
	tree.AddNode(3) // Crea ID 8 (Hijo de 3)

	// Nivel 3
	tree.AddNode(4) // Crea ID 9 (Hijo de 4)
	tree.AddNode(4) // Crea ID 10 (Hijo de 4)
	tree.AddNode(7) // Crea ID 11 (Hijo de 7)
	tree.AddNode(7) // Crea ID 12 (Hijo de 7)

	bfsSteps, err := tree.TraversalBfsSteps()
	if err != nil {
		log.Fatalf("Error al calcular BFS: %v", err)
	}

	bgColor := color.RGBA{R: 20, G: 22, B: 30, A: 255} // Un tono oscuro elegante
	game := Game.NewGame(bgColor, 1280, 720, tree)

	game.TraversalSteps = bfsSteps

	game.IsPlaying = false
	game.CurrentStep = 0

	log.Println("Iniciando simulación. Controles:")
	log.Println(" - Ratón: Arrastrar nodos")
	log.Println(" - Espacio: Reproducir / Pausar animación BFS")
	log.Println(" - Flechas Izq/Der: Avanzar paso a paso manualmente")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Error al ejecutar el juego: ", err)
	}
}
