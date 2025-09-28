package editor

import (
	"fmt"
	"math"
	
	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
)

// DrawingTool represents different drawing tools available in the editor
type DrawingTool int

const (
	ToolPencil DrawingTool = iota
	ToolBrush
	ToolLine
	ToolRectangle
	ToolFillRect
	ToolCircle
	ToolText
	ToolFill
	ToolErase
)

// DrawingState holds the current drawing configuration
type DrawingState struct {
	Tool        DrawingTool
	Character   rune
	FgColor     int
	BgColor     int
	Attributes  int
	BrushSize   int
	Text        string
}

// NewDrawingState creates a new drawing state with defaults
func NewDrawingState() *DrawingState {
	return &DrawingState{
		Tool:       ToolPencil,
		Character:  '█',
		FgColor:    canvas.ColorWhite,
		BgColor:    canvas.ColorBlack,
		Attributes: 0,
		BrushSize:  1,
		Text:       "",
	}
}

// DrawPoint draws a single point at the specified coordinates
func (e *Editor) DrawPoint(x, y int, state *DrawingState) error {
	page := e.GetCurrentPage()
	if page == nil {
		return fmt.Errorf("no active page")
	}
	
	switch state.Tool {
	case ToolPencil:
		return e.drawPencil(x, y, state)
	case ToolBrush:
		return e.drawBrush(x, y, state)
	case ToolErase:
		return e.drawErase(x, y, state)
	case ToolFill:
		return e.drawFill(x, y, state)
	default:
		return fmt.Errorf("unsupported tool for point drawing: %v", state.Tool)
	}
}

// DrawLine draws a line between two points
func (e *Editor) DrawLine(x1, y1, x2, y2 int, state *DrawingState) error {
	page := e.GetCurrentPage()
	if page == nil {
		return fmt.Errorf("no active page")
	}
	
	// Bresenham's line algorithm
	dx := int(math.Abs(float64(x2 - x1)))
	dy := int(math.Abs(float64(y2 - y1)))
	
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	
	err := dx - dy
	x, y := x1, y1
	
	for {
		page.Canvas.SetCell(x, y, state.Character, state.FgColor, state.BgColor, state.Attributes)
		
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
	
	e.MarkCurrentPageModified()
	return nil
}

// DrawRectangle draws a rectangle outline
func (e *Editor) DrawRectangle(x, y, width, height int, state *DrawingState) error {
	page := e.GetCurrentPage()
	if page == nil {
		return fmt.Errorf("no active page")
	}
	
	// Draw horizontal lines
	for i := 0; i < width; i++ {
		page.Canvas.SetCell(x+i, y, state.Character, state.FgColor, state.BgColor, state.Attributes)
		page.Canvas.SetCell(x+i, y+height-1, state.Character, state.FgColor, state.BgColor, state.Attributes)
	}
	
	// Draw vertical lines
	for i := 0; i < height; i++ {
		page.Canvas.SetCell(x, y+i, state.Character, state.FgColor, state.BgColor, state.Attributes)
		page.Canvas.SetCell(x+width-1, y+i, state.Character, state.FgColor, state.BgColor, state.Attributes)
	}
	
	e.MarkCurrentPageModified()
	return nil
}

// DrawFilledRectangle draws a filled rectangle
func (e *Editor) DrawFilledRectangle(x, y, width, height int, state *DrawingState) error {
	page := e.GetCurrentPage()
	if page == nil {
		return fmt.Errorf("no active page")
	}
	
	page.Canvas.FillRect(x, y, width, height, state.Character, state.FgColor, state.BgColor, state.Attributes)
	e.MarkCurrentPageModified()
	return nil
}

// DrawCircle draws a circle outline
func (e *Editor) DrawCircle(centerX, centerY, radius int, state *DrawingState) error {
	page := e.GetCurrentPage()
	if page == nil {
		return fmt.Errorf("no active page")
	}
	
	// Bresenham's circle algorithm
	x := 0
	y := radius
	d := 3 - 2*radius
	
	for x <= y {
		// Draw 8 points of the circle
		e.setCirclePixel(page.Canvas, centerX+x, centerY+y, state)
		e.setCirclePixel(page.Canvas, centerX-x, centerY+y, state)
		e.setCirclePixel(page.Canvas, centerX+x, centerY-y, state)
		e.setCirclePixel(page.Canvas, centerX-x, centerY-y, state)
		e.setCirclePixel(page.Canvas, centerX+y, centerY+x, state)
		e.setCirclePixel(page.Canvas, centerX-y, centerY+x, state)
		e.setCirclePixel(page.Canvas, centerX+y, centerY-x, state)
		e.setCirclePixel(page.Canvas, centerX-y, centerY-x, state)
		
		if d < 0 {
			d += 4*x + 6
		} else {
			d += 4*(x-y) + 10
			y--
		}
		x++
	}
	
	e.MarkCurrentPageModified()
	return nil
}

// DrawText draws text at the specified position
func (e *Editor) DrawText(x, y int, text string, state *DrawingState) error {
	page := e.GetCurrentPage()
	if page == nil {
		return fmt.Errorf("no active page")
	}
	
	page.Canvas.DrawText(x, y, text, state.FgColor, state.BgColor, state.Attributes)
	e.MarkCurrentPageModified()
	return nil
}

// drawPencil implements pencil tool drawing
func (e *Editor) drawPencil(x, y int, state *DrawingState) error {
	page := e.GetCurrentPage()
	err := page.Canvas.SetCell(x, y, state.Character, state.FgColor, state.BgColor, state.Attributes)
	if err == nil {
		e.MarkCurrentPageModified()
	}
	return err
}

// drawBrush implements brush tool drawing with configurable size
func (e *Editor) drawBrush(x, y int, state *DrawingState) error {
	page := e.GetCurrentPage()
	
	size := state.BrushSize
	for dy := -size/2; dy <= size/2; dy++ {
		for dx := -size/2; dx <= size/2; dx++ {
			// Only draw within circular brush
			if dx*dx+dy*dy <= (size/2)*(size/2) {
				page.Canvas.SetCell(x+dx, y+dy, state.Character, state.FgColor, state.BgColor, state.Attributes)
			}
		}
	}
	
	e.MarkCurrentPageModified()
	return nil
}

// drawErase implements eraser tool
func (e *Editor) drawErase(x, y int, state *DrawingState) error {
	page := e.GetCurrentPage()
	err := page.Canvas.SetCell(x, y, ' ', canvas.ColorWhite, canvas.ColorBlack, 0)
	if err == nil {
		e.MarkCurrentPageModified()
	}
	return err
}

// drawFill implements flood fill tool
func (e *Editor) drawFill(x, y int, state *DrawingState) error {
	page := e.GetCurrentPage()
	
	// Get the original cell to determine what to replace
	originalCell, err := page.Canvas.GetCell(x, y)
	if err != nil {
		return err
	}
	
	// Don't fill if the target is already the desired color/char
	if originalCell.Char == state.Character && 
	   originalCell.FgColor == state.FgColor && 
	   originalCell.BgColor == state.BgColor &&
	   originalCell.Attributes == state.Attributes {
		return nil
	}
	
	// Perform flood fill using stack-based approach
	stack := []struct{ x, y int }{{x, y}}
	
	for len(stack) > 0 {
		// Pop from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		
		cx, cy := current.x, current.y
		
		// Check bounds
		if cx < 0 || cx >= page.Canvas.Width || cy < 0 || cy >= page.Canvas.Height {
			continue
		}
		
		// Get current cell
		cell, err := page.Canvas.GetCell(cx, cy)
		if err != nil {
			continue
		}
		
		// Check if this cell matches the original
		if cell.Char != originalCell.Char || 
		   cell.FgColor != originalCell.FgColor || 
		   cell.BgColor != originalCell.BgColor ||
		   cell.Attributes != originalCell.Attributes {
			continue
		}
		
		// Fill this cell
		page.Canvas.SetCell(cx, cy, state.Character, state.FgColor, state.BgColor, state.Attributes)
		
		// Add neighbors to stack
		stack = append(stack, 
			struct{ x, y int }{cx + 1, cy},
			struct{ x, y int }{cx - 1, cy},
			struct{ x, y int }{cx, cy + 1},
			struct{ x, y int }{cx, cy - 1},
		)
	}
	
	e.MarkCurrentPageModified()
	return nil
}

// setCirclePixel helper function for circle drawing
func (e *Editor) setCirclePixel(c *canvas.Canvas, x, y int, state *DrawingState) {
	if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
		c.SetCell(x, y, state.Character, state.FgColor, state.BgColor, state.Attributes)
	}
}