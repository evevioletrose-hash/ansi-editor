# ANSI Editor Framework

A powerful Go-based framework for creating interactive ANSI art applications that can be exported as WebAssembly bundles for the web.

## 🎨 Features

- **Canvas System**: 2D grid-based drawing with full ANSI color support
- **Object Management**: Create and manage moveable, animated ANSI art objects  
- **Scene System**: Compose complex scenes with backgrounds and layered objects
- **Animation Support**: Frame-based animation system for dynamic content
- **ANSI File I/O**: Load and save standard ANSI art files (.ans format)
- **WASM Export**: Export complete projects as web-ready WebAssembly bundles
- **Interactive Editor**: Command-line tool for creating and editing ANSI art
- **Web Integration**: Generated bundles include HTML, JavaScript, and build scripts

## 🚀 Quick Start

### Installation

```bash
git clone https://github.com/evevioletrose-hash/ansi-editor
cd ansi-editor
go mod tidy
```

### Hello World Example

```go
package main

import (
    "fmt"
    "github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
    "github.com/evevioletrose-hash/ansi-editor/pkg/objects"
    "github.com/evevioletrose-hash/ansi-editor/pkg/export"
)

func main() {
    // Create a scene
    scene := objects.NewScene("hello", 40, 10)
    
    // Draw background
    scene.Background.DrawString(5, 3, "Hello ANSI World!", canvas.BrightGreen, canvas.Black)
    scene.Background.DrawRect(3, 2, 25, 3, '*', canvas.Yellow, canvas.Black)
    
    // Add animated object
    obj := objects.NewObject("spinner", "Spinner", 30, 5, 1, 1)
    frames := []rune{'|', '/', '-', '\\'}
    for _, r := range frames {
        frame := canvas.NewCanvas(1, 1)
        frame.DrawChar(0, 0, r, canvas.Cyan, canvas.Black)
        obj.AddFrame(frame)
    }
    scene.AddObject(obj)
    
    // Export as web app
    export.ExportScene(scene, "./hello-web", "hello-world", "1.0.0")
}
```

### Try the Interactive Editor

```bash
go run cmd/editor/main.go
```

Then try these commands:
```
new 60 20
text 10 5 "Welcome to ANSI Editor!" green black
rect 5 3 50 8 # blue black
object add player Player 25 10 3 3
show
export wasm ./my-app demo 1.0
```

## 📖 Documentation

- **[Getting Started Guide](docs/GETTING_STARTED.md)** - Your first steps with the framework
- **[Full Documentation](docs/README.md)** - Complete API reference and guides
- **[Examples](examples/)** - Sample applications and code snippets

## 🎮 Use Cases

- **BBS-Style Interfaces**: Create retro terminal interfaces
- **Simple Games**: Build text-based games with animations
- **Interactive Art**: Create dynamic ANSI artwork
- **Educational Tools**: Teach programming concepts with visual feedback
- **Prototyping**: Quickly mock up console-style applications
- **Web Integration**: Embed ANSI art in modern web applications

## 🏗️ Architecture

The framework consists of several key packages:

- `pkg/canvas` - Core drawing and rendering system
- `pkg/objects` - Object and scene management
- `pkg/ansi` - ANSI file format parsing and generation
- `pkg/export` - WebAssembly export functionality
- `cmd/editor` - Interactive command-line editor
- `cmd/example` - Example application

## 🌐 Web Export

When you export a project, you get a complete web application:

```
output-folder/
├── index.html          # Main web page
├── ansi-runtime.js     # JavaScript runtime
├── main.go            # Go WASM source
├── bundle.json        # Scene data
├── build.sh           # Unix build script
├── build.bat          # Windows build script
└── README.md          # Instructions
```

To build and run:
```bash
cd output-folder
./build.sh
python3 -m http.server 8080
# Open http://localhost:8080
```

## 🎨 Color Support

Full ANSI color palette support:
- Standard colors: Black, Red, Green, Yellow, Blue, Magenta, Cyan, White
- Bright variants: BrightRed, BrightGreen, etc.
- Text attributes: Bold, Italic, Underline

## 🔧 Examples

### Simple Drawing
```go
canvas := canvas.NewCanvas(20, 10)
canvas.DrawRect(2, 2, 16, 6, '#', canvas.Blue, canvas.Black)
canvas.DrawString(4, 4, "Hello!", canvas.Yellow, canvas.Black)
```

### Animated Sprite
```go
sprite := objects.NewObject("ball", "Bouncing Ball", 10, 5, 1, 1)
for _, char := range []rune{'o', 'O', '0', 'O'} {
    frame := canvas.NewCanvas(1, 1)
    frame.DrawChar(0, 0, char, canvas.Red, canvas.Black)
    sprite.AddFrame(frame)
}
```

### Load ANSI File
```go
file, _ := os.Open("artwork.ans")
parser := ansi.NewParser()
canvas, _ := parser.ParseFile(file, 80, 25)
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

## 📄 License

This project is open source. See LICENSE file for details.

---

**Built with Go** • **Powered by WebAssembly** • **Inspired by BBS culture**
