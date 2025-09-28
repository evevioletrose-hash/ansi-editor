package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/evevioletrose-hash/ansi-editor/pkg/editor"
)

// ExportFormat represents different export formats
type ExportFormat int

const (
	FormatJSON ExportFormat = iota
	FormatWebAssembly
	FormatJavaScript
	FormatHTML
	FormatCSS
)

// ExportOptions contains configuration for export operations
type ExportOptions struct {
	Format      ExportFormat
	OutputDir   string
	ProjectName string
	Template    string
	Minify      bool
	Animations  bool
}

// Exporter handles exporting editor content to various formats
type Exporter struct {
	options *ExportOptions
}

// NewExporter creates a new exporter with the specified options
func NewExporter(options *ExportOptions) *Exporter {
	return &Exporter{
		options: options,
	}
}

// ExportProject exports an entire editor project
func (exp *Exporter) ExportProject(editor *editor.Editor) error {
	switch exp.options.Format {
	case FormatJSON:
		return exp.exportJSON(editor)
	case FormatHTML:
		return exp.exportHTML(editor)
	case FormatJavaScript:
		return exp.exportJavaScript(editor)
	case FormatWebAssembly:
		return exp.exportWASM(editor)
	default:
		return fmt.Errorf("unsupported export format: %v", exp.options.Format)
	}
}

// exportJSON exports the project as JSON data
func (exp *Exporter) exportJSON(editor *editor.Editor) error {
	// Create output directory
	err := os.MkdirAll(exp.options.OutputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Convert editor to JSON-serializable format
	project := exp.editorToProjectData(editor)
	
	// Write project data
	projectFile := filepath.Join(exp.options.OutputDir, "project.json")
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project data: %w", err)
	}
	
	err = os.WriteFile(projectFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write project file: %w", err)
	}
	
	// Export individual pages
	for i, page := range editor.Pages {
		pageData := exp.pageToData(page, i)
		pageFile := filepath.Join(exp.options.OutputDir, fmt.Sprintf("page_%d.json", i))
		
		data, err := json.MarshalIndent(pageData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal page %d: %w", i, err)
		}
		
		err = os.WriteFile(pageFile, data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write page %d: %w", i, err)
		}
	}
	
	// Create a manifest file
	manifest := exp.createManifest(editor)
	manifestFile := filepath.Join(exp.options.OutputDir, "manifest.json")
	data, err = json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	
	err = os.WriteFile(manifestFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	
	return nil
}

// exportHTML exports the project as HTML with embedded ANSI content
func (exp *Exporter) exportHTML(editor *editor.Editor) error {
	err := os.MkdirAll(exp.options.OutputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Generate HTML for each page
	for i, page := range editor.Pages {
		html := exp.pageToHTML(page, i)
		htmlFile := filepath.Join(exp.options.OutputDir, fmt.Sprintf("page_%d.html", i))
		
		err = os.WriteFile(htmlFile, []byte(html), 0644)
		if err != nil {
			return fmt.Errorf("failed to write HTML page %d: %w", i, err)
		}
	}
	
	// Generate index page
	index := exp.createIndexHTML(editor)
	indexFile := filepath.Join(exp.options.OutputDir, "index.html")
	err = os.WriteFile(indexFile, []byte(index), 0644)
	if err != nil {
		return fmt.Errorf("failed to write index HTML: %w", err)
	}
	
	// Generate CSS file
	css := exp.generateCSS()
	cssFile := filepath.Join(exp.options.OutputDir, "ansi.css")
	err = os.WriteFile(cssFile, []byte(css), 0644)
	if err != nil {
		return fmt.Errorf("failed to write CSS: %w", err)
	}
	
	return nil
}

// exportJavaScript exports the project as JavaScript modules
func (exp *Exporter) exportJavaScript(editor *editor.Editor) error {
	err := os.MkdirAll(exp.options.OutputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Generate main JavaScript module
	js := exp.generateJavaScript(editor)
	jsFile := filepath.Join(exp.options.OutputDir, "ansi-editor.js")
	err = os.WriteFile(jsFile, []byte(js), 0644)
	if err != nil {
		return fmt.Errorf("failed to write JavaScript: %w", err)
	}
	
	// Export data as JSON for JavaScript consumption
	return exp.exportJSON(editor)
}

// exportWASM is a placeholder for WebAssembly export (would require Go WASM build)
func (exp *Exporter) exportWASM(editor *editor.Editor) error {
	// This would involve compiling the Go code to WebAssembly
	// For now, we'll generate the necessary structure and a placeholder
	
	err := os.MkdirAll(exp.options.OutputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Export JSON data for WASM consumption
	err = exp.exportJSON(editor)
	if err != nil {
		return err
	}
	
	// Generate WASM loader HTML
	wasmHTML := exp.generateWASMLoader(editor)
	htmlFile := filepath.Join(exp.options.OutputDir, "wasm.html")
	err = os.WriteFile(htmlFile, []byte(wasmHTML), 0644)
	if err != nil {
		return fmt.Errorf("failed to write WASM HTML: %w", err)
	}
	
	// Create build instructions
	buildScript := exp.generateWASMBuildScript()
	scriptFile := filepath.Join(exp.options.OutputDir, "build-wasm.sh")
	err = os.WriteFile(scriptFile, []byte(buildScript), 0755)
	if err != nil {
		return fmt.Errorf("failed to write build script: %w", err)
	}
	
	return nil
}

// Helper functions for data conversion

func (exp *Exporter) editorToProjectData(editor *editor.Editor) map[string]interface{} {
	pages := make([]map[string]interface{}, len(editor.Pages))
	for i, page := range editor.Pages {
		pages[i] = exp.pageToData(page, i)
	}
	
	return map[string]interface{}{
		"name":         exp.options.ProjectName,
		"pages":        pages,
		"currentPage":  editor.CurrentPage,
		"modified":     editor.Modified,
		"exportedAt":   "placeholder-timestamp",
		"version":      "1.0.0",
	}
}

func (exp *Exporter) pageToData(page *editor.Page, index int) map[string]interface{} {
	// Convert canvas grid to a more compact format
	cells := make([]map[string]interface{}, 0)
	
	for y := 0; y < page.Canvas.Height; y++ {
		for x := 0; x < page.Canvas.Width; x++ {
			cell, _ := page.Canvas.GetCell(x, y)
			// Only include non-default cells to reduce size
			if cell.Char != ' ' || cell.FgColor != 37 || cell.BgColor != 40 || cell.Attributes != 0 {
				cells = append(cells, map[string]interface{}{
					"x":    x,
					"y":    y,
					"char": string(cell.Char),
					"fg":   cell.FgColor,
					"bg":   cell.BgColor,
					"attr": cell.Attributes,
				})
			}
		}
	}
	
	return map[string]interface{}{
		"index":    index,
		"name":     page.Name,
		"width":    page.Canvas.Width,
		"height":   page.Canvas.Height,
		"modified": page.Modified,
		"cells":    cells,
	}
}

func (exp *Exporter) createManifest(editor *editor.Editor) map[string]interface{} {
	pages := make([]string, len(editor.Pages))
	for i, page := range editor.Pages {
		pages[i] = page.Name
	}
	
	return map[string]interface{}{
		"name":        exp.options.ProjectName,
		"description": "ANSI art project exported from ansi-editor",
		"version":     "1.0.0",
		"format":      "ansi-editor-export",
		"pages":       pages,
		"files": map[string]string{
			"project": "project.json",
			"pages":   "page_*.json",
		},
	}
}

func (exp *Exporter) pageToHTML(page *editor.Page, index int) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s - %s</title>
    <link rel="stylesheet" href="ansi.css">
    <meta charset="UTF-8">
</head>
<body>
    <div class="ansi-container">
        <h1>%s</h1>
        <pre class="ansi-content">`, exp.options.ProjectName, page.Name, page.Name))
	
	// Convert canvas to HTML with ANSI styling
	ansiContent := page.Canvas.RenderANSI()
	sb.WriteString(exp.ansiToHTMLSpans(ansiContent))
	
	sb.WriteString(`</pre>
    </div>
</body>
</html>`)
	
	return sb.String()
}

func (exp *Exporter) ansiToHTMLSpans(ansi string) string {
	// This is a simplified ANSI to HTML converter
	// In a full implementation, this would properly parse ANSI codes
	// and convert them to HTML spans with CSS classes
	return ansi // Placeholder - would need full ANSI parsing
}

func (exp *Exporter) createIndexHTML(editor *editor.Editor) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s - Index</title>
    <link rel="stylesheet" href="ansi.css">
</head>
<body>
    <div class="index-container">
        <h1>%s</h1>
        <h2>Pages</h2>
        <ul class="page-list">`, exp.options.ProjectName, exp.options.ProjectName))
	
	for i, page := range editor.Pages {
		sb.WriteString(fmt.Sprintf(`
            <li><a href="page_%d.html">%s</a> (%dx%d)</li>`, 
			i, page.Name, page.Canvas.Width, page.Canvas.Height))
	}
	
	sb.WriteString(`
        </ul>
    </div>
</body>
</html>`)
	
	return sb.String()
}

func (exp *Exporter) generateCSS() string {
	return `/* ANSI Editor Export CSS */
.ansi-container {
    font-family: 'Courier New', monospace;
    background-color: #000;
    color: #fff;
    padding: 20px;
}

.ansi-content {
    background-color: #000;
    color: #fff;
    font-family: 'Courier New', monospace;
    font-size: 14px;
    line-height: 1.2;
    white-space: pre;
    margin: 0;
}

.index-container {
    font-family: Arial, sans-serif;
    padding: 20px;
}

.page-list {
    list-style-type: none;
    padding: 0;
}

.page-list li {
    margin: 10px 0;
    padding: 10px;
    background-color: #f5f5f5;
    border-radius: 5px;
}

.page-list a {
    text-decoration: none;
    color: #333;
    font-weight: bold;
}`
}

func (exp *Exporter) generateJavaScript(editor *editor.Editor) string {
	return `// ANSI Editor Export JavaScript Module
class ANSIRenderer {
    constructor() {
        this.canvas = null;
        this.currentPage = 0;
        this.pages = [];
    }
    
    async loadProject(projectPath) {
        try {
            const response = await fetch(projectPath + '/project.json');
            const project = await response.json();
            this.pages = project.pages;
            this.currentPage = project.currentPage || 0;
            return project;
        } catch (error) {
            console.error('Failed to load project:', error);
            throw error;
        }
    }
    
    renderPage(pageIndex, containerId) {
        if (pageIndex < 0 || pageIndex >= this.pages.length) {
            throw new Error('Invalid page index');
        }
        
        const container = document.getElementById(containerId);
        if (!container) {
            throw new Error('Container not found');
        }
        
        const page = this.pages[pageIndex];
        const pre = document.createElement('pre');
        pre.className = 'ansi-content';
        
        // This would need full ANSI rendering implementation
        pre.textContent = this.renderPageContent(page);
        
        container.innerHTML = '';
        container.appendChild(pre);
    }
    
    renderPageContent(page) {
        // Placeholder - would need full implementation
        return 'ANSI content rendering would be implemented here';
    }
}

// Export for use
if (typeof module !== 'undefined' && module.exports) {
    module.exports = ANSIRenderer;
} else if (typeof window !== 'undefined') {
    window.ANSIRenderer = ANSIRenderer;
}`
}

func (exp *Exporter) generateWASMLoader(editor *editor.Editor) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s - WebAssembly</title>
    <meta charset="UTF-8">
</head>
<body>
    <div id="wasm-container">
        <h1>%s</h1>
        <p>Loading WebAssembly module...</p>
        <canvas id="ansi-canvas" width="800" height="600"></canvas>
    </div>
    
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("ansi-editor.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        });
    </script>
</body>
</html>`, exp.options.ProjectName, exp.options.ProjectName)
}

func (exp *Exporter) generateWASMBuildScript() string {
	return `#!/bin/bash
# Build script for WebAssembly export
# This would compile the Go code to WebAssembly

echo "Building WebAssembly module..."
GOOS=js GOARCH=wasm go build -o ansi-editor.wasm main.go

echo "Copying wasm_exec.js..."
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

echo "WebAssembly build complete!"
echo "Serve the files with a web server to test."
`
}