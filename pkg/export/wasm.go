// Package export provides functionality for exporting ANSI projects to WASM bundles
package export

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/evevioletrose-hash/ansi-editor/pkg/objects"
)

// WASMBundle represents an exported WASM bundle
type WASMBundle struct {
	Name        string           `json:"name"`
	Version     string           `json:"version"`
	Scenes      []*objects.Scene `json:"scenes"`
	CurrentScene string          `json:"currentScene"`
	Config      BundleConfig     `json:"config"`
}

// BundleConfig contains configuration for the WASM bundle
type BundleConfig struct {
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	FPS            int    `json:"fps"`
	CanvasID       string `json:"canvasId"`
	AutoStart      bool   `json:"autoStart"`
	EnableKeyboard bool   `json:"enableKeyboard"`
	EnableMouse    bool   `json:"enableMouse"`
}

// DefaultConfig returns a default bundle configuration
func DefaultConfig() BundleConfig {
	return BundleConfig{
		Width:          80,
		Height:         25,
		FPS:            30,
		CanvasID:       "ansi-canvas",
		AutoStart:      true,
		EnableKeyboard: true,
		EnableMouse:    true,
	}
}

// NewWASMBundle creates a new WASM bundle
func NewWASMBundle(name, version string) *WASMBundle {
	return &WASMBundle{
		Name:    name,
		Version: version,
		Scenes:  make([]*objects.Scene, 0),
		Config:  DefaultConfig(),
	}
}

// AddScene adds a scene to the bundle
func (wb *WASMBundle) AddScene(scene *objects.Scene) {
	wb.Scenes = append(wb.Scenes, scene)
	if wb.CurrentScene == "" {
		wb.CurrentScene = scene.Name
	}
}

// Export exports the bundle to a directory
func (wb *WASMBundle) Export(outputDir string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Export bundle data as JSON
	bundleData, err := json.MarshalIndent(wb, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal bundle data: %w", err)
	}

	bundleFile := filepath.Join(outputDir, "bundle.json")
	if err := os.WriteFile(bundleFile, bundleData, 0644); err != nil {
		return fmt.Errorf("failed to write bundle data: %w", err)
	}

	// Generate HTML file
	if err := wb.generateHTML(outputDir); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Generate JavaScript runtime
	if err := wb.generateJS(outputDir); err != nil {
		return fmt.Errorf("failed to generate JavaScript: %w", err)
	}

	// Generate Go WASM source
	if err := wb.generateWASMSource(outputDir); err != nil {
		return fmt.Errorf("failed to generate WASM source: %w", err)
	}

	// Generate build script
	if err := wb.generateBuildScript(outputDir); err != nil {
		return fmt.Errorf("failed to generate build script: %w", err)
	}

	return nil
}

// generateHTML creates the HTML wrapper
func (wb *WASMBundle) generateHTML(outputDir string) error {
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>{{.Name}}</title>
    <style>
        body {
            margin: 0;
            padding: 20px;
            background-color: #000;
            color: #fff;
            font-family: 'Courier New', monospace;
        }
        #{{.Config.CanvasID}} {
            border: 1px solid #333;
            background-color: #000;
            font-family: 'Courier New', monospace;
            font-size: 12px;
            line-height: 1.2;
            white-space: pre;
        }
        .controls {
            margin-top: 10px;
        }
        button {
            background-color: #333;
            color: #fff;
            border: 1px solid #555;
            padding: 5px 10px;
            margin-right: 5px;
            cursor: pointer;
        }
        button:hover {
            background-color: #555;
        }
    </style>
</head>
<body>
    <h1>{{.Name}}</h1>
    <div id="{{.Config.CanvasID}}"></div>
    <div class="controls">
        <button onclick="ansiEditor.start()">Start</button>
        <button onclick="ansiEditor.stop()">Stop</button>
        <button onclick="ansiEditor.reset()">Reset</button>
    </div>
    
    <script src="wasm_exec.js"></script>
    <script src="ansi-runtime.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("ansi-bundle.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        });
    </script>
</body>
</html>`

	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	htmlFile := filepath.Join(outputDir, "index.html")
	file, err := os.Create(htmlFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, wb)
}

// generateJS creates the JavaScript runtime
func (wb *WASMBundle) generateJS(outputDir string) error {
	jsContent := `// ANSI Editor Runtime
class ANSIEditor {
    constructor(config) {
        this.config = config;
        this.canvas = document.getElementById(config.canvasId);
        this.running = false;
        this.animationId = null;
        this.lastFrameTime = 0;
        this.frameInterval = 1000 / config.fps;
        
        // Load bundle data
        this.loadBundle();
        
        // Setup event listeners
        this.setupEventListeners();
    }
    
    async loadBundle() {
        try {
            const response = await fetch('bundle.json');
            this.bundle = await response.json();
            console.log('Bundle loaded:', this.bundle);
        } catch (error) {
            console.error('Failed to load bundle:', error);
        }
    }
    
    setupEventListeners() {
        if (this.config.enableKeyboard) {
            document.addEventListener('keydown', (e) => this.handleKeyDown(e));
            document.addEventListener('keyup', (e) => this.handleKeyUp(e));
        }
        
        if (this.config.enableMouse) {
            this.canvas.addEventListener('click', (e) => this.handleMouseClick(e));
            this.canvas.addEventListener('mousemove', (e) => this.handleMouseMove(e));
        }
    }
    
    start() {
        if (!this.running) {
            this.running = true;
            this.animate();
            console.log('ANSI Editor started');
        }
    }
    
    stop() {
        this.running = false;
        if (this.animationId) {
            cancelAnimationFrame(this.animationId);
            this.animationId = null;
        }
        console.log('ANSI Editor stopped');
    }
    
    reset() {
        this.stop();
        // Reset to initial state
        if (this.bundle && window.resetWASM) {
            window.resetWASM();
        }
        console.log('ANSI Editor reset');
    }
    
    animate(currentTime = 0) {
        if (!this.running) return;
        
        if (currentTime - this.lastFrameTime >= this.frameInterval) {
            this.render();
            this.lastFrameTime = currentTime;
        }
        
        this.animationId = requestAnimationFrame((time) => this.animate(time));
    }
    
    render() {
        // Call WASM render function if available
        if (window.renderWASM) {
            const output = window.renderWASM();
            if (output && this.canvas) {
                this.canvas.innerHTML = output;
            }
        }
    }
    
    handleKeyDown(event) {
        if (window.handleKeyDownWASM) {
            window.handleKeyDownWASM(event.keyCode, event.key);
        }
    }
    
    handleKeyUp(event) {
        if (window.handleKeyUpWASM) {
            window.handleKeyUpWASM(event.keyCode, event.key);
        }
    }
    
    handleMouseClick(event) {
        const rect = this.canvas.getBoundingClientRect();
        const x = Math.floor((event.clientX - rect.left) / 8); // Assuming 8px char width
        const y = Math.floor((event.clientY - rect.top) / 16);  // Assuming 16px char height
        
        if (window.handleMouseClickWASM) {
            window.handleMouseClickWASM(x, y, event.button);
        }
    }
    
    handleMouseMove(event) {
        const rect = this.canvas.getBoundingClientRect();
        const x = Math.floor((event.clientX - rect.left) / 8);
        const y = Math.floor((event.clientY - rect.top) / 16);
        
        if (window.handleMouseMoveWASM) {
            window.handleMouseMoveWASM(x, y);
        }
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    const config = {
        canvasId: '` + wb.Config.CanvasID + `',
        fps: ` + fmt.Sprintf("%d", wb.Config.FPS) + `,
        enableKeyboard: ` + fmt.Sprintf("%t", wb.Config.EnableKeyboard) + `,
        enableMouse: ` + fmt.Sprintf("%t", wb.Config.EnableMouse) + `
    };
    
    window.ansiEditor = new ANSIEditor(config);
    
    if (` + fmt.Sprintf("%t", wb.Config.AutoStart) + `) {
        window.ansiEditor.start();
    }
});`

	jsFile := filepath.Join(outputDir, "ansi-runtime.js")
	return os.WriteFile(jsFile, []byte(jsContent), 0644)
}

// generateWASMSource creates the Go source for WASM compilation
func (wb *WASMBundle) generateWASMSource(outputDir string) error {
	wasmSource := `package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	
	"github.com/evevioletrose-hash/ansi-editor/pkg/objects"
	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
)

var (
	currentScene *objects.Scene
	bundle       map[string]interface{}
)

func main() {
	// Register WASM functions
	js.Global().Set("renderWASM", js.FuncOf(renderWASM))
	js.Global().Set("handleKeyDownWASM", js.FuncOf(handleKeyDownWASM))
	js.Global().Set("handleKeyUpWASM", js.FuncOf(handleKeyUpWASM))
	js.Global().Set("handleMouseClickWASM", js.FuncOf(handleMouseClickWASM))
	js.Global().Set("handleMouseMoveWASM", js.FuncOf(handleMouseMoveWASM))
	js.Global().Set("resetWASM", js.FuncOf(resetWASM))
	
	// Load bundle
	loadBundle()
	
	// Keep the program running
	select {}
}

func loadBundle() error {
	// In a real implementation, this would load the bundle.json
	// For now, create a simple demo scene
	scene := objects.NewScene("demo", ` + fmt.Sprintf("%d", wb.Config.Width) + `, ` + fmt.Sprintf("%d", wb.Config.Height) + `)
	
	// Add some demo content
	scene.Background.DrawString(10, 5, "Hello ANSI World!", canvas.Green, canvas.Black)
	scene.Background.DrawRect(5, 3, 30, 8, '#', canvas.Yellow, canvas.Black)
	
	// Add a demo object
	obj := objects.NewObject("demo-obj", "Demo Object", 20, 10, 10, 5)
	obj.Canvas.FillRect(0, 0, 10, 5, '*', canvas.Red, canvas.Black)
	scene.AddObject(obj)
	
	currentScene = scene
	return nil
}

func renderWASM(this js.Value, args []js.Value) interface{} {
	if currentScene == nil {
		return "<pre>Loading...</pre>"
	}
	
	// Update animations
	currentScene.UpdateAnimations()
	
	// Render scene
	rendered := currentScene.Render()
	ansiOutput := rendered.ToANSI()
	
	// Convert to HTML
	htmlOutput := fmt.Sprintf("<pre>%s</pre>", ansiOutput)
	return htmlOutput
}

func handleKeyDownWASM(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return nil
	}
	
	keyCode := args[0].Int()
	key := args[1].String()
	
	// Handle key input - customize as needed
	fmt.Printf("Key down: %d (%s)\n", keyCode, key)
	
	return nil
}

func handleKeyUpWASM(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return nil
	}
	
	keyCode := args[0].Int()
	key := args[1].String()
	
	fmt.Printf("Key up: %d (%s)\n", keyCode, key)
	
	return nil
}

func handleMouseClickWASM(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return nil
	}
	
	x := args[0].Int()
	y := args[1].Int()
	button := args[2].Int()
	
	fmt.Printf("Mouse click at (%d, %d) button: %d\n", x, y, button)
	
	return nil
}

func handleMouseMoveWASM(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return nil
	}
	
	x := args[0].Int()
	y := args[1].Int()
	
	// Update mouse position - customize as needed
	// fmt.Printf("Mouse move: (%d, %d)\n", x, y)
	
	return nil
}

func resetWASM(this js.Value, args []js.Value) interface{} {
	// Reset to initial state
	loadBundle()
	return nil
}`

	wasmFile := filepath.Join(outputDir, "main.go")
	return os.WriteFile(wasmFile, []byte(wasmSource), 0644)
}

// generateBuildScript creates a build script for the WASM bundle
func (wb *WASMBundle) generateBuildScript(outputDir string) error {
	buildScript := `#!/bin/bash
# Build script for ANSI Editor WASM bundle

echo "Building WASM bundle..."

# Set WASM environment
export GOOS=js
export GOARCH=wasm

# Build the WASM binary
go build -o ansi-bundle.wasm main.go

# Copy wasm_exec.js from Go installation
WASM_EXEC=$(go env GOROOT)/misc/wasm/wasm_exec.js
if [ -f "$WASM_EXEC" ]; then
    cp "$WASM_EXEC" .
else
    echo "Warning: wasm_exec.js not found. You may need to copy it manually."
fi

echo "Build complete! Open index.html in a web server to run the bundle."
echo "Note: Due to CORS restrictions, you cannot open index.html directly in a browser."
echo "Use a local web server, for example:"
echo "  python3 -m http.server 8080"
echo "  # or"
echo "  go run -m http.server -addr :8080 ."
`

	buildFile := filepath.Join(outputDir, "build.sh")
	if err := os.WriteFile(buildFile, []byte(buildScript), 0755); err != nil {
		return err
	}

	// Also create a Windows batch file
	winBuildScript := `@echo off
REM Build script for ANSI Editor WASM bundle

echo Building WASM bundle...

REM Set WASM environment
set GOOS=js
set GOARCH=wasm

REM Build the WASM binary
go build -o ansi-bundle.wasm main.go

REM Copy wasm_exec.js from Go installation
for /f "tokens=*" %%i in ('go env GOROOT') do set GOROOT=%%i
set WASM_EXEC=%GOROOT%\misc\wasm\wasm_exec.js
if exist "%WASM_EXEC%" (
    copy "%WASM_EXEC%" .
) else (
    echo Warning: wasm_exec.js not found. You may need to copy it manually.
)

echo Build complete! Open index.html in a web server to run the bundle.
echo Note: Due to CORS restrictions, you cannot open index.html directly in a browser.
echo Use a local web server, for example:
echo   python -m http.server 8080
pause
`

	winBuildFile := filepath.Join(outputDir, "build.bat")
	return os.WriteFile(winBuildFile, []byte(winBuildScript), 0644)
}

// ExportScene is a convenience function to export a single scene
func ExportScene(scene *objects.Scene, outputDir, name, version string) error {
	bundle := NewWASMBundle(name, version)
	bundle.AddScene(scene)
	return bundle.Export(outputDir)
}

// LoadBundle loads a WASM bundle from a JSON file
func LoadBundle(reader io.Reader) (*WASMBundle, error) {
	var bundle WASMBundle
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&bundle)
	return &bundle, err
}