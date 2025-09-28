package main

import (
	"fmt"
	"log"
	
	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
	"github.com/evevioletrose-hash/ansi-editor/pkg/editor"
	"github.com/evevioletrose-hash/ansi-editor/pkg/export"
)

func main() {
	// Create a new editor
	ed := editor.NewEditor()
	
	// Create a new page
	page := ed.NewPage("Demo Page", 60, 20)
	
	// Get drawing state
	drawState := editor.NewDrawingState()
	drawState.Character = '█'
	drawState.FgColor = canvas.ColorRed
	drawState.BgColor = canvas.ColorBlack
	
	// Draw a title
	ed.DrawText(5, 2, "ANSI Editor Demo", &editor.DrawingState{
		Character:  ' ',
		FgColor:    canvas.ColorYellow,
		BgColor:    canvas.ColorBlue,
		Attributes: canvas.AttrBold,
	})
	
	// Draw a border
	err := ed.DrawRectangle(2, 4, 50, 10, &editor.DrawingState{
		Character: '█',
		FgColor:   canvas.ColorCyan,
		BgColor:   canvas.ColorBlack,
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// Fill an area
	err = ed.DrawFilledRectangle(10, 6, 20, 4, &editor.DrawingState{
		Character: '▓',
		FgColor:   canvas.ColorGreen,
		BgColor:   canvas.ColorBlack,
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// Draw some text inside
	ed.DrawText(15, 8, "Hello ANSI!", &editor.DrawingState{
		Character:  ' ',
		FgColor:    canvas.ColorWhite,
		BgColor:    canvas.ColorGreen,
		Attributes: canvas.AttrBold,
	})
	
	// Draw a circle
	err = ed.DrawCircle(40, 12, 5, &editor.DrawingState{
		Character: '●',
		FgColor:   canvas.ColorMagenta,
		BgColor:   canvas.ColorBlack,
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// Display the result
	fmt.Println("=== ANSI Editor Demo ===")
	fmt.Println(page.Canvas.RenderANSI())
	
	// Save the page
	err = ed.SaveCurrentPage("examples/demo.ans")
	if err != nil {
		log.Printf("Warning: Could not save demo file: %v", err)
	} else {
		fmt.Println("\nSaved demo to examples/demo.ans")
	}
	
	// Export as HTML
	options := &export.ExportOptions{
		Format:      export.FormatHTML,
		OutputDir:   "examples/html_output",
		ProjectName: "ANSI Demo",
		Minify:      false,
		Animations:  false,
	}
	
	exporter := export.NewExporter(options)
	err = exporter.ExportProject(ed)
	if err != nil {
		log.Printf("Warning: Could not export to HTML: %v", err)
	} else {
		fmt.Println("Exported to examples/html_output/")
	}
	
	// Show page information
	fmt.Println("\n=== Page Information ===")
	pageInfo := ed.GetPageInfo()
	for _, info := range pageInfo {
		fmt.Printf("Page %d: %s (%dx%d) Modified: %v Active: %v\n",
			info.Index, info.Name, info.Width, info.Height, info.Modified, info.Active)
	}
}