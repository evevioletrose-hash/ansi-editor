package editor

import (
	"fmt"
	
	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
	"github.com/evevioletrose-hash/ansi-editor/pkg/ansi"
)

// Editor represents the main ANSI editor with multiple pages/scenes
type Editor struct {
	Pages       []*Page
	CurrentPage int
	Modified    bool
}

// Page represents a single canvas/scene in the editor
type Page struct {
	Name     string
	Canvas   *canvas.Canvas
	Modified bool
}

// NewEditor creates a new editor instance
func NewEditor() *Editor {
	return &Editor{
		Pages:       make([]*Page, 0),
		CurrentPage: -1,
		Modified:    false,
	}
}

// NewPage creates a new page with the specified dimensions
func (e *Editor) NewPage(name string, width, height int) *Page {
	page := &Page{
		Name:     name,
		Canvas:   canvas.NewCanvas(width, height),
		Modified: false,
	}
	
	e.Pages = append(e.Pages, page)
	e.CurrentPage = len(e.Pages) - 1
	e.Modified = true
	
	return page
}

// GetCurrentPage returns the currently active page
func (e *Editor) GetCurrentPage() *Page {
	if e.CurrentPage < 0 || e.CurrentPage >= len(e.Pages) {
		return nil
	}
	return e.Pages[e.CurrentPage]
}

// SetCurrentPage changes the active page by index
func (e *Editor) SetCurrentPage(index int) error {
	if index < 0 || index >= len(e.Pages) {
		return fmt.Errorf("page index %d is out of range", index)
	}
	e.CurrentPage = index
	return nil
}

// SetCurrentPageByName changes the active page by name
func (e *Editor) SetCurrentPageByName(name string) error {
	for i, page := range e.Pages {
		if page.Name == name {
			e.CurrentPage = i
			return nil
		}
	}
	return fmt.Errorf("page with name '%s' not found", name)
}

// LoadPage loads an ANSI file as a new page
func (e *Editor) LoadPage(filename, pageName string) error {
	c, err := ansi.LoadFromFile(filename)
	if err != nil {
		return fmt.Errorf("failed to load file %s: %w", filename, err)
	}
	
	if pageName == "" {
		pageName = fmt.Sprintf("Page %d", len(e.Pages)+1)
	}
	
	page := &Page{
		Name:     pageName,
		Canvas:   c,
		Modified: false,
	}
	
	e.Pages = append(e.Pages, page)
	e.CurrentPage = len(e.Pages) - 1
	
	return nil
}

// SaveCurrentPage saves the current page to a file
func (e *Editor) SaveCurrentPage(filename string) error {
	page := e.GetCurrentPage()
	if page == nil {
		return fmt.Errorf("no active page to save")
	}
	
	err := ansi.SaveToFile(page.Canvas, filename)
	if err != nil {
		return fmt.Errorf("failed to save page: %w", err)
	}
	
	page.Modified = false
	e.updateModifiedStatus()
	
	return nil
}

// DeletePage removes a page by index
func (e *Editor) DeletePage(index int) error {
	if index < 0 || index >= len(e.Pages) {
		return fmt.Errorf("page index %d is out of range", index)
	}
	
	// Remove the page
	e.Pages = append(e.Pages[:index], e.Pages[index+1:]...)
	
	// Adjust current page index
	if e.CurrentPage >= len(e.Pages) {
		e.CurrentPage = len(e.Pages) - 1
	}
	if e.CurrentPage < 0 && len(e.Pages) > 0 {
		e.CurrentPage = 0
	}
	
	e.Modified = true
	return nil
}

// DuplicatePage creates a copy of the specified page
func (e *Editor) DuplicatePage(index int, newName string) error {
	if index < 0 || index >= len(e.Pages) {
		return fmt.Errorf("page index %d is out of range", index)
	}
	
	originalPage := e.Pages[index]
	
	// Create new canvas and copy content
	newCanvas := canvas.NewCanvas(originalPage.Canvas.Width, originalPage.Canvas.Height)
	for y := 0; y < originalPage.Canvas.Height; y++ {
		for x := 0; x < originalPage.Canvas.Width; x++ {
			cell, _ := originalPage.Canvas.GetCell(x, y)
			newCanvas.SetCell(x, y, cell.Char, cell.FgColor, cell.BgColor, cell.Attributes)
		}
	}
	
	newPage := &Page{
		Name:     newName,
		Canvas:   newCanvas,
		Modified: false,
	}
	
	e.Pages = append(e.Pages, newPage)
	e.Modified = true
	
	return nil
}

// GetPageInfo returns information about all pages
func (e *Editor) GetPageInfo() []PageInfo {
	info := make([]PageInfo, len(e.Pages))
	for i, page := range e.Pages {
		info[i] = PageInfo{
			Index:    i,
			Name:     page.Name,
			Width:    page.Canvas.Width,
			Height:   page.Canvas.Height,
			Modified: page.Modified,
			Active:   i == e.CurrentPage,
		}
	}
	return info
}

// PageInfo contains metadata about a page
type PageInfo struct {
	Index    int
	Name     string
	Width    int
	Height   int
	Modified bool
	Active   bool
}

// MarkCurrentPageModified marks the current page as modified
func (e *Editor) MarkCurrentPageModified() {
	page := e.GetCurrentPage()
	if page != nil {
		page.Modified = true
		e.Modified = true
	}
}

// updateModifiedStatus updates the editor's overall modified status
func (e *Editor) updateModifiedStatus() {
	e.Modified = false
	for _, page := range e.Pages {
		if page.Modified {
			e.Modified = true
			break
		}
	}
}

// HasModifiedPages returns true if any pages have unsaved changes
func (e *Editor) HasModifiedPages() bool {
	return e.Modified
}

// GetModifiedPageNames returns names of all modified pages
func (e *Editor) GetModifiedPageNames() []string {
	var names []string
	for _, page := range e.Pages {
		if page.Modified {
			names = append(names, page.Name)
		}
	}
	return names
}