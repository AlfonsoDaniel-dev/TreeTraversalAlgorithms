package main

import (
	"encoding/json"
	"image/color"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	// IMPORTANTE: Ajusta esta ruta a tu módulo
	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Game"
	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Tree"
)

type ArbolJSON struct {
	Nombre string `json:"nombre"`
	Padres []int  `json:"padres"`
}

func GenerarArbolAleatorio(numNodos int) *Tree.Tree {
	t := Tree.NewTree()
	for i := 1; i < numNodos; i++ {
		t.AddNode(rand.Intn(i))
	}
	return t
}

func GenerarArbolProfundo(niveles int) *Tree.Tree {
	t := Tree.NewTree()
	for i := 0; i < niveles; i++ {
		t.AddNode(i)
	}
	return t
}

func main() {
	ebiten.SetWindowTitle("Visualizador BFS/DFS - Ebitengine")
	ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	var catalogTrees []*Tree.Tree
	var catalogNames []string

	// 1. Leer JSON
	data, err := os.ReadFile("./cmd/arboles.json")
	if err == nil {
		var catalogo []ArbolJSON
		if err := json.Unmarshal(data, &catalogo); err == nil {
			for _, arbolData := range catalogo {
				t := Tree.NewTree()
				for _, padreID := range arbolData.Padres {
					t.AddNode(padreID)
				}
				catalogTrees = append(catalogTrees, t)
				catalogNames = append(catalogNames, arbolData.Nombre)
			}
		} else {
			log.Println("Error parseando JSON:", err)
		}
	} else {
		log.Println("No se encontró arboles.json, continuando con generadores locales.")
	}

	// 2. Generar Árboles Masivos por Código
	catalogTrees = append(catalogTrees, GenerarArbolAleatorio(150))
	catalogNames = append(catalogNames, "Monstruo Aleatorio (150 Nodos)")

	catalogTrees = append(catalogTrees, GenerarArbolProfundo(50))
	catalogNames = append(catalogNames, "Látigo Profundo (50 Nodos)")

	// Fallback de seguridad por si no hay árboles
	if len(catalogTrees) == 0 {
		catalogTrees = append(catalogTrees, Tree.NewTree())
		catalogNames = append(catalogNames, "Árbol Vacío")
	}

	// 3. Inicializar el Juego
	bgColor := color.RGBA{R: 15, G: 15, B: 25, A: 255}
	game := Game.NewGame(bgColor, 1280, 720, catalogTrees[0])

	game.CatalogTrees = catalogTrees
	game.CatalogNames = catalogNames
	game.CurrentTreeIndex = 0

	log.Println("Sistema iniciado correctamente. Catálogo cargado:")
	for i, name := range catalogNames {
		log.Printf(" [%d] %s", i+1, name)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Error al ejecutar el juego:", err)
	}
}
