// Simple Game Example - Demonstrates basic game mechanics with ANSI Editor
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
	"github.com/evevioletrose-hash/ansi-editor/pkg/export"
	"github.com/evevioletrose-hash/ansi-editor/pkg/objects"
)

func main() {
	fmt.Println("Creating Simple Game Example...")

	// Create the game scene
	scene := objects.NewScene("simple-game", 60, 20)

	// Set up the background
	setupBackground(scene)

	// Create player
	player := createPlayer(scene)

	// Create some enemies
	createEnemies(scene)

	// Create collectibles
	createCollectibles(scene)

	// Create UI elements
	createUI(scene)

	// Display the game
	fmt.Println("\nGame Scene:")
	rendered := scene.Render()
	fmt.Print(rendered.ToANSI())

	// Demonstrate some gameplay
	fmt.Println("\nSimulating gameplay...")
	simulateGameplay(scene, player)

	// Export as WASM
	fmt.Println("\nExporting game as WASM bundle...")
	err := export.ExportScene(scene, "./simple-game-export", "simple-game", "1.0.0")
	if err != nil {
		fmt.Printf("Export failed: %v\n", err)
		return
	}

	fmt.Println("Game exported successfully!")
	fmt.Println("\nTo play in browser:")
	fmt.Println("  cd simple-game-export")
	fmt.Println("  ./build.sh")
	fmt.Println("  python3 -m http.server 8080")
	fmt.Println("  # Open http://localhost:8080")
}

func setupBackground(scene *objects.Scene) {
	// Fill with grass
	scene.Background.FillRect(0, 0, 60, 20, '.', canvas.Green, canvas.Black)

	// Add some trees
	for i := 0; i < 8; i++ {
		x := rand.Intn(50) + 5
		y := rand.Intn(15) + 2
		scene.Background.DrawChar(x, y, 'T', canvas.BrightGreen, canvas.Black)
	}

	// Add some rocks
	for i := 0; i < 12; i++ {
		x := rand.Intn(55) + 2
		y := rand.Intn(17) + 1
		scene.Background.DrawChar(x, y, 'o', canvas.BrightBlack, canvas.Black)
	}

	// Create borders
	for x := 0; x < 60; x++ {
		scene.Background.DrawChar(x, 0, '#', canvas.Yellow, canvas.Black)
		scene.Background.DrawChar(x, 19, '#', canvas.Yellow, canvas.Black)
	}
	for y := 0; y < 20; y++ {
		scene.Background.DrawChar(0, y, '#', canvas.Yellow, canvas.Black)
		scene.Background.DrawChar(59, y, '#', canvas.Yellow, canvas.Black)
	}
}

func createPlayer(scene *objects.Scene) *objects.Object {
	player := objects.NewObject("player", "Player", 30, 10, 1, 1)
	
	// Create player animation frames (walking)
	frames := []struct {
		char rune
		color canvas.ANSIColor
	}{
		{'@', canvas.BrightYellow},
		{'&', canvas.Yellow},
		{'@', canvas.BrightYellow},
		{'%', canvas.Yellow},
	}

	for _, frame := range frames {
		frameCanvas := canvas.NewCanvas(1, 1)
		frameCanvas.DrawChar(0, 0, frame.char, frame.color, canvas.Black)
		player.AddFrame(frameCanvas)
	}

	scene.AddObject(player)
	return player
}

func createEnemies(scene *objects.Scene) {
	enemies := []struct {
		id   string
		x, y int
		char rune
	}{
		{"enemy1", 10, 5, 'X'},
		{"enemy2", 45, 8, 'Z'},
		{"enemy3", 25, 15, 'X'},
		{"enemy4", 50, 12, 'Z'},
	}

	for _, e := range enemies {
		enemy := objects.NewObject(e.id, "Enemy", e.x, e.y, 1, 1)
		
		// Create enemy animation
		frame1 := canvas.NewCanvas(1, 1)
		frame1.DrawChar(0, 0, e.char, canvas.Red, canvas.Black)
		
		frame2 := canvas.NewCanvas(1, 1)
		frame2.DrawChar(0, 0, e.char, canvas.BrightRed, canvas.Black)

		enemy.AddFrame(frame1)
		enemy.AddFrame(frame2)
		
		scene.AddObject(enemy)
	}
}

func createCollectibles(scene *objects.Scene) {
	collectibles := []struct {
		id   string
		x, y int
		char rune
	}{
		{"coin1", 15, 7, '$'},
		{"coin2", 35, 12, '$'},
		{"coin3", 20, 16, '*'},
		{"coin4", 40, 6, '*'},
		{"coin5", 12, 14, '$'},
	}

	for _, c := range collectibles {
		coin := objects.NewObject(c.id, "Collectible", c.x, c.y, 1, 1)
		
		// Glittering animation
		frame1 := canvas.NewCanvas(1, 1)
		frame1.DrawChar(0, 0, c.char, canvas.Yellow, canvas.Black)
		
		frame2 := canvas.NewCanvas(1, 1)
		frame2.DrawChar(0, 0, c.char, canvas.BrightYellow, canvas.Black)
		
		frame3 := canvas.NewCanvas(1, 1)
		frame3.DrawChar(0, 0, c.char, canvas.White, canvas.Black)

		coin.AddFrame(frame1)
		coin.AddFrame(frame2)
		coin.AddFrame(frame3)
		coin.AddFrame(frame2)
		
		scene.AddObject(coin)
	}
}

func createUI(scene *objects.Scene) {
	// Status bar at top
	statusBar := objects.NewObject("status", "Status Bar", 1, 1, 58, 1)
	statusBar.Canvas.FillRect(0, 0, 58, 1, ' ', canvas.Black, canvas.BrightWhite)
	statusBar.Canvas.DrawString(2, 0, "Score: 0  Lives: 3  Level: 1", canvas.Black, canvas.BrightWhite)
	statusBar.ZIndex = 100 // Always on top
	scene.AddObject(statusBar)

	// Instructions at bottom
	instructions := objects.NewObject("instructions", "Instructions", 1, 18, 58, 1)
	instructions.Canvas.FillRect(0, 0, 58, 1, ' ', canvas.Black, canvas.Blue)
	instructions.Canvas.DrawString(2, 0, "Arrow keys to move, SPACE to jump, ESC to quit", canvas.White, canvas.Blue)
	instructions.ZIndex = 100
	scene.AddObject(instructions)
}

func simulateGameplay(scene *objects.Scene, player *objects.Object) {
	// Simulate a few seconds of gameplay
	for i := 0; i < 20; i++ {
		// Update animations
		scene.UpdateAnimations()

		// Move player in a simple pattern
		switch i % 8 {
		case 0, 4:
			player.Move(player.X+1, player.Y)
		case 1, 5:
			player.Move(player.X, player.Y-1)
		case 2, 6:
			player.Move(player.X-1, player.Y)
		case 3, 7:
			player.Move(player.X, player.Y+1)
		}

		// Move some enemies
		if enemy, exists := scene.GetObject("enemy1"); exists {
			if i%4 == 0 {
				newX := enemy.X + (rand.Intn(3) - 1) // -1, 0, or 1
				newY := enemy.Y + (rand.Intn(3) - 1)
				if newX > 1 && newX < 58 && newY > 1 && newY < 18 {
					enemy.Move(newX, newY)
				}
			}
		}

		// Simple frame display (in a real game, you'd update the display)
		if i%5 == 0 {
			fmt.Printf("Frame %d: Player at (%d, %d)\n", i, player.X, player.Y)
		}

		// Small delay
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Gameplay simulation complete.")
}