package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	// Asegúrate de que estas rutas coincidan con tu go.mod
	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Game"
	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Tree"
)

func main() {
	// 1. Configuración de la Ventana para el maestro
	ebiten.SetWindowTitle("Visualizador de Algoritmos - Estructuras de Datos")
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// 2. Construcción Controlada del Árbol
	tree := Tree.NewTree()

	// 3. Preparación de la Visualización
	// Por defecto calculamos el BFS para la primera impresión
	bfsSteps, err := tree.TraversalBfsSteps()
	if err != nil {
		log.Fatal("No se pudieron pre-calcular los pasos del BFS:", err)
	}

	// 4. Inicialización del Motor Visual
	// Usamos un fondo oscuro azulado (Dark Mode) que es más profesional
	bgColor := color.RGBA{R: 15, G: 15, B: 25, A: 255}
	game := Game.NewGame(bgColor, 1280, 720, tree)

	// Cargamos los pasos en el reproductor
	game.TraversalSteps = bfsSteps
	game.Mode = 0 // Forzamos Modo Edición para que el maestro pueda jugar primero

	log.Println("Sistema iniciado correctamente.")
	log.Println("Controles para el usuario:")
	log.Println(" - [C]: Agregar hijo al nodo seleccionado")
	log.Println(" - [B]: Iniciar simulación BFS")
	log.Println(" - [D]: Iniciar simulación DFS")
	log.Println(" - [Espacio]: Pausar/Reproducir")
	log.Println(" - [ESC]: Volver al modo edición")

	// 5. ¡Ejecución!
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("El motor gráfico se detuvo:", err)
	}
}
