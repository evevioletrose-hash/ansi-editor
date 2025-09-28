package ansi

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
)

// SupportedFormats lists the ANSI file formats we support
var SupportedFormats = []string{".ans", ".ansi", ".txt"}

// LoadFromFile loads an ANSI file and returns a canvas
func LoadFromFile(filename string) (*canvas.Canvas, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()
	
	parser := NewANSIParser()
	return parser.ParseFromReader(file)
}

// SaveToFile saves a canvas to an ANSI file
func SaveToFile(c *canvas.Canvas, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()
	
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	
	return SaveToWriter(c, writer)
}

// SaveToWriter saves a canvas to a writer in ANSI format
func SaveToWriter(c *canvas.Canvas, writer *bufio.Writer) error {
	// Write ANSI content
	ansiContent := c.RenderANSI()
	_, err := writer.WriteString(ansiContent)
	if err != nil {
		return fmt.Errorf("failed to write ANSI content: %w", err)
	}
	
	return nil
}

// SaveAsPlainText saves a canvas as plain text (no ANSI codes)
func SaveAsPlainText(c *canvas.Canvas, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()
	
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	
	plainContent := c.RenderPlainText()
	_, err = writer.WriteString(plainContent)
	if err != nil {
		return fmt.Errorf("failed to write plain text content: %w", err)
	}
	
	return nil
}

// IsValidFormat checks if a filename has a supported ANSI format extension
func IsValidFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, format := range SupportedFormats {
		if ext == format {
			return true
		}
	}
	return false
}

// GetFileInfo returns basic information about an ANSI file
func GetFileInfo(filename string) (*FileInfo, error) {
	if !IsValidFormat(filename) {
		return nil, fmt.Errorf("unsupported file format: %s", filepath.Ext(filename))
	}
	
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	
	// Try to parse the file to get dimensions
	c, err := LoadFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file for info: %w", err)
	}
	
	return &FileInfo{
		Filename:   filename,
		Size:       stat.Size(),
		ModTime:    stat.ModTime(),
		Width:      c.Width,
		Height:     c.Height,
		Extension:  filepath.Ext(filename),
	}, nil
}

// FileInfo contains metadata about an ANSI file
type FileInfo struct {
	Filename  string
	Size      int64
	ModTime   interface{}
	Width     int
	Height    int
	Extension string
}

// CreateSampleFile creates a sample ANSI file for testing
func CreateSampleFile(filename string) error {
	c := canvas.NewCanvas(40, 10)
	
	// Add some sample content
	c.DrawText(2, 1, "ANSI Editor Sample", canvas.ColorWhite, canvas.ColorBlue, canvas.AttrBold)
	c.DrawText(2, 3, "This is a test file", canvas.ColorYellow, canvas.ColorBlack, 0)
	c.DrawText(2, 4, "with colors and text!", canvas.ColorGreen, canvas.ColorBlack, 0)
	
	// Add a border
	for x := 0; x < c.Width; x++ {
		c.SetCell(x, 0, '═', canvas.ColorCyan, canvas.ColorBlack, 0)
		c.SetCell(x, c.Height-1, '═', canvas.ColorCyan, canvas.ColorBlack, 0)
	}
	for y := 0; y < c.Height; y++ {
		c.SetCell(0, y, '║', canvas.ColorCyan, canvas.ColorBlack, 0)
		c.SetCell(c.Width-1, y, '║', canvas.ColorCyan, canvas.ColorBlack, 0)
	}
	
	// Corners
	c.SetCell(0, 0, '╔', canvas.ColorCyan, canvas.ColorBlack, 0)
	c.SetCell(c.Width-1, 0, '╗', canvas.ColorCyan, canvas.ColorBlack, 0)
	c.SetCell(0, c.Height-1, '╚', canvas.ColorCyan, canvas.ColorBlack, 0)
	c.SetCell(c.Width-1, c.Height-1, '╝', canvas.ColorCyan, canvas.ColorBlack, 0)
	
	return SaveToFile(c, filename)
}