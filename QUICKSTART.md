# Quick Start Guide - Tree Traversal Visualization

## Building the Project

```bash
cd /home/alfonso/Documents/uni/metodosNumericos/TreeTraversal
go build -o tree-visualization ./cmd
./tree-visualization
```

## Running the Application

The application will launch with a sample tree pre-loaded:
- **7 nodes** structured in a tree hierarchy
- **BFS traversal** visualization enabled by default

### Default Sample Tree Structure:
```
        0 (Root)
       / \
      1   2
     / \   \
    3   4   5
   / \
  6   7
```

## Interactive Controls

### Navigation:
| Key | Action |
|-----|--------|
| **→** (Right Arrow) | Next traversal step |
| **←** (Left Arrow) | Previous traversal step |
| **Space** | Pause/Resume (prepared for future implementation) |

## Understanding the Visualization

### Node Colors During BFS/DFS:
1. **Red Nodes**: Currently being processed
2. **Blue Nodes**: Already visited/explored
3. **Yellow Nodes**: In queue/stack waiting to be processed
4. **Gray Nodes**: Not yet discovered

### Physical Representation:
- **Gray Lines**: Connections between parent and child nodes
- **Node Circles**: Each represents a tree node
- **Nodes Spread Out**: Due to physics simulation (repulsion/attraction)
- **Nodes Settle Down**: Gravity and friction stabilize the layout

## How It Works

### 1. Physics Simulation (Real-time)
Each frame, the system applies:
- **Repulsion** between all nodes (Coulomb force)
- **Attraction** between parent-child nodes (Hook's law)
- **Gravity** pulling nodes downward
- **Friction** damping movement
- **Boundary constraints** keeping nodes on-screen

### 2. Traversal Playback (Step-by-step)
Use arrow keys to navigate through algorithm steps:
- Each step shows current node state
- Colors update to reflect algorithm progress
- Nodes remain in their physics-simulated positions

## Customizing the Visualization

### Adding More Nodes to the Tree
Edit `cmd/main.go`:

```go
tree := Tree.NewTree()
tree.AddNodeFromRoot()      // Add node 1 to root
tree.AddNodeFromRoot()      // Add node 2 to root
tree.AddNode(1)             // Add node 3 to node 1
tree.AddNode(1)             // Add node 4 to node 1
// Add more as needed...
```

### Switching to DFS
Edit `cmd/main.go`, replace:
```go
bfsSteps, err := tree.TraversalBfsSteps()
```
with:
```go
bfsSteps, err := tree.TraversalDfsSteps()
```

### Adjusting Physics Parameters
Edit `src/Game/game.go`, in the `NewGame` function:

```go
g.Physics = PhysicsEngine{
    CoulombConstant:       500.0,   // ↑ = more repulsion
    HookConstant:          0.1,     // ↑ = tighter clustering
    EquilibriumDistance:   100.0,   // ↑ = nodes spread farther
    GravityForce:          100.0,   // ↑ = stronger downward pull
    Friction:              0.85,    // ↓ = more bouncy, ↑ = more damped
    DeltaTime:             1.0/60.0, // Keep as is for 60 FPS
}
```

## Example Physics Tuning

### For Tightly Clustered Layout:
```go
CoulombConstant:     200.0,   // Reduce repulsion
HookConstant:        0.3,     // Increase attraction
Friction:            0.95,    // More damping
```

### For Spread Out Layout:
```go
CoulombConstant:     1000.0,  // Increase repulsion
HookConstant:        0.05,    // Reduce attraction
EquilibriumDistance: 150.0,   // Larger spacing
```

### For Bouncy/Energetic Movement:
```go
Friction: 0.7,                // Less damping
GravityForce: 50.0,           // Less gravity
```

## File Overview

| File | Purpose |
|------|---------|
| `src/Game/game.go` | Main game loop, physics engine, rendering |
| `src/Node/node.go` | Node data structure with physics properties |
| `src/Tree/tree.go` | Tree data structure and management |
| `src/Tree/bfs.go` | BFS algorithm with step recording |
| `src/Tree/dfs.go` | DFS algorithm with step recording |
| `cmd/main.go` | Application entry point |

## Troubleshooting

### Nodes Disappearing:
- Decrease `CoulombConstant` (too much repulsion)
- Increase `EquilibriumDistance` (nodes too far apart)

### Nodes Not Moving:
- Increase `GravityForce`
- Decrease `Friction`
- Increase `CoulombConstant`

### Layout Too Chaotic:
- Increase `Friction` (more damping)
- Increase `HookConstant` (stronger parent-child bonds)
- Decrease `CoulombConstant` (less repulsion)

### Traversal Not Showing:
- Press Right Arrow to advance steps
- Verify `game.TraversalSteps` is populated in main.go

## Architecture Overview

```
┌─────────────────────────────────┐
│     Ebiten Game Loop (60 FPS)   │
├─────────────────────────────────┤
│                                 │
│  Update():                      │  Draw():
│  ├─ Physics Simulation          │  ├─ Categorize Nodes by State
│  │  ├─ Coulomb Repulsion        │  ├─ Draw Edges
│  │  ├─ Hook Attraction          │  └─ Draw Nodes
│  │  ├─ Gravity                  │
│  │  ├─ Friction                 │
│  │  └─ Boundary Enforcing       │
│  └─ Input Processing            │
│                                 │
└─────────────────────────────────┘
         ↓ Uses ↓
┌─────────────────────────────────┐
│      Tree Structure             │
│                                 │
│  ├─ Root Nodes                  │
│  ├─ Child Relationships         │
│  └─ Node Map (by ID)            │
└─────────────────────────────────┘
         ↓ Processes ↓
┌─────────────────────────────────┐
│   BFS/DFS Traversal Algorithms  │
│                                 │
│  ├─ Records Steps               │
│  ├─ Tracks State Changes        │
│  └─ Generates TraversalSteps    │
└─────────────────────────────────┘
```

## Performance Notes

- **Suitable for**: Trees with < 1000 nodes
- **Smooth at**: ~60 FPS on most systems
- **Physics iterations**: 1 per frame

For trees with > 1000 nodes, consider:
1. Reducing `CoulombConstant` (fewer collision calculations)
2. Using spatial partitioning
3. Simplifying rendering (reduce circle resolution)

## Advanced: Custom Traversal Colorization

To modify state colors, edit the `drawNodes` method in `game.go`:

```go
// Current state color
g.drawNode(screen, node, color.RGBA{R: 255, G: 50, B: 50, A: 255}, nodeRadius)

// Change to any color: color.RGBA{R: value, G: value, B: value, A: 255}
```

## Key Concepts

1. **State Machine**: Nodes transition between states during traversal
2. **Physics Engine**: Independent force-based layout system
3. **Recorder Pattern**: Algorithm records steps, game replays visualization
4. **Color Coding**: Immediate visual feedback on algorithm progress
5. **Force-Directed Graph**: Natural spacing based on physics laws

Enjoy your tree traversal visualization!

