package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	
	"github.com/evevioletrose-hash/ansi-editor/pkg/ansi"
	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
	"github.com/evevioletrose-hash/ansi-editor/pkg/editor"
	"github.com/evevioletrose-hash/ansi-editor/pkg/export"
)

func main() {
	var (
		create     = flag.String("create", "", "Create a new ANSI file with specified dimensions (WxH)")
		load       = flag.String("load", "", "Load an ANSI file")
		export_    = flag.String("export", "", "Export project (format:output_dir)")
		info       = flag.String("info", "", "Show information about an ANSI file")
		sample     = flag.String("sample", "", "Create a sample ANSI file")
		convert    = flag.String("convert", "", "Convert ANSI file to plain text")
		dimensions = flag.String("dimensions", "80x25", "Canvas dimensions for new files (WxH)")
		text       = flag.String("text", "", "Text to add when creating")
		color      = flag.String("color", "white", "Text color (white, red, green, blue, etc.)")
		bg         = flag.String("bg", "black", "Background color")
		help       = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()
	
	if *help {
		showHelp()
		return
	}
	
	// Handle different operations
	switch {
	case *create != "":
		handleCreate(*create, *dimensions, *text, *color, *bg)
	case *info != "":
		handleInfo(*info)
	case *sample != "":
		handleSample(*sample)
	case *convert != "":
		handleConvert(*convert)
	case *export_ != "" && *load != "":
		handleExport(*export_, *load)
	case *load != "":
		handleLoad(*load)
	default:
		fmt.Println("ANSI Editor - A tool for creating and editing ANSI art")
		fmt.Println("Use -help for more information")
		showUsageExamples()
	}
}

func showHelp() {
	fmt.Println(`ANSI Editor - A Go-based ANSI artwork editor

USAGE:
    ansi-editor [OPTIONS]

OPTIONS:
    -create <file>      Create a new ANSI file
    -load <file>        Load an existing ANSI file
    -save <file>        Save current work to file (not implemented in CLI)
    -export <format:dir> Export project (json:output, html:output, js:output, wasm:output)
    -info <file>        Show file information
    -sample <file>      Create a sample ANSI file
    -convert <file>     Convert ANSI to plain text
    -dimensions <WxH>   Set canvas dimensions (default: 80x25)
    -text <text>        Add text when creating
    -color <color>      Text color (white, red, green, blue, yellow, magenta, cyan, black)
    -bg <color>         Background color
    -help               Show this help

EXAMPLES:
    ansi-editor -create myart.ans -dimensions 40x20 -text "Hello World" -color red
    ansi-editor -load artwork.ans -save modified.ans
    ansi-editor -sample demo.ans
    ansi-editor -info artwork.ans
    ansi-editor -export json:output_dir -load project.ans
    ansi-editor -convert artwork.ans > plain.txt`)
}

func showUsageExamples() {
	fmt.Println("\nQuick start examples:")
	fmt.Println("  ansi-editor -sample demo.ans              # Create a sample file")
	fmt.Println("  ansi-editor -create myart.ans             # Create a new 80x25 canvas")
	fmt.Println("  ansi-editor -info demo.ans                # Show file information")
	fmt.Println("  ansi-editor -export html:web_output       # Export for web")
}

func handleCreate(filename, dims, text, colorName, bgName string) {
	width, height, err := parseDimensions(dims)
	if err != nil {
		fmt.Printf("Error parsing dimensions: %v\n", err)
		os.Exit(1)
	}
	
	// Create canvas
	c := canvas.NewCanvas(width, height)
	
	// Add text if specified
	if text != "" {
		fgColor := parseColor(colorName)
		bgColor := parseColor(bgName)
		
		// Center the text
		x := (width - len(text)) / 2
		y := height / 2
		
		if x < 0 {
			x = 0
		}
		if y < 0 {
			y = 0
		}
		
		c.DrawText(x, y, text, fgColor, bgColor, 0)
	}
	
	// Save the file
	err = ansi.SaveToFile(c, filename)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Created ANSI file: %s (%dx%d)\n", filename, width, height)
}

func handleLoad(filename string) {
	c, err := ansi.LoadFromFile(filename)
	if err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		os.Exit(1)
	}
	
	// Display the content
	fmt.Printf("Loaded: %s (%dx%d)\n", filename, c.Width, c.Height)
	fmt.Println("Content:")
	fmt.Println(c.RenderANSI())
}

func handleInfo(filename string) {
	info, err := ansi.GetFileInfo(filename)
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("File: %s\n", info.Filename)
	fmt.Printf("Size: %d bytes\n", info.Size)
	fmt.Printf("Dimensions: %dx%d\n", info.Width, info.Height)
	fmt.Printf("Format: %s\n", info.Extension)
}

func handleSample(filename string) {
	err := ansi.CreateSampleFile(filename)
	if err != nil {
		fmt.Printf("Error creating sample file: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Created sample file: %s\n", filename)
}

func handleConvert(filename string) {
	c, err := ansi.LoadFromFile(filename)
	if err != nil {
		fmt.Printf("Error loading file: %v\n", err)
		os.Exit(1)
	}
	
	// Output plain text to stdout
	fmt.Print(c.RenderPlainText())
}

func handleExport(exportSpec, loadFile string) {
	if loadFile == "" {
		fmt.Println("Error: -load is required when using -export")
		os.Exit(1)
	}
	
	// Parse export specification
	parts := strings.Split(exportSpec, ":")
	if len(parts) != 2 {
		fmt.Println("Error: export format should be 'format:output_dir'")
		os.Exit(1)
	}
	
	formatStr, outputDir := parts[0], parts[1]
	
	// Determine export format
	var format export.ExportFormat
	switch strings.ToLower(formatStr) {
	case "json":
		format = export.FormatJSON
	case "html":
		format = export.FormatHTML
	case "js", "javascript":
		format = export.FormatJavaScript
	case "wasm", "webassembly":
		format = export.FormatWebAssembly
	default:
		fmt.Printf("Error: unsupported export format '%s'\n", formatStr)
		os.Exit(1)
	}
	
	// Load the file into an editor
	ed := editor.NewEditor()
	err := ed.LoadPage(loadFile, filepath.Base(loadFile))
	if err != nil {
		fmt.Printf("Error loading file for export: %v\n", err)
		os.Exit(1)
	}
	
	// Set up export options
	options := &export.ExportOptions{
		Format:      format,
		OutputDir:   outputDir,
		ProjectName: strings.TrimSuffix(filepath.Base(loadFile), filepath.Ext(loadFile)),
		Minify:      false,
		Animations:  false,
	}
	
	// Export the project
	exporter := export.NewExporter(options)
	err = exporter.ExportProject(ed)
	if err != nil {
		fmt.Printf("Error exporting project: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Exported project to %s (format: %s)\n", outputDir, formatStr)
}

func parseDimensions(dims string) (int, int, error) {
	parts := strings.Split(dims, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("dimensions must be in format WxH (e.g., 80x25)")
	}
	
	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %s", parts[0])
	}
	
	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %s", parts[1])
	}
	
	if width <= 0 || height <= 0 {
		return 0, 0, fmt.Errorf("dimensions must be positive")
	}
	
	return width, height, nil
}

func parseColor(colorName string) int {
	colors := map[string]int{
		"black":   canvas.ColorBlack,
		"red":     canvas.ColorRed,
		"green":   canvas.ColorGreen,
		"yellow":  canvas.ColorYellow,
		"blue":    canvas.ColorBlue,
		"magenta": canvas.ColorMagenta,
		"cyan":    canvas.ColorCyan,
		"white":   canvas.ColorWhite,
	}
	
	if color, ok := colors[strings.ToLower(colorName)]; ok {
		return color
	}
	
	return canvas.ColorWhite // Default to white
}