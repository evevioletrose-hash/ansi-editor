package canvas

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	// Standard colors
	ColorBlack   = 30
	ColorRed     = 31
	ColorGreen   = 32
	ColorYellow  = 33
	ColorBlue    = 34
	ColorMagenta = 35
	ColorCyan    = 36
	ColorWhite   = 37
	
	// Bright colors
	ColorBrightBlack   = 90
	ColorBrightRed     = 91
	ColorBrightGreen   = 92
	ColorBrightYellow  = 93
	ColorBrightBlue    = 94
	ColorBrightMagenta = 95
	ColorBrightCyan    = 96
	ColorBrightWhite   = 97
)

// ANSI attributes
const (
	AttrReset     = 0
	AttrBold      = 1
	AttrDim       = 2
	AttrItalic    = 3
	AttrUnderline = 4
	AttrBlink     = 5
	AttrReverse   = 7
	AttrStrike    = 9
)

// RenderANSI converts the canvas to an ANSI string for terminal display
func (c *Canvas) RenderANSI() string {
	var builder strings.Builder
	
	// Reset terminal at the start
	builder.WriteString("\033[0m")
	
	var lastFg, lastBg, lastAttr = -1, -1, -1
	
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			cell := c.Grid[y][x]
			
			// Only output color/attribute codes if they changed
			if cell.FgColor != lastFg || cell.BgColor != lastBg || cell.Attributes != lastAttr {
				// Reset if needed
				if lastAttr != 0 && cell.Attributes != lastAttr {
					builder.WriteString("\033[0m")
					lastFg, lastBg, lastAttr = -1, -1, -1
				}
				
				// Set attributes
				if cell.Attributes > 0 && cell.Attributes != lastAttr {
					builder.WriteString(fmt.Sprintf("\033[%dm", cell.Attributes))
				}
				
				// Set foreground color
				if cell.FgColor != lastFg {
					builder.WriteString(fmt.Sprintf("\033[%dm", cell.FgColor))
				}
				
				// Set background color (add 10 to foreground color code)
				if cell.BgColor != lastBg {
					builder.WriteString(fmt.Sprintf("\033[%dm", cell.BgColor+10))
				}
				
				lastFg = cell.FgColor
				lastBg = cell.BgColor
				lastAttr = cell.Attributes
			}
			
			builder.WriteRune(cell.Char)
		}
		
		// Add newline at end of each row (except the last)
		if y < c.Height-1 {
			builder.WriteString("\n")
		}
	}
	
	// Reset terminal at the end
	builder.WriteString("\033[0m")
	
	return builder.String()
}

// RenderPlainText converts the canvas to plain text (characters only)
func (c *Canvas) RenderPlainText() string {
	var builder strings.Builder
	
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			builder.WriteRune(c.Grid[y][x].Char)
		}
		if y < c.Height-1 {
			builder.WriteString("\n")
		}
	}
	
	return builder.String()
}

// GetANSISequence returns the ANSI escape sequence for a cell's formatting
func GetANSISequence(fgColor, bgColor, attributes int) string {
	var parts []string
	
	if attributes > 0 {
		parts = append(parts, fmt.Sprintf("%d", attributes))
	}
	
	if fgColor >= 0 {
		parts = append(parts, fmt.Sprintf("%d", fgColor))
	}
	
	if bgColor >= 0 {
		parts = append(parts, fmt.Sprintf("%d", bgColor+10))
	}
	
	if len(parts) == 0 {
		return ""
	}
	
	return fmt.Sprintf("\033[%sm", strings.Join(parts, ";"))
}