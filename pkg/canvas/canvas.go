package canvas

import (
	"fmt"
)

// Cell represents a single character cell in the canvas with color information
type Cell struct {
	Char       rune   // The character to display
	FgColor    int    // Foreground color (ANSI color code)
	BgColor    int    // Background color (ANSI color code)
	Attributes int    // Text attributes (bold, italic, etc.)
}

// Canvas represents a grid-based drawing surface for ANSI art
type Canvas struct {
	Width  int     // Canvas width in characters
	Height int     // Canvas height in characters
	Grid   [][]Cell // 2D grid of cells
}

// NewCanvas creates a new canvas with the specified dimensions
func NewCanvas(width, height int) *Canvas {
	grid := make([][]Cell, height)
	for i := range grid {
		grid[i] = make([]Cell, width)
		// Initialize with default values (space character, default colors)
		for j := range grid[i] {
			grid[i][j] = Cell{
				Char:       ' ',
				FgColor:    37, // White foreground
				BgColor:    40, // Black background
				Attributes: 0,  // No attributes
			}
		}
	}
	
	return &Canvas{
		Width:  width,
		Height: height,
		Grid:   grid,
	}
}

// SetCell sets a character and its properties at the specified position
func (c *Canvas) SetCell(x, y int, char rune, fgColor, bgColor, attributes int) error {
	if x < 0 || x >= c.Width || y < 0 || y >= c.Height {
		return fmt.Errorf("position (%d, %d) is out of bounds", x, y)
	}
	
	c.Grid[y][x] = Cell{
		Char:       char,
		FgColor:    fgColor,
		BgColor:    bgColor,
		Attributes: attributes,
	}
	
	return nil
}

// GetCell returns the cell at the specified position
func (c *Canvas) GetCell(x, y int) (Cell, error) {
	if x < 0 || x >= c.Width || y < 0 || y >= c.Height {
		return Cell{}, fmt.Errorf("position (%d, %d) is out of bounds", x, y)
	}
	
	return c.Grid[y][x], nil
}

// Clear resets all cells to default values
func (c *Canvas) Clear() {
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			c.Grid[y][x] = Cell{
				Char:       ' ',
				FgColor:    37,
				BgColor:    40,
				Attributes: 0,
			}
		}
	}
}

// Resize changes the canvas dimensions, preserving existing content where possible
func (c *Canvas) Resize(newWidth, newHeight int) {
	newGrid := make([][]Cell, newHeight)
	
	for i := 0; i < newHeight; i++ {
		newGrid[i] = make([]Cell, newWidth)
		
		for j := 0; j < newWidth; j++ {
			if i < c.Height && j < c.Width {
				// Copy existing cell
				newGrid[i][j] = c.Grid[i][j]
			} else {
				// Initialize new cell with defaults
				newGrid[i][j] = Cell{
					Char:       ' ',
					FgColor:    37,
					BgColor:    40,
					Attributes: 0,
				}
			}
		}
	}
	
	c.Width = newWidth
	c.Height = newHeight
	c.Grid = newGrid
}

// DrawText draws text starting at the specified position with given colors
func (c *Canvas) DrawText(x, y int, text string, fgColor, bgColor, attributes int) error {
	if y < 0 || y >= c.Height {
		return fmt.Errorf("y position %d is out of bounds", y)
	}
	
	for i, char := range text {
		if x+i >= c.Width {
			break // Stop if we reach the edge
		}
		if x+i >= 0 {
			c.Grid[y][x+i] = Cell{
				Char:       char,
				FgColor:    fgColor,
				BgColor:    bgColor,
				Attributes: attributes,
			}
		}
	}
	
	return nil
}

// FillRect fills a rectangular area with the specified character and colors
func (c *Canvas) FillRect(x, y, width, height int, char rune, fgColor, bgColor, attributes int) error {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			px, py := x+dx, y+dy
			if px >= 0 && px < c.Width && py >= 0 && py < c.Height {
				c.Grid[py][px] = Cell{
					Char:       char,
					FgColor:    fgColor,
					BgColor:    bgColor,
					Attributes: attributes,
				}
			}
		}
	}
	
	return nil
}