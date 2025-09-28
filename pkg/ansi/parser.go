// Package ansi provides functionality for parsing and generating ANSI art files
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

// ANSI escape sequence patterns
var (
	escapePattern = regexp.MustCompile(`\033\[([0-9;]+)m`)
	cursorPattern = regexp.MustCompile(`\033\[(\d+);(\d+)H`)
)

// Parser handles ANSI file parsing
type Parser struct {
	canvas *canvas.Canvas
	x, y   int
	fg, bg canvas.ANSIColor
	bold   bool
	italic bool
	underline bool
}

// NewParser creates a new ANSI parser
func NewParser() *Parser {
	return &Parser{
		fg: canvas.White,
		bg: canvas.Black,
	}
}

// ParseFile parses an ANSI file and returns a canvas
func (p *Parser) ParseFile(reader io.Reader, width, height int) (*canvas.Canvas, error) {
	p.canvas = canvas.NewCanvas(width, height)
	p.x, p.y = 0, 0
	p.fg, p.bg = canvas.White, canvas.Black
	p.bold, p.italic, p.underline = false, false, false

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if err := p.parseLine(line); err != nil {
			return nil, fmt.Errorf("error parsing line: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return p.canvas, nil
}

// parseLine processes a single line of ANSI text
func (p *Parser) parseLine(line string) error {
	i := 0
	for i < len(line) {
		if line[i] == '\033' && i+1 < len(line) && line[i+1] == '[' {
			// Found escape sequence
			end := strings.Index(line[i:], "m")
			if end == -1 {
				// Check for cursor positioning
				cursorEnd := strings.Index(line[i:], "H")
				if cursorEnd != -1 {
					if err := p.processCursorSequence(line[i : i+cursorEnd+1]); err != nil {
						return err
					}
					i += cursorEnd + 1
					continue
				}
				// No valid sequence found, treat as regular character
				p.drawChar(rune(line[i]))
				i++
				continue
			}
			
			sequence := line[i : i+end+1]
			if err := p.processEscapeSequence(sequence); err != nil {
				return err
			}
			i += end + 1
		} else {
			// Regular character
			p.drawChar(rune(line[i]))
			i++
		}
	}
	
	// Move to next line
	p.x = 0
	p.y++
	
	return nil
}

// processEscapeSequence handles ANSI escape sequences
func (p *Parser) processEscapeSequence(sequence string) error {
	matches := escapePattern.FindStringSubmatch(sequence)
	if len(matches) != 2 {
		return fmt.Errorf("invalid escape sequence: %s", sequence)
	}

	codes := strings.Split(matches[1], ";")
	for _, codeStr := range codes {
		code, err := strconv.Atoi(codeStr)
		if err != nil {
			return fmt.Errorf("invalid color code: %s", codeStr)
		}

		switch {
		case code == 0:
			// Reset all attributes
			p.fg, p.bg = canvas.White, canvas.Black
			p.bold, p.italic, p.underline = false, false, false
		case code == 1:
			p.bold = true
		case code == 3:
			p.italic = true
		case code == 4:
			p.underline = true
		case code >= 30 && code <= 37:
			// Standard foreground colors
			p.fg = canvas.ANSIColor(code - 30)
		case code >= 40 && code <= 47:
			// Standard background colors
			p.bg = canvas.ANSIColor(code - 40)
		case code >= 90 && code <= 97:
			// Bright foreground colors
			p.fg = canvas.ANSIColor(code - 90 + 8)
		case code >= 100 && code <= 107:
			// Bright background colors
			p.bg = canvas.ANSIColor(code - 100 + 8)
		}
	}

	return nil
}

// processCursorSequence handles cursor positioning sequences
func (p *Parser) processCursorSequence(sequence string) error {
	matches := cursorPattern.FindStringSubmatch(sequence)
	if len(matches) != 3 {
		return fmt.Errorf("invalid cursor sequence: %s", sequence)
	}

	y, err := strconv.Atoi(matches[1])
	if err != nil {
		return err
	}
	x, err := strconv.Atoi(matches[2])
	if err != nil {
		return err
	}

	// ANSI coordinates are 1-based, convert to 0-based
	p.x = x - 1
	p.y = y - 1

	return nil
}

// drawChar draws a character at the current position with current attributes
func (p *Parser) drawChar(char rune) {
	if p.x >= 0 && p.x < p.canvas.Width && p.y >= 0 && p.y < p.canvas.Height {
		cell := canvas.Cell{
			Char:       char,
			Foreground: p.fg,
			Background: p.bg,
			Bold:       p.bold,
			Italic:     p.italic,
			Underline:  p.underline,
		}
		p.canvas.SetCell(p.x, p.y, cell)
	}
	p.x++
}

// SaveANSI saves a canvas as an ANSI file
func SaveANSI(canvas *canvas.Canvas, writer io.Writer) error {
	ansiString := canvas.ToANSI()
	_, err := writer.Write([]byte(ansiString))
	return err
}