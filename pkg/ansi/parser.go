package ansi

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	
	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
)

// ANSIParser handles parsing of ANSI files and converting them to canvas format
type ANSIParser struct {
	// Current parsing state
	currentFg   int
	currentBg   int
	currentAttr int
	cursorX     int
	cursorY     int
}

// NewANSIParser creates a new ANSI parser with default state
func NewANSIParser() *ANSIParser {
	return &ANSIParser{
		currentFg:   37, // White
		currentBg:   40, // Black
		currentAttr: 0,  // No attributes
		cursorX:     0,
		cursorY:     0,
	}
}

// ParseFromReader parses ANSI content from a reader and returns a canvas
func (p *ANSIParser) ParseFromReader(reader io.Reader) (*canvas.Canvas, error) {
	// First pass: determine canvas dimensions
	width, height, err := p.measureContent(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to measure content: %w", err)
	}
	
	// Reset reader if it's seekable
	if seeker, ok := reader.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}
	
	// Create canvas with measured dimensions
	c := canvas.NewCanvas(width, height)
	
	// Second pass: populate the canvas
	err = p.parseContent(reader, c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content: %w", err)
	}
	
	return c, nil
}

// measureContent determines the required canvas dimensions
func (p *ANSIParser) measureContent(reader io.Reader) (int, int, error) {
	scanner := bufio.NewScanner(reader)
	maxWidth := 0
	height := 0
	currentWidth := 0
	
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[mHJK]`)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Remove ANSI escape sequences for width calculation
		cleanLine := ansiRegex.ReplaceAllString(line, "")
		currentWidth = len([]rune(cleanLine))
		
		if currentWidth > maxWidth {
			maxWidth = currentWidth
		}
		height++
	}
	
	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}
	
	// Ensure minimum dimensions
	if maxWidth == 0 {
		maxWidth = 80
	}
	if height == 0 {
		height = 25
	}
	
	return maxWidth, height, nil
}

// parseContent populates the canvas with parsed ANSI content
func (p *ANSIParser) parseContent(reader io.Reader, c *canvas.Canvas) error {
	p.resetState()
	
	scanner := bufio.NewScanner(reader)
	
	for scanner.Scan() {
		line := scanner.Text()
		err := p.parseLine(line, c)
		if err != nil {
			return err
		}
		
		// Move to next line
		p.cursorY++
		p.cursorX = 0
	}
	
	return scanner.Err()
}

// parseLine processes a single line of ANSI content
func (p *ANSIParser) parseLine(line string, c *canvas.Canvas) error {
	i := 0
	runes := []rune(line)
	
	for i < len(runes) {
		if runes[i] == '\x1b' && i+1 < len(runes) && runes[i+1] == '[' {
			// Parse ANSI escape sequence
			seqEnd := i + 2
			for seqEnd < len(runes) && !isANSITerminator(runes[seqEnd]) {
				seqEnd++
			}
			
			if seqEnd < len(runes) {
				sequence := string(runes[i:seqEnd+1])
				err := p.processANSISequence(sequence)
				if err != nil {
					return err
				}
				i = seqEnd + 1
			} else {
				i++
			}
		} else {
			// Regular character
			if p.cursorY < c.Height && p.cursorX < c.Width {
				c.SetCell(p.cursorX, p.cursorY, runes[i], p.currentFg, p.currentBg, p.currentAttr)
			}
			p.cursorX++
			i++
		}
	}
	
	return nil
}

// processANSISequence handles ANSI escape sequences
func (p *ANSIParser) processANSISequence(sequence string) error {
	if len(sequence) < 3 {
		return nil
	}
	
	command := sequence[len(sequence)-1]
	params := sequence[2 : len(sequence)-1]
	
	switch command {
	case 'm': // SGR (Select Graphic Rendition)
		return p.processSGR(params)
	case 'H', 'f': // Cursor position
		return p.processCursorPosition(params)
	case 'J': // Erase display
		// Ignore for now
	case 'K': // Erase line
		// Ignore for now
	}
	
	return nil
}

// processSGR handles Select Graphic Rendition sequences
func (p *ANSIParser) processSGR(params string) error {
	if params == "" {
		params = "0" // Default to reset
	}
	
	codes := strings.Split(params, ";")
	
	for _, codeStr := range codes {
		code, err := strconv.Atoi(codeStr)
		if err != nil {
			continue // Skip invalid codes
		}
		
		switch {
		case code == 0: // Reset
			p.currentFg = 37
			p.currentBg = 40
			p.currentAttr = 0
		case code >= 1 && code <= 9: // Attributes
			p.currentAttr = code
		case code >= 30 && code <= 37: // Foreground colors
			p.currentFg = code
		case code >= 40 && code <= 47: // Background colors
			p.currentBg = code - 10
		case code >= 90 && code <= 97: // Bright foreground colors
			p.currentFg = code
		case code >= 100 && code <= 107: // Bright background colors
			p.currentBg = code - 10
		}
	}
	
	return nil
}

// processCursorPosition handles cursor positioning sequences
func (p *ANSIParser) processCursorPosition(params string) error {
	if params == "" {
		p.cursorX = 0
		p.cursorY = 0
		return nil
	}
	
	parts := strings.Split(params, ";")
	if len(parts) >= 2 {
		if y, err := strconv.Atoi(parts[0]); err == nil {
			p.cursorY = y - 1 // ANSI is 1-based
		}
		if x, err := strconv.Atoi(parts[1]); err == nil {
			p.cursorX = x - 1 // ANSI is 1-based
		}
	}
	
	return nil
}

// resetState resets the parser to initial state
func (p *ANSIParser) resetState() {
	p.currentFg = 37
	p.currentBg = 40
	p.currentAttr = 0
	p.cursorX = 0
	p.cursorY = 0
}

// isANSITerminator checks if a character terminates an ANSI sequence
func isANSITerminator(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}