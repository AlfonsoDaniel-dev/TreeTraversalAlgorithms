# Quick Start Guide - Tree Traversal Visualization

## Building the Project

```bash
cd /home/alfonso/Documents/uni/metodosNumericos/TreeTraversal
go build -o tree-viz ./cmd
./tree-viz
```

## Running the Application

The application will launch in **Edit Mode** with a sample tree pre-loaded:
- **7 nodes** (from arboles.json if available)
- **Dark theme** UI with modern color palette
- **Ready to explore** with keyboard and mouse

### Default Sample Tree Structure (from arbitlos.json):
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

### Edit Mode (Default)
| Key | Action |
|-----|--------|
| **Mouse Click** | Select a node (white outline) |
| **Mouse Drag** | Move selected node around |
| **C** | Add child to selected node |
| **1-5** | Load pre-configured test trees |
| **B** | Start BFS from selected node |
| **D** | Start DFS from selected node |
| **↑** | Speed up animation |
| **↓** | Slow down animation |

### Playback Mode (During algorithm visualization)
| Key | Action |
|-----|--------|
| **Space** | Play / Pause animation |
| **→** (Right Arrow) | Next traversal step |
| **←** (Left Arrow) | Previous traversal step (rewind) |
| **↑** | Speed up playback |
| **↓** | Slow down playback |
| **ESC** | Exit playback, return to edit mode |
| **Mouse Drag** | Still available - move nodes |

## Understanding the Visualization

### Node Colors During Traversal
The color of each node reflects its state during algorithm execution:

| Color | Meaning | Hex Value |
|-------|---------|-----------|
| **🔴 Coral Red** | Currently processing | #FF6B6B |
| **🟣 Soft Purple** | Already visited/explored | #9575CD |
| **🟦 Neon Cyan** | In queue/stack (frontier) | #4ECDC4 |
| **⚪ Slate Gray** | Not yet discovered | #546E7A |

**Color State Progression**:
1. Start: All nodes are **Gray** (undiscovered)
2. Frontier added: Nodes turn **Cyan** (waiting to be processed)
3. Position reached: Node turns **Red** (currently processing)
4. Processing complete: Node turns **Purple** (visited)

### Physical Representation:
- **Gray Lines**: Connections between parent and child nodes
- **Node Circles**: Each represents a tree node with ID displayed
- **Nodes Spread Out**: Due to physics simulation (repulsion/attraction)
- **Nodes Settle Down**: Gravity and friction stabilize the layout

## How It Works

### Step 1: Physics Simulation (Real-time 60 FPS)
Each frame, the system applies four physical forces:

1. **Coulomb Repulsion** (F = k/r²)
   - All nodes push each other away
   - Prevents overlap and spreads the tree
   - Constant: 4000.0

2. **Hook Attraction** (F = -k·Δx)
   - Parent-child nodes pull together
   - Maintains tree structure visibility
   - Spring length: 100px | Constant: 0.06

3. **Gravity** (F = g·m)
   - Gentle pull toward screen center
   - Prevents drift to edges
   - Constant: 0.015

4. **Friction Damping** (v *= coefficient)
   - Reduces velocity each frame
   - Stabilizes chaotic motion
   - Coefficient: 0.82 | Max speed: 20px/frame

**Result**: Organic, force-directed tree layout that settles smoothly

### Step 2: Algorithm Execution (Recorder Pattern)
When you press **B** or **D**:

1. **Algorithm runs** BFS or DFS from your selected node
2. **Each iteration**, complete state snapshot is recorded:
   - Current node being processed
   - List of visited nodes
   - Frontier (queue/stack)
   - Undiscovered nodes
   - Full path taken
3. **Steps array** returned contains full algorithm history
4. **Game loads** steps into playback mode

Code example:
```go
// Recording BFS steps
type TraversalStep struct {
    Id    int
    State TraversalState  // Interface with GetCurrent(), GetVisited(), etc.
}

type BfsState struct {
    CurrentNode *Node
    Visited     []*Node
    queue       []*Node         // Frontier
    Unvisited   []*Node
    PathTaken   []int
}
```

### Step 3: Playback & Visualization  
- **Arrow keys** let you navigate through recorded steps
- **Each frame** displays the state from current step
- **Node colors** update instantly based on their role in current step
- **Pause/Resume** with Space, rewind with Left arrow

## Customizing the Visualization

### Adding More Nodes to the Tree
Edit `cmd/main.go`:

```go
tree := Tree.NewTree()  // Creates root node 0 automatically

// Add direct children to root
tree.AddNodeFromRoot()      // Creates node 1 as child of root
tree.AddNodeFromRoot()      // Creates node 2 as child of root

// Add children to specific nodes
tree.AddNode(1)             // Creates node 3 as child of node 1
tree.AddNode(1)             // Creates node 4 as child of node 1
tree.AddNode(2)             // Creates node 5 as child of node 2
tree.AddNode(3)             // Creates node 6 as child of node 3

// Rebuild and run
// go build -o tree-viz ./cmd && ./tree-viz
```

### Using Pre-configured Test Trees
In-game (Edit Mode):
- **Press 1** → Random 150-node tree
- **Press 2** → Deep 50-level chain tree
- **Press 3-5** → Custom trees from `cmd/arboles.json` (if available)

### Starting from Different Nodes
1. **In-game**: Click a node to select it (white outline appears)
2. Press **B** for BFS or **D** for DFS starting from that node

**In code** (editing `cmd/main.go`):
```go
// Get algorithm steps starting from node 3
bfsSteps, err := tree.TraversalBfsSteps(3)
if err != nil {
    log.Fatal(err)
}

// Or DFS from node 5
dfsSteps, err := tree.TraversalDfsSteps(5)
if err != nil {
    log.Fatal(err)
}
```

### Adjusting Physics Parameters
Edit `src/Game/game.go`, in the `updatePhysics()` method:

```go
const (
    repulsion       = 4000.0  // Coulomb constant (↑ = more spread)
    springLen       = 100.0   // Hook equilibrium distance
    springK         = 0.06    // Hook constant (↑ = tighter clustering)
    gravity         = 0.015   // Center attraction (↑ = pull stronger)
    friction        = 0.82    // Velocity damping (↑ = less bouncy)
    maxSpeed        = 20.0    // Speed cap
)
```

## Example Physics Tuning

### Tight Clustering (Dense tree layout)
Edit `src/Game/game.go` line ~243 in `updatePhysics()`:
```go
const (
    repulsion = 800.0       // Reduce repulsion
    springK   = 0.15        // Increase attraction
    friction  = 0.95        // More damping
)
```

### Spread Out Layout (Spacious tree layout)
```go
const (
    repulsion = 8000.0      // Increase repulsion
    springLen = 150.0        // Larger spacing
    springK   = 0.02        // Reduce attraction
)
```

### Bouncy/Energetic Movement (Dramatic animations)
```go
const (
    friction  = 0.65        // Less damping
    gravity   = 0.001        // Minimal center pull
    maxSpeed  = 50.0        // Higher speed cap
)
```

### Mobile-friendly (Smoother, gentler)
```go
const (
    repulsion = 2000.0      // Less repulsion
    springK   = 0.1         // Standard attraction
    friction  = 0.9         // High damping
)
```

## File Overview

| File | Purpose | Lines |
|------|---------|-------|
| `src/Game/game.go` | **Main engine**: Physics simulation, rendering, input handling, UI | ~361 |
| `src/Node/node.go` | Node data structure with physics properties (position, velocity, mass) | ~61 |
| `src/Tree/tree.go` | Tree structure, node management, add/remove operations | ~97 |
| `src/Tree/bfs.go` | BFS algorithm with complete state recording (Recorder pattern) | ~111 |
| `src/Tree/dfs.go` | DFS algorithm with complete state recording (Recorder pattern) | ~105 |
| `src/Tree/TraversalState.go` | Interface for algorithm state (GetCurrent, GetVisited, GetFrontier, etc.) | ~12 |
| `src/Tree/TraversalStep.go` | Data structure holding step ID and state snapshot | ~7 |
| `cmd/main.go` | Application entry point, tree loading, catalog management | ~96 |
| `cmd/arboles.json` | JSON file with pre-configured tree structures (optional) | |

## Troubleshooting

### Problem: Nodes Disappearing Off-Screen
**Cause**: Too much repulsion vs. center gravity
- Decrease `repulsion` constant (less spread)
- Increase `gravity` constant (stronger pull to center)
- Decrease `springLen` (closer spacing)

### Problem: Nodes Not Moving/Static Layout
**Cause**: Physics frozen or friction too high
- Decrease `friction` (more bouncy)
- Increase `repulsion` (more force)
- Try dragging a node manually to test physics

### Problem: Layout Too Chaotic/Bouncy
**Cause**: Too many forces fighting
- Increase `friction` (more damping)
- Decrease `repulsion` (less aggressive spreading)
- Increase `springK` (stronger parent-child bonds)

### Problem: Traversal Steps Not Advancing
**Cause**: Algorithm didn't record steps or playback not starting
- Verify you **selected a node** before pressing B/D
- Check console for error messages
- Try a different tree (press 1-5)

### Problem: Colors Not Updating
**Cause**: Step playback not active
- Verify you're in **Playback Mode** (top status bar shows "REPRODUCCION")
- Try pressing Space to play the animation
- Check that current step < total steps

## Architecture Overview

### System Architecture
```
┌────────────────────────────────────────────┐
│      Ebiten Game Loop (60 FPS Target)      │
├────────────────────────────────────────────┤
│                                            │
│  Update()                          Draw()  │
│  ├─ Handle Input                  ├─ Fill background
│  │  ├─ Mouse (drag/select)        ├─ Draw edges (gray)
│  │  ├─ Keyboard (B/D/arrows)      ├─ Categorize nodes by state
│  │  └─ Mode management (Edit/    ├─ Draw nodes (colored circles)
│  │     Playback)                  ├─ Draw node IDs
│  │                                └─ Draw UI panel
│  ├─ Physics Simulation (if not              
│  │   dragging)                    
│  │  ├─ Coulomb repulsion (all pairs)
│  │  ├─ Hook attraction (parent-child)
│  │  ├─ Gravity (center pull)
│  │  ├─ Friction (velocity damping)
│  │  └─ Update positions
│  │
│  └─ Playback Logic (if playing)
│     └─ Advance step counter
│
└────────────────────────────────────────────┘
                    ↑ Uses ↑
        ┌───────────────────────────┐
        │    Tree Data Structure    │
        ├───────────────────────────┤
        │ ├─ Root: Node 0          │
        │ ├─ Nodes: map[id]Node    │
        │ └─ LastNodeId: int       │
        └───────────────────────────┘
                    ↑ Processes ↑
        ┌───────────────────────────┐
        │  BFS/DFS Algorithms       │
        ├───────────────────────────┤
        │ Records complete state    │
        │ at EACH iteration:        │
        │ ├─ Current node          │
        │ ├─ Visited list          │
        │ ├─ Frontier (queue/stack)│
        │ ├─ Unseen nodes          │
        │ └─ Path taken            │
        │                           │
        │ Returns: []TraversalStep │
        └───────────────────────────┘
```

### Recorder Pattern Flow
```
[1] User presses B/D
        ↓
[2] Algorithm.TraversalBfsSteps(nodeId)
        ├─ Runs algorithm iteration
        ├─ Creates snapshot state
        └─ Appends to history
        ↓
[3] Returns []TraversalStep
        ├─ Id: 0, State: BfsState{...}
        ├─ Id: 1, State: BfsState{...}
        └─ Id: N, State: BfsState{...}
        ↓
[4] Game enters Playback Mode
        ├─ CurrentStep = 0
        └─ IsPlaying = true
        ↓
[5] Each frame:
        ├─ Get state = TraversalSteps[CurrentStep].State
        └─ Render nodes with colors from state
        ↓
[6] User presses → or Space progresses
        └─ Increment CurrentStep
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

### 1. Recorder Pattern (Core Feature)
- **Decouples** algorithm execution from visualization
- Algorithm records complete state at EVERY iteration
- Each TraversalStep = snapshot in time
- Enables: pause, rewind, playback at any speed
- Alternative: Simulate real-time (would run algorithm every frame)

### 2. State Machine (Game Modes)
- **Edit Mode**: Build trees, select nodes, configure algorithms
- **Playback Mode**: Visualize recorded algorithm steps
- Can toggle between modes (ESC exits playback)

### 3. Physics Engine (Independent System)
- **Runs independently** of algorithm playback
- Positions change smoothly even while viewing single step
- Four forces create force-directed layout:
  - Repulsion (spread) + Attraction (hierarchy) + Gravity (center) + Friction (stability)
- Natural, organic appearance

### 4. Color Coding (Visual Feedback)
- Immediate understanding of algorithm state
- Four distinct colors for four node roles
- Updates instant as you navigate steps

### 5. Force-Directed Graph (Layout Algorithm)
- Better than fixed/hierarchical layouts
- Natural spacing without manual positioning
- Parent-child relationships visible through physics

### 6. Interface-Based Design (TraversalState)
- Both BFS and DFS implement same interface
- Easy to add new algorithms (implement GetCurrent, GetVisited, etc.)
- Game code doesn't need to know specific algorithm details

---

## Next Steps

1. **Try the defaults** - press B or D to see algorithms in action
2. **Load test trees** - press 1-5 to see different structures
3. **Adjust physics** - edit constants in `game.go` to customize appearance
4. **Create your own tree** - edit `cmd/main.go` with AddNode calls
5. **Read IMPLEMENTATION.md** - for deep technical details
6. **Check TECHNICAL_REVIEW.txt** - for complete code samples

**Enjoy exploring tree traversals visually!** 🌳✨

