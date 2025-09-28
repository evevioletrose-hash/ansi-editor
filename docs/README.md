# ANSI Editor Framework Documentation

## Overview

The ANSI Editor Framework is a Go-based toolkit for creating interactive ANSI art applications that can be exported as WebAssembly (WASM) bundles. It provides a complete system for drawing, animating, and exporting ANSI artwork for use in web applications, games, and interactive interfaces.

## Features

- **Canvas System**: 2D grid-based drawing with full ANSI color support
- **Object Management**: Create and manage moveable ANSI art objects
- **Animation Support**: Frame-based animation system
- **ANSI File I/O**: Load and save standard ANSI art files
- **WASM Export**: Export projects as complete web-ready WASM bundles
- **Interactive Editor**: Command-line tool for creating and editing projects

## Quick Start

### Installation

```bash
go get github.com/evevioletrose-hash/ansi-editor
```

### Basic Usage

```go
package main

import (
    "github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
    "github.com/evevioletrose-hash/ansi-editor/pkg/objects"
    "github.com/evevioletrose-hash/ansi-editor/pkg/export"
)

func main() {
    // Create a new scene
    scene := objects.NewScene("my-scene", 80, 25)
    
    // Draw on the background
    scene.Background.DrawString(10, 5, "Hello ANSI World!", canvas.Green, canvas.Black)
    
    // Create a moveable object
    obj := objects.NewObject("player", "Player Sprite", 20, 10, 5, 3)
    obj.Canvas.DrawString(1, 1, "@", canvas.Yellow, canvas.Black)
    scene.AddObject(obj)
    
    // Export as WASM bundle
    export.ExportScene(scene, "./output", "my-app", "1.0.0")
}
```

### Using the Command-Line Editor

```bash
# Build and run the editor
go run cmd/editor/main.go

# Or build it first
go build -o ansi-editor cmd/editor/main.go
./ansi-editor
```

## Core Components

### Canvas System

The canvas system provides a 2D grid for drawing ANSI characters with colors and attributes.

#### Creating a Canvas

```go
// Create a new 80x25 canvas
canvas := canvas.NewCanvas(80, 25)

// Draw individual characters
canvas.DrawChar(10, 5, '@', canvas.Red, canvas.Black)

// Draw strings
canvas.DrawString(0, 0, "Hello World", canvas.Green, canvas.Black)

// Draw shapes
canvas.DrawRect(5, 5, 20, 10, '#', canvas.Blue, canvas.Black)
canvas.DrawLine(0, 0, 79, 24, '*', canvas.Yellow, canvas.Black)
```

#### Canvas Methods

- `NewCanvas(width, height int)` - Create new canvas
- `GetCell(x, y int)` - Get cell at coordinates
- `SetCell(x, y int, cell Cell)` - Set cell at coordinates
- `DrawChar(x, y int, char rune, fg, bg ANSIColor)` - Draw single character
- `DrawString(x, y int, text string, fg, bg ANSIColor)` - Draw text string
- `DrawRect(x, y, w, h int, char rune, fg, bg ANSIColor)` - Draw rectangle outline
- `FillRect(x, y, w, h int, char rune, fg, bg ANSIColor)` - Draw filled rectangle
- `DrawLine(x1, y1, x2, y2 int, char rune, fg, bg ANSIColor)` - Draw line
- `Clear(char rune, fg, bg ANSIColor)` - Clear entire canvas
- `ToANSI()` - Convert to ANSI escape sequences

### Object System

Objects are moveable ANSI art elements that can be positioned, animated, and layered.

#### Creating Objects

```go
// Create a new object
obj := objects.NewObject("player", "Player Character", 10, 5, 3, 3)

// Draw on the object's canvas
obj.Canvas.DrawChar(1, 1, '@', canvas.Yellow, canvas.Black)

// Move the object
obj.Move(15, 8)

// Add to scene
scene.AddObject(obj)
```

#### Animation

```go
// Create animation frames
frame1 := canvas.NewCanvas(3, 3)
frame1.DrawChar(1, 1, 'o', canvas.Red, canvas.Black)

frame2 := canvas.NewCanvas(3, 3)
frame2.DrawChar(1, 1, 'O', canvas.BrightRed, canvas.Black)

// Add frames to object
obj.AddFrame(frame1)
obj.AddFrame(frame2)

// Update animations (call this in your game loop)
scene.UpdateAnimations()
```

### Scene Management

Scenes combine a background canvas with multiple objects and handle rendering.

```go
// Create scene
scene := objects.NewScene("level1", 80, 25)

// Draw background
scene.Background.FillRect(0, 0, 80, 25, '.', canvas.Green, canvas.Black)

// Add objects
player := objects.NewObject("player", "Player", 40, 12, 1, 1)
scene.AddObject(player)

// Render complete scene
rendered := scene.Render()
```

### ANSI File I/O

Load and save standard ANSI art files.

```go
import "github.com/evevioletrose-hash/ansi-editor/pkg/ansi"

// Load ANSI file
file, _ := os.Open("artwork.ans")
parser := ansi.NewParser()
canvas, _ := parser.ParseFile(file, 80, 25)

// Save ANSI file
file, _ := os.Create("output.ans")
ansi.SaveANSI(canvas, file)
```

### WASM Export

Export your scenes as complete web applications.

```go
import "github.com/evevioletrose-hash/ansi-editor/pkg/export"

// Export single scene
export.ExportScene(scene, "./web-output", "my-app", "1.0.0")

// Or create a bundle with multiple scenes
bundle := export.NewWASMBundle("my-game", "1.0.0")
bundle.AddScene(scene1)
bundle.AddScene(scene2)
bundle.Export("./web-output")
```

The exported bundle includes:
- `index.html` - Main web page
- `ansi-runtime.js` - JavaScript runtime
- `main.go` - Go WASM source
- `build.sh` / `build.bat` - Build scripts
- `bundle.json` - Scene data

## Command-Line Editor

The included editor provides an interactive way to create ANSI art:

### Basic Commands

```
new 80 25                    # Create new 80x25 canvas
draw 10 5 @ red black        # Draw @ character at (10,5)
text 0 0 "Hello" green black # Draw text
rect 5 5 20 10 # blue black  # Draw rectangle
line 0 0 79 24 * yellow black # Draw line
show                         # Display current canvas
clear                        # Clear canvas
```

### Object Commands

```
object add player Player 10 5 3 3  # Add new object
object move player 15 8            # Move object
object list                        # List all objects
```

### File Operations

```
load artwork.ans    # Load ANSI file
save output.ans     # Save current canvas
```

### Export

```
export wasm ./output my-app 1.0.0  # Export WASM bundle
```

## ANSI Colors

The framework supports all standard ANSI colors:

**Standard Colors:**
- `Black`, `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`, `White`

**Bright Colors:**
- `BrightBlack`, `BrightRed`, `BrightGreen`, `BrightYellow`
- `BrightBlue`, `BrightMagenta`, `BrightCyan`, `BrightWhite`

## Cell Attributes

Each character cell supports:
- `Char` - The character to display
- `Foreground` - Text color
- `Background` - Background color
- `Bold` - Bold text attribute
- `Italic` - Italic text attribute
- `Underline` - Underlined text attribute

## Web Integration

When you export a WASM bundle, you get a complete web application that can be integrated into larger projects:

### Building the WASM Bundle

```bash
cd output-directory
./build.sh  # or build.bat on Windows
```

### Serving the Bundle

```bash
# Simple HTTP server
python3 -m http.server 8080

# Or use Go
go run -m http.server -addr :8080 .

# Then open http://localhost:8080
```

### JavaScript API

The exported bundle provides JavaScript hooks for interaction:

```javascript
// Available functions (automatically set up)
window.renderWASM()                    // Render current frame
window.handleKeyDownWASM(keyCode, key) // Handle key press
window.handleKeyUpWASM(keyCode, key)   // Handle key release
window.handleMouseClickWASM(x, y, btn) // Handle mouse click
window.handleMouseMoveWASM(x, y)       // Handle mouse move
window.resetWASM()                     // Reset to initial state

// Control the runtime
ansiEditor.start()  // Start animation loop
ansiEditor.stop()   // Stop animation loop
ansiEditor.reset()  // Reset to initial state
```

## Examples

### Simple Drawing Application

```go
scene := objects.NewScene("drawing", 40, 20)

// Draw a simple house
house := objects.NewObject("house", "House", 15, 8, 10, 8)
house.Canvas.DrawRect(0, 2, 10, 6, '#', canvas.Brown, canvas.Black)
house.Canvas.DrawLine(0, 2, 5, 0, '#', canvas.Red, canvas.Black)
house.Canvas.DrawLine(5, 0, 10, 2, '#', canvas.Red, canvas.Black)
house.Canvas.DrawChar(2, 5, 'D', canvas.Yellow, canvas.Black)

scene.AddObject(house)
export.ExportScene(scene, "./house-app", "house", "1.0.0")
```

### Animated Character

```go
// Create walking animation
char := objects.NewObject("walker", "Walking Character", 10, 10, 1, 1)

frames := []rune{'|', '/', '-', '\\'}
for _, frame := range frames {
    canvas := canvas.NewCanvas(1, 1)
    canvas.DrawChar(0, 0, frame, canvas.Green, canvas.Black)
    char.AddFrame(canvas)
}

scene.AddObject(char)

// In your game loop:
scene.UpdateAnimations() // Advance to next frame
```

### Interactive Game

The WASM export includes event handling for creating interactive applications:

```go
// Your main.go in the exported bundle can handle events:
func handleKeyDownWASM(this js.Value, args []js.Value) interface{} {
    keyCode := args[0].Int()
    
    switch keyCode {
    case 37: // Left arrow
        player.Move(player.X-1, player.Y)
    case 39: // Right arrow
        player.Move(player.X+1, player.Y)
    }
    
    return nil
}
```

## Best Practices

### Performance
- Keep canvas sizes reasonable (80x25 to 120x40 typically)
- Limit the number of animated objects
- Use object layering (Z-index) efficiently

### Design
- Use consistent color schemes
- Keep text readable with good contrast
- Consider the target display environment

### Export
- Test WASM bundles in the target environment
- Optimize bundle size by removing unused scenes
- Use appropriate frame rates (30fps is usually sufficient)

## API Reference

See the individual package documentation for complete API details:

- [Canvas Package](../pkg/canvas/) - Drawing and canvas operations
- [Objects Package](../pkg/objects/) - Object and scene management
- [ANSI Package](../pkg/ansi/) - File format handling
- [Export Package](../pkg/export/) - WASM export functionality

## Contributing

Contributions are welcome! Please see the main repository for contribution guidelines.

## License

This project is open source. See the LICENSE file for details.