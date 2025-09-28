package main

import (
	"fmt"
	"log"
	
	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
	"github.com/evevioletrose-hash/ansi-editor/pkg/editor"
	"github.com/evevioletrose-hash/ansi-editor/pkg/export"
)

// Example showing how to build a Doom-style terminal interface
func main() {
	fmt.Println("=== Doom-Style Terminal Interface Example ===")
	
	// Create editor for game scenes
	gameEditor := editor.NewEditor()
	
	// Create main menu scene
	gameEditor.NewPage("Main Menu", 80, 25)
	createMainMenu(gameEditor)
	
	// Create game HUD scene
	gameEditor.NewPage("Game HUD", 80, 25)
	gameEditor.SetCurrentPage(1) // Switch to HUD page
	createGameHUD(gameEditor)
	
	// Create level scene
	gameEditor.NewPage("Level 1", 80, 25)
	gameEditor.SetCurrentPage(2) // Switch to Level page
	createLevel(gameEditor)
	
	// Display all scenes
	displayScene(gameEditor, 0, "MAIN MENU")
	displayScene(gameEditor, 1, "GAME HUD")
	displayScene(gameEditor, 2, "LEVEL 1")
	
	// Export as web package for browser-based terminal emulation
	exportOptions := &export.ExportOptions{
		Format:      export.FormatHTML,
		OutputDir:   "examples/doom_web",
		ProjectName: "Doom Terminal",
		Minify:      false,
		Animations:  true,
	}
	
	exporter := export.NewExporter(exportOptions)
	err := exporter.ExportProject(gameEditor)
	if err != nil {
		log.Printf("Warning: Could not export web package: %v", err)
	} else {
		fmt.Println("\n🌐 Exported web package to examples/doom_web/")
		fmt.Println("   Open examples/doom_web/index.html in a browser")
	}
	
	// Export as JSON for game engine consumption
	exportOptions.Format = export.FormatJSON
	exportOptions.OutputDir = "examples/doom_data"
	err = exporter.ExportProject(gameEditor)
	if err != nil {
		log.Printf("Warning: Could not export JSON data: %v", err)
	} else {
		fmt.Println("📊 Exported game data to examples/doom_data/")
	}
	
	// Show how this could be used in a game loop
	fmt.Println("\n=== Integration Example ===")
	fmt.Println("This framework can be integrated into game engines like:")
	fmt.Println("• Terminal-based games using ncurses")
	fmt.Println("• Web-based terminal emulators")
	fmt.Println("• Console applications with ANSI support")
	fmt.Println("• WebAssembly games for browsers")
}

func createMainMenu(ed *editor.Editor) {
	ed.SetCurrentPage(0)
	
	// Create title art
	titleState := &editor.DrawingState{
		Character:  ' ',
		FgColor:    canvas.ColorYellow,
		BgColor:    canvas.ColorRed,
		Attributes: canvas.AttrBold,
	}
	
	// Draw title background
	ed.DrawFilledRectangle(20, 3, 40, 7, titleState)
	
	// Draw title text
	titleState.FgColor = canvas.ColorWhite
	titleState.BgColor = canvas.ColorRed
	ed.DrawText(35, 5, "D O O M", titleState)
	ed.DrawText(32, 6, "TERMINAL EDITION", titleState)
	
	// Menu options
	menuState := &editor.DrawingState{
		Character:  ' ',
		FgColor:    canvas.ColorWhite,
		BgColor:    canvas.ColorBlack,
		Attributes: 0,
	}
	
	ed.DrawText(35, 12, "> NEW GAME", menuState)
	ed.DrawText(35, 14, "  LOAD GAME", menuState)
	ed.DrawText(35, 16, "  OPTIONS", menuState)
	ed.DrawText(35, 18, "  QUIT", menuState)
	
	// Draw border
	borderState := &editor.DrawingState{
		Character: '█',
		FgColor:   canvas.ColorRed,
		BgColor:   canvas.ColorBlack,
	}
	ed.DrawRectangle(15, 1, 50, 22, borderState)
}

func createGameHUD(ed *editor.Editor) {
	// Health bar
	healthState := &editor.DrawingState{
		Character: '█',
		FgColor:   canvas.ColorGreen,
		BgColor:   canvas.ColorBlack,
	}
	
	ed.DrawText(2, 22, "HEALTH:", &editor.DrawingState{
		FgColor: canvas.ColorWhite,
		BgColor: canvas.ColorBlack,
	})
	
	// Full health bar (100%)
	for i := 0; i < 20; i++ {
		ed.DrawPoint(10+i, 22, healthState)
	}
	
	// Ammo counter
	ed.DrawText(40, 22, "AMMO: 150/200", &editor.DrawingState{
		FgColor: canvas.ColorYellow,
		BgColor: canvas.ColorBlack,
	})
	
	// Weapon display
	ed.DrawText(2, 23, "WEAPON: SHOTGUN", &editor.DrawingState{
		FgColor: canvas.ColorCyan,
		BgColor: canvas.ColorBlack,
	})
	
	// Armor
	ed.DrawText(40, 23, "ARMOR: 75%", &editor.DrawingState{
		FgColor: canvas.ColorBlue,
		BgColor: canvas.ColorBlack,
	})
	
	// Game view area (simplified 3D representation)
	viewState := &editor.DrawingState{
		Character: '▓',
		FgColor:   canvas.ColorBlack,
		BgColor:   canvas.ColorWhite,
	}
	
	// Draw "walls" using ASCII art technique
	for y := 2; y < 20; y++ {
		// Left wall
		for x := 10; x < 15; x++ {
			ed.DrawPoint(x, y, viewState)
		}
		// Right wall  
		for x := 65; x < 70; x++ {
			ed.DrawPoint(x, y, viewState)
		}
	}
	
	// Floor pattern
	floorState := &editor.DrawingState{
		Character: '.',
		FgColor:   canvas.ColorBlack,
		BgColor:   canvas.ColorWhite,
	}
	
	for y := 15; y < 20; y++ {
		for x := 15; x < 65; x += 2 {
			ed.DrawPoint(x, y, floorState)
		}
	}
	
	// "Enemy" representation
	enemyState := &editor.DrawingState{
		Character: 'E',
		FgColor:   canvas.ColorRed,
		BgColor:   canvas.ColorWhite,
		Attributes: canvas.AttrBold,
	}
	ed.DrawPoint(40, 10, enemyState)
}

func createLevel(ed *editor.Editor) {
	// Create a simple level map using ASCII characters
	
	// Walls
	wallState := &editor.DrawingState{
		Character: '█',
		FgColor:   canvas.ColorWhite,
		BgColor:   canvas.ColorBlack,
	}
	
	// Top and bottom walls
	for x := 0; x < 80; x++ {
		ed.DrawPoint(x, 0, wallState)
		ed.DrawPoint(x, 24, wallState)
	}
	
	// Left and right walls
	for y := 0; y < 25; y++ {
		ed.DrawPoint(0, y, wallState)
		ed.DrawPoint(79, y, wallState)
	}
	
	// Internal walls to create rooms
	for y := 5; y < 20; y++ {
		ed.DrawPoint(20, y, wallState)
		ed.DrawPoint(60, y, wallState)
	}
	
	// Doors (gaps in walls)
	ed.DrawPoint(20, 12, &editor.DrawingState{Character: ' ', FgColor: canvas.ColorBlack, BgColor: canvas.ColorBlack})
	ed.DrawPoint(60, 8, &editor.DrawingState{Character: ' ', FgColor: canvas.ColorBlack, BgColor: canvas.ColorBlack})
	
	// Player position
	playerState := &editor.DrawingState{
		Character:  '@',
		FgColor:    canvas.ColorGreen,
		BgColor:    canvas.ColorBlack,
		Attributes: canvas.AttrBold,
	}
	ed.DrawPoint(10, 12, playerState)
	
	// Enemies
	enemyState := &editor.DrawingState{
		Character:  'D',
		FgColor:    canvas.ColorRed,
		BgColor:    canvas.ColorBlack,
		Attributes: canvas.AttrBold,
	}
	ed.DrawPoint(30, 8, enemyState)
	ed.DrawPoint(50, 15, enemyState)
	ed.DrawPoint(70, 10, enemyState)
	
	// Items
	itemState := &editor.DrawingState{
		Character: '*',
		FgColor:   canvas.ColorYellow,
		BgColor:   canvas.ColorBlack,
	}
	ed.DrawPoint(15, 6, itemState)  // Health pack
	ed.DrawPoint(45, 18, itemState) // Ammo
	ed.DrawPoint(65, 5, itemState)  // Weapon
	
	// Exit
	exitState := &editor.DrawingState{
		Character:  'X',
		FgColor:    canvas.ColorCyan,
		BgColor:    canvas.ColorBlack,
		Attributes: canvas.AttrBold,
	}
	ed.DrawPoint(70, 20, exitState)
	
	// Legend
	ed.DrawText(2, 1, "@ = Player", &editor.DrawingState{FgColor: canvas.ColorGreen, BgColor: canvas.ColorBlack})
	ed.DrawText(2, 2, "D = Demon", &editor.DrawingState{FgColor: canvas.ColorRed, BgColor: canvas.ColorBlack})
	ed.DrawText(15, 1, "* = Item", &editor.DrawingState{FgColor: canvas.ColorYellow, BgColor: canvas.ColorBlack})
	ed.DrawText(15, 2, "X = Exit", &editor.DrawingState{FgColor: canvas.ColorCyan, BgColor: canvas.ColorBlack})
}

func displayScene(ed *editor.Editor, pageIndex int, title string) {
	ed.SetCurrentPage(pageIndex)
	page := ed.GetCurrentPage()
	
	fmt.Printf("\n=== %s ===\n", title)
	fmt.Println(page.Canvas.RenderANSI())
}