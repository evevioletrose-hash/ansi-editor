# ANSI Editor

A Go-based ANSI artwork editor framework that provides a generic starter kit for building ANSI art editors. This library allows you to create ANSI front-end packages that can be exported to various formats including WebAssembly, HTML, JavaScript, and JSON for use in terminal interfaces, web applications, and games.

## Features

### Core Framework
- **Grid-based Canvas System**: Paint tool-like interface for character placement
- **Multi-page Support**: Create and manage multiple ANSI art pages/scenes
- **ANSI File Format Support**: Load and save standard ANSI files (.ans, .ansi, .txt)
- **Drawing Tools**: Pencil, brush, line, rectangle, circle, text, fill, and erase tools
- **Color Management**: Full ANSI color palette support (16 colors + bright variants)
- **Export System**: Generate frontend packages for various target platforms

### Export Formats
- **JSON**: Data format for programmatic access
- **HTML**: Web-ready ANSI art with CSS styling
- **JavaScript**: Interactive modules for web applications
- **WebAssembly**: High-performance WASM packages (build system included)
- **Plain Text**: Character-only output

### Use Cases
- Terminal-based games (like Doom in ANSI)
- ANSI art creation and editing
- Frontend package generation for web applications
- Text-based user interfaces
- ASCII/ANSI animation systems

## Installation

```bash
go get github.com/evevioletrose-hash/ansi-editor
```

## Quick Start

### Command Line Tool

Create a sample ANSI file:
```bash
go run cmd/ansi-editor/main.go -sample demo.ans
```

Create a new ANSI canvas:
```bash
go run cmd/ansi-editor/main.go -create myart.ans -dimensions 60x20 -text "Hello ANSI" -color red
```

Load and display an ANSI file:
```bash
go run cmd/ansi-editor/main.go -load demo.ans
```

Export to HTML:
```bash
go run cmd/ansi-editor/main.go -load demo.ans -export html:web_output
```

### Programmatic Usage

```go
package main

import (
    "github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
    "github.com/evevioletrose-hash/ansi-editor/pkg/editor"
    "github.com/evevioletrose-hash/ansi-editor/pkg/export"
)

func main() {
    // Create a new editor
    ed := editor.NewEditor()
    
    // Create a new page
    page := ed.NewPage("My Art", 80, 25)
    
    // Draw some content
    drawState := editor.NewDrawingState()
    drawState.Character = '█'
    drawState.FgColor = canvas.ColorRed
    
    ed.DrawText(10, 5, "Hello ANSI!", drawState)
    ed.DrawRectangle(5, 3, 30, 5, drawState)
    
    // Save the work
    ed.SaveCurrentPage("myart.ans")
    
    // Export as HTML
    options := &export.ExportOptions{
        Format:      export.FormatHTML,
        OutputDir:   "html_output",
        ProjectName: "My ANSI Art",
    }
    
    exporter := export.NewExporter(options)
    exporter.ExportProject(ed)
}
```

## Architecture

The framework is designed with modularity and extensibility in mind:

### Core Components

- **`pkg/canvas/`**: Grid-based drawing surface with ANSI rendering
- **`pkg/editor/`**: Multi-page editor with drawing tools
- **`pkg/ansi/`**: ANSI file format parser and writer
- **`pkg/export/`**: Export system for various output formats
- **`cmd/ansi-editor/`**: Command-line interface tool

### Canvas System

The canvas provides a grid-based drawing surface where each cell contains:
- Character (rune)
- Foreground color (ANSI code)
- Background color (ANSI code)  
- Text attributes (bold, italic, underline, etc.)

```go
// Create a 80x25 canvas
canvas := canvas.NewCanvas(80, 25)

// Set individual cells
canvas.SetCell(10, 5, '█', canvas.ColorRed, canvas.ColorBlack, 0)

// Draw text
canvas.DrawText(0, 0, "Hello World", canvas.ColorWhite, canvas.ColorBlue, canvas.AttrBold)

// Render as ANSI
ansiOutput := canvas.RenderANSI()
```

### Drawing Tools

The editor provides various drawing tools:

```go
drawState := editor.NewDrawingState()
drawState.Tool = editor.ToolBrush
drawState.Character = '▓'
drawState.FgColor = canvas.ColorGreen
drawState.BrushSize = 3

// Draw with the configured tool
editor.DrawPoint(x, y, drawState)
editor.DrawLine(x1, y1, x2, y2, drawState)
editor.DrawRectangle(x, y, width, height, drawState)
editor.DrawCircle(centerX, centerY, radius, drawState)
```

### File Format Support

Load and save ANSI files:

```go
// Load an ANSI file
canvas, err := ansi.LoadFromFile("artwork.ans")

// Save a canvas
err = ansi.SaveToFile(canvas, "output.ans")

// Get file information
info, err := ansi.GetFileInfo("artwork.ans")
fmt.Printf("Dimensions: %dx%d\n", info.Width, info.Height)
```

### Export System

Generate frontend packages:

```go
// Export as JSON
options := &export.ExportOptions{
    Format:      export.FormatJSON,
    OutputDir:   "json_output",
    ProjectName: "My Project",
}

exporter := export.NewExporter(options)
err := exporter.ExportProject(editor)

// Export as WebAssembly (includes build scripts)
options.Format = export.FormatWebAssembly
options.OutputDir = "wasm_output"
err = exporter.ExportProject(editor)
```

## API Reference

### Canvas API

```go
// Creation
func NewCanvas(width, height int) *Canvas

// Cell manipulation
func (c *Canvas) SetCell(x, y int, char rune, fgColor, bgColor, attributes int) error
func (c *Canvas) GetCell(x, y int) (Cell, error)
func (c *Canvas) Clear()
func (c *Canvas) Resize(newWidth, newHeight int)

// Drawing operations
func (c *Canvas) DrawText(x, y int, text string, fgColor, bgColor, attributes int) error
func (c *Canvas) FillRect(x, y, width, height int, char rune, fgColor, bgColor, attributes int) error

// Rendering
func (c *Canvas) RenderANSI() string
func (c *Canvas) RenderPlainText() string
```

### Editor API

```go
// Editor management
func NewEditor() *Editor
func (e *Editor) NewPage(name string, width, height int) *Page
func (e *Editor) GetCurrentPage() *Page
func (e *Editor) SetCurrentPage(index int) error

// File operations
func (e *Editor) LoadPage(filename, pageName string) error
func (e *Editor) SaveCurrentPage(filename string) error

// Drawing tools
func (e *Editor) DrawPoint(x, y int, state *DrawingState) error
func (e *Editor) DrawLine(x1, y1, x2, y2 int, state *DrawingState) error
func (e *Editor) DrawRectangle(x, y, width, height int, state *DrawingState) error
func (e *Editor) DrawCircle(centerX, centerY, radius int, state *DrawingState) error
func (e *Editor) DrawText(x, y int, text string, state *DrawingState) error
```

## Color Palette

The framework supports the full ANSI color palette:

### Standard Colors (30-37, 40-47)
- Black, Red, Green, Yellow, Blue, Magenta, Cyan, White

### Bright Colors (90-97, 100-107)  
- Bright versions of all standard colors

### Text Attributes
- Bold, Dim, Italic, Underline, Blink, Reverse, Strikethrough

## Building and Testing

### Run Tests
```bash
go test ./tests/...
```

### Build CLI Tool
```bash
go build -o ansi-editor cmd/ansi-editor/main.go
```

### Run Examples
```bash
go run examples/basic_usage.go
```

### Build for WebAssembly
```bash
GOOS=js GOARCH=wasm go build -o ansi-editor.wasm cmd/ansi-editor/main.go
```

## Examples

See the `examples/` directory for complete usage examples:

- `basic_usage.go`: Demonstrates core functionality
- `demo.ans`: Sample ANSI file (generated by CLI tool)
- `html_output/`: Example HTML export

## Contributing

This is a starter framework designed to be extended. Key areas for contribution:

1. **Additional Export Formats**: Add support for more output formats
2. **Advanced Drawing Tools**: Implement more sophisticated drawing tools
3. **Animation Support**: Add frame-based animation capabilities
4. **Palette Management**: Enhanced color palette tools
5. **Performance Optimizations**: Optimize for large canvases
6. **ANSI Parser**: Extend support for more ANSI escape sequences

## License

This project is open source. See LICENSE file for details.

## Use Cases in Detail

### Terminal Games
Build text-based games like Doom using the ANSI framework:
```go
// Create game scenes
gameEditor := editor.NewEditor()
level1 := gameEditor.NewPage("Level 1", 80, 25)
level2 := gameEditor.NewPage("Level 2", 80, 25)

// Export for terminal rendering
export.ExportProject(gameEditor)
```

### Web Applications
Generate web-ready ANSI art:
```go
// Export as HTML with CSS styling
options := &export.ExportOptions{
    Format: export.FormatHTML,
    OutputDir: "web_assets",
}
```

### Frontend Packages
Create reusable ANSI components:
```go
// Export as JavaScript modules
options := &export.ExportOptions{
    Format: export.FormatJavaScript,
    OutputDir: "js_modules",
}
```
