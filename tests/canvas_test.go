package tests

import (
	"testing"
	
	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
)

func TestCanvasCreation(t *testing.T) {
	c := canvas.NewCanvas(10, 5)
	
	if c.Width != 10 {
		t.Errorf("Expected width 10, got %d", c.Width)
	}
	
	if c.Height != 5 {
		t.Errorf("Expected height 5, got %d", c.Height)
	}
	
	if len(c.Grid) != 5 {
		t.Errorf("Expected grid height 5, got %d", len(c.Grid))
	}
	
	if len(c.Grid[0]) != 10 {
		t.Errorf("Expected grid width 10, got %d", len(c.Grid[0]))
	}
}

func TestSetGetCell(t *testing.T) {
	c := canvas.NewCanvas(10, 5)
	
	// Test setting a cell
	err := c.SetCell(2, 3, 'A', canvas.ColorRed, canvas.ColorBlue, canvas.AttrBold)
	if err != nil {
		t.Errorf("Unexpected error setting cell: %v", err)
	}
	
	// Test getting the cell
	cell, err := c.GetCell(2, 3)
	if err != nil {
		t.Errorf("Unexpected error getting cell: %v", err)
	}
	
	if cell.Char != 'A' {
		t.Errorf("Expected char 'A', got %c", cell.Char)
	}
	
	if cell.FgColor != canvas.ColorRed {
		t.Errorf("Expected fg color %d, got %d", canvas.ColorRed, cell.FgColor)
	}
	
	if cell.BgColor != canvas.ColorBlue {
		t.Errorf("Expected bg color %d, got %d", canvas.ColorBlue, cell.BgColor)
	}
	
	if cell.Attributes != canvas.AttrBold {
		t.Errorf("Expected attributes %d, got %d", canvas.AttrBold, cell.Attributes)
	}
}

func TestOutOfBounds(t *testing.T) {
	c := canvas.NewCanvas(10, 5)
	
	// Test setting out of bounds
	err := c.SetCell(-1, 0, 'A', canvas.ColorWhite, canvas.ColorBlack, 0)
	if err == nil {
		t.Error("Expected error for negative x coordinate")
	}
	
	err = c.SetCell(0, -1, 'A', canvas.ColorWhite, canvas.ColorBlack, 0)
	if err == nil {
		t.Error("Expected error for negative y coordinate")
	}
	
	err = c.SetCell(10, 0, 'A', canvas.ColorWhite, canvas.ColorBlack, 0)
	if err == nil {
		t.Error("Expected error for x coordinate >= width")
	}
	
	err = c.SetCell(0, 5, 'A', canvas.ColorWhite, canvas.ColorBlack, 0)
	if err == nil {
		t.Error("Expected error for y coordinate >= height")
	}
	
	// Test getting out of bounds
	_, err = c.GetCell(-1, 0)
	if err == nil {
		t.Error("Expected error for negative x coordinate")
	}
	
	_, err = c.GetCell(10, 0)
	if err == nil {
		t.Error("Expected error for x coordinate >= width")
	}
}

func TestDrawText(t *testing.T) {
	c := canvas.NewCanvas(20, 5)
	
	text := "Hello"
	err := c.DrawText(2, 1, text, canvas.ColorGreen, canvas.ColorBlack, 0)
	if err != nil {
		t.Errorf("Unexpected error drawing text: %v", err)
	}
	
	// Check each character
	for i, expectedChar := range text {
		cell, err := c.GetCell(2+i, 1)
		if err != nil {
			t.Errorf("Unexpected error getting cell at (%d, 1): %v", 2+i, err)
		}
		
		if cell.Char != expectedChar {
			t.Errorf("Expected char '%c' at position %d, got '%c'", expectedChar, i, cell.Char)
		}
		
		if cell.FgColor != canvas.ColorGreen {
			t.Errorf("Expected green color at position %d, got %d", i, cell.FgColor)
		}
	}
}

func TestFillRect(t *testing.T) {
	c := canvas.NewCanvas(10, 10)
	
	err := c.FillRect(2, 2, 4, 3, '#', canvas.ColorYellow, canvas.ColorRed, 0)
	if err != nil {
		t.Errorf("Unexpected error filling rectangle: %v", err)
	}
	
	// Check that the rectangle is filled
	for y := 2; y < 5; y++ {
		for x := 2; x < 6; x++ {
			cell, err := c.GetCell(x, y)
			if err != nil {
				t.Errorf("Unexpected error getting cell at (%d, %d): %v", x, y, err)
			}
			
			if cell.Char != '#' {
				t.Errorf("Expected '#' at (%d, %d), got '%c'", x, y, cell.Char)
			}
			
			if cell.FgColor != canvas.ColorYellow {
				t.Errorf("Expected yellow fg at (%d, %d), got %d", x, y, cell.FgColor)
			}
			
			if cell.BgColor != canvas.ColorRed {
				t.Errorf("Expected red bg at (%d, %d), got %d", x, y, cell.BgColor)
			}
		}
	}
}

func TestResize(t *testing.T) {
	c := canvas.NewCanvas(5, 5)
	
	// Set a character in the original canvas
	c.SetCell(2, 2, 'X', canvas.ColorBlue, canvas.ColorWhite, 0)
	
	// Resize to larger
	c.Resize(10, 8)
	
	if c.Width != 10 {
		t.Errorf("Expected width 10 after resize, got %d", c.Width)
	}
	
	if c.Height != 8 {
		t.Errorf("Expected height 8 after resize, got %d", c.Height)
	}
	
	// Check that the original character is preserved
	cell, err := c.GetCell(2, 2)
	if err != nil {
		t.Errorf("Unexpected error getting cell after resize: %v", err)
	}
	
	if cell.Char != 'X' {
		t.Errorf("Expected 'X' to be preserved after resize, got '%c'", cell.Char)
	}
	
	// Check that new areas are initialized with defaults
	cell, err = c.GetCell(9, 7)
	if err != nil {
		t.Errorf("Unexpected error getting new cell after resize: %v", err)
	}
	
	if cell.Char != ' ' {
		t.Errorf("Expected ' ' in new area after resize, got '%c'", cell.Char)
	}
}

func TestClear(t *testing.T) {
	c := canvas.NewCanvas(5, 5)
	
	// Fill with some content
	c.SetCell(2, 2, 'X', canvas.ColorRed, canvas.ColorBlue, canvas.AttrBold)
	
	// Clear the canvas
	c.Clear()
	
	// Check that everything is back to defaults
	cell, err := c.GetCell(2, 2)
	if err != nil {
		t.Errorf("Unexpected error getting cell after clear: %v", err)
	}
	
	if cell.Char != ' ' {
		t.Errorf("Expected ' ' after clear, got '%c'", cell.Char)
	}
	
	if cell.FgColor != 37 {
		t.Errorf("Expected fg color 37 after clear, got %d", cell.FgColor)
	}
	
	if cell.BgColor != 40 {
		t.Errorf("Expected bg color 40 after clear, got %d", cell.BgColor)
	}
	
	if cell.Attributes != 0 {
		t.Errorf("Expected attributes 0 after clear, got %d", cell.Attributes)
	}
}