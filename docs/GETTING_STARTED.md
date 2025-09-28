# Getting Started with ANSI Editor

This guide will help you get started with the ANSI Editor framework in just a few minutes.

## What is ANSI Editor?

ANSI Editor is a Go framework for creating paint tool applications that work with ANSI art. You can:

- Create and edit ANSI artwork programmatically
- Build interactive applications with moveable objects
- Export your creations as WebAssembly bundles that run in browsers
- Create BBS-style interfaces, simple games, and animated art

## Installation

Make sure you have Go 1.19 or later installed, then:

```bash
# Clone or download the framework
git clone https://github.com/evevioletrose-hash/ansi-editor
cd ansi-editor

# Initialize the Go module (if needed)
go mod tidy
```

## Your First ANSI Art

Let's create a simple "Hello World" example:

```go
package main

import (
    "fmt"
    "github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
    "github.com/evevioletrose-hash/ansi-editor/pkg/objects"
)

func main() {
    // Create a scene (like a stage for your art)
    scene := objects.NewScene("hello", 40, 10)
    
    // Draw text on the background
    scene.Background.DrawString(5, 3, "Hello ANSI World!", canvas.BrightGreen, canvas.Black)
    
    // Draw a border around it
    scene.Background.DrawRect(3, 2, 25, 3, '*', canvas.Yellow, canvas.Black)
    
    // Show what we created
    rendered := scene.Render()
    fmt.Println(rendered.ToANSI())
}
```

Save this as `hello.go` and run:

```bash
go run hello.go
```

You should see colorful text with ANSI colors displayed in your terminal!

## Adding Interactive Objects

Now let's add a moveable object:

```go
package main

import (
    "fmt"
    "github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
    "github.com/evevioletrose-hash/ansi-editor/pkg/objects"
)

func main() {
    // Create scene
    scene := objects.NewScene("game", 30, 15)
    
    // Draw a simple background
    scene.Background.FillRect(0, 0, 30, 15, '.', canvas.Green, canvas.Black)
    scene.Background.DrawString(8, 1, "Simple Game Area", canvas.White, canvas.Black)
    
    // Create a player character
    player := objects.NewObject("player", "Player", 15, 7, 1, 1)
    player.Canvas.DrawChar(0, 0, '@', canvas.BrightYellow, canvas.Black)
    scene.AddObject(player)
    
    // Create an enemy
    enemy := objects.NewObject("enemy", "Enemy", 10, 10, 1, 1)
    enemy.Canvas.DrawChar(0, 0, 'X', canvas.Red, canvas.Black)
    scene.AddObject(enemy)
    
    // Show the scene
    rendered := scene.Render()
    fmt.Println(rendered.ToANSI())
    
    // Move the player and show again
    player.Move(18, 9)
    fmt.Println("\nAfter moving player:")
    rendered = scene.Render()
    fmt.Println(rendered.ToANSI())
}
```

## Using the Command-Line Editor

The framework includes an interactive editor you can use right away:

```bash
# Run the editor
go run cmd/editor/main.go
```

Try these commands in the editor:

```
new 60 20                           # Create a 60x20 canvas
text 5 5 "Welcome!" green black     # Draw text
rect 10 8 20 6 # blue black         # Draw a rectangle
draw 15 10 @ yellow black           # Draw a character
show                                # Display your artwork
export wasm ./my-app demo 1.0       # Export as web app
```

## Creating Your First Web App

Let's export a simple scene as a web application:

```go
package main

import (
    "github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
    "github.com/evevioletrose-hash/ansi-editor/pkg/objects"
    "github.com/evevioletrose-hash/ansi-editor/pkg/export"
)

func main() {
    // Create an animated demo
    scene := objects.NewScene("demo", 50, 20)
    
    // Background
    scene.Background.DrawString(10, 3, "My First ANSI Web App!", canvas.BrightCyan, canvas.Black)
    scene.Background.DrawRect(5, 5, 40, 10, '=', canvas.Blue, canvas.Black)
    scene.Background.DrawString(8, 8, "This runs in your browser!", canvas.White, canvas.Black)
    scene.Background.DrawString(8, 10, "Built with Go + WebAssembly", canvas.Green, canvas.Black)
    
    // Add a bouncing ball
    ball := objects.NewObject("ball", "Ball", 25, 12, 1, 1)
    ball.Canvas.DrawChar(0, 0, 'o', canvas.Red, canvas.Black)
    scene.AddObject(ball)
    
    // Export as web app
    export.ExportScene(scene, "./my-web-app", "demo-app", "1.0.0")
    
    println("Web app exported to ./my-web-app")
    println("To run it:")
    println("  cd my-web-app")
    println("  ./build.sh")
    println("  python3 -m http.server 8080")
    println("  # Open http://localhost:8080")
}
```

Run this, then follow the instructions to build and serve your web app!

## Next Steps

### Learn More
- Read the full [Documentation](README.md)
- Check out the [Examples](../examples/)
- Explore the [Command Reference](COMMANDS.md)

### Try These Ideas
1. **Create a simple game** - Add keyboard controls to move objects around
2. **Build an interface** - Create menus and buttons using rectangles and text
3. **Make animations** - Use multiple frames to create animated sprites
4. **Load ANSI files** - Import existing ANSI art and add interactivity

### Common Patterns

**Drawing shapes:**
```go
canvas.DrawRect(x, y, width, height, '#', color, background)
canvas.FillRect(x, y, width, height, ' ', color, background)
canvas.DrawLine(x1, y1, x2, y2, '*', color, background)
```

**Managing objects:**
```go
obj := objects.NewObject("id", "name", x, y, width, height)
obj.Move(newX, newY)
scene.AddObject(obj)
```

**Animation:**
```go
// Create frames
frame1 := canvas.NewCanvas(width, height)
frame2 := canvas.NewCanvas(width, height)
// ... draw on frames ...

// Add to object
obj.AddFrame(frame1)
obj.AddFrame(frame2)

// In your loop:
scene.UpdateAnimations()
```

## Getting Help

- Check the documentation in the `docs/` folder
- Look at examples in the `examples/` folder  
- Run the interactive editor with `go run cmd/editor/main.go`
- Use the `help` command in the editor for quick reference

Have fun creating ANSI art applications!