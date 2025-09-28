// Package canvas provides a drawing surface for ANSI art with support for
// colors, characters, and basic drawing operations.
package canvas

import (
	"fmt"
	"strings"
)

// ANSIColor represents ANSI color codes
type ANSIColor int

const (
	Black ANSIColor = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

// Cell represents a single character cell in the canvas
type Cell struct {
	Char       rune      `json:"char"`
	Foreground ANSIColor `json:"foreground"`
	Background ANSIColor `json:"background"`
	Bold       bool      `json:"bold"`
	Italic     bool      `json:"italic"`
	Underline  bool      `json:"underline"`
}

// Canvas represents a 2D grid of ANSI characters
type Canvas struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Cells  []Cell `json:"cells"`
}

// NewCanvas creates a new canvas with the specified dimensions
func NewCanvas(width, height int) *Canvas {
	cells := make([]Cell, width*height)
	// Initialize with spaces and default colors
	for i := range cells {
		cells[i] = Cell{
			Char:       ' ',
			Foreground: White,
			Background: Black,
		}
	}
	
	return &Canvas{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// GetCell returns the cell at the specified coordinates
func (c *Canvas) GetCell(x, y int) (*Cell, error) {
	if x < 0 || x >= c.Width || y < 0 || y >= c.Height {
		return nil, fmt.Errorf("coordinates (%d, %d) out of bounds", x, y)
	}
	return &c.Cells[y*c.Width+x], nil
}

// SetCell sets the cell at the specified coordinates
func (c *Canvas) SetCell(x, y int, cell Cell) error {
	if x < 0 || x >= c.Width || y < 0 || y >= c.Height {
		return fmt.Errorf("coordinates (%d, %d) out of bounds", x, y)
	}
	c.Cells[y*c.Width+x] = cell
	return nil
}

// DrawChar draws a character at the specified position with given colors
func (c *Canvas) DrawChar(x, y int, char rune, fg, bg ANSIColor) error {
	return c.SetCell(x, y, Cell{
		Char:       char,
		Foreground: fg,
		Background: bg,
	})
}

// DrawString draws a string starting at the specified position
func (c *Canvas) DrawString(x, y int, text string, fg, bg ANSIColor) error {
	for i, char := range text {
		if x+i >= c.Width {
			break // Stop if we exceed canvas width
		}
		if err := c.DrawChar(x+i, y, char, fg, bg); err != nil {
			return err
		}
	}
	return nil
}

// Clear fills the entire canvas with the specified character and colors
func (c *Canvas) Clear(char rune, fg, bg ANSIColor) {
	cell := Cell{
		Char:       char,
		Foreground: fg,
		Background: bg,
	}
	for i := range c.Cells {
		c.Cells[i] = cell
	}
}

// DrawLine draws a line between two points using the specified character
func (c *Canvas) DrawLine(x1, y1, x2, y2 int, char rune, fg, bg ANSIColor) {
	// Simple Bresenham's line algorithm
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := sign(x2 - x1)
	sy := sign(y2 - y1)
	err := dx - dy

	x, y := x1, y1
	for {
		c.DrawChar(x, y, char, fg, bg)
		if x == x2 && y == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// DrawRect draws a rectangle outline
func (c *Canvas) DrawRect(x, y, width, height int, char rune, fg, bg ANSIColor) {
	// Top and bottom lines
	for i := 0; i < width; i++ {
		c.DrawChar(x+i, y, char, fg, bg)
		c.DrawChar(x+i, y+height-1, char, fg, bg)
	}
	// Left and right lines
	for i := 0; i < height; i++ {
		c.DrawChar(x, y+i, char, fg, bg)
		c.DrawChar(x+width-1, y+i, char, fg, bg)
	}
}

// FillRect fills a rectangle with the specified character
func (c *Canvas) FillRect(x, y, width, height int, char rune, fg, bg ANSIColor) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			c.DrawChar(x+dx, y+dy, char, fg, bg)
		}
	}
}

// ToANSI converts the canvas to ANSI escape sequence string
func (c *Canvas) ToANSI() string {
	var result strings.Builder
	
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			cell := c.Cells[y*c.Width+x]
			
			// Set foreground color
			result.WriteString(fmt.Sprintf("\033[%dm", 30+int(cell.Foreground)))
			
			// Set background color
			result.WriteString(fmt.Sprintf("\033[%dm", 40+int(cell.Background)))
			
			// Set text attributes
			if cell.Bold {
				result.WriteString("\033[1m")
			}
			if cell.Italic {
				result.WriteString("\033[3m")
			}
			if cell.Underline {
				result.WriteString("\033[4m")
			}
			
			// Write the character
			result.WriteRune(cell.Char)
			
			// Reset attributes
			result.WriteString("\033[0m")
		}
		result.WriteString("\n")
	}
	
	return result.String()
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}