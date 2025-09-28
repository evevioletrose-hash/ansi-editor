// Package objects provides a system for creating and managing moveable ANSI objects
package objects

import (
	"encoding/json"
	"fmt"

	"github.com/evevioletrose-hash/ansi-editor/pkg/canvas"
)

// Object represents a moveable ANSI art object
type Object struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	X        int            `json:"x"`
	Y        int            `json:"y"`
	Width    int            `json:"width"`
	Height   int            `json:"height"`
	Canvas   *canvas.Canvas `json:"canvas"`
	Visible  bool           `json:"visible"`
	ZIndex   int            `json:"zIndex"`
	Animated bool           `json:"animated"`
	Frames   []*canvas.Canvas `json:"frames,omitempty"`
	CurrentFrame int        `json:"currentFrame"`
}

// Scene manages a collection of objects and a background canvas
type Scene struct {
	Name       string             `json:"name"`
	Background *canvas.Canvas     `json:"background"`
	Objects    map[string]*Object `json:"objects"`
	Width      int                `json:"width"`
	Height     int                `json:"height"`
}

// NewObject creates a new ANSI object
func NewObject(id, name string, x, y, width, height int) *Object {
	return &Object{
		ID:       id,
		Name:     name,
		X:        x,
		Y:        y,
		Width:    width,
		Height:   height,
		Canvas:   canvas.NewCanvas(width, height),
		Visible:  true,
		ZIndex:   0,
		Animated: false,
		CurrentFrame: 0,
	}
}

// NewScene creates a new scene with the specified dimensions
func NewScene(name string, width, height int) *Scene {
	return &Scene{
		Name:       name,
		Background: canvas.NewCanvas(width, height),
		Objects:    make(map[string]*Object),
		Width:      width,
		Height:     height,
	}
}

// Move repositions the object to new coordinates
func (o *Object) Move(x, y int) {
	o.X = x
	o.Y = y
}

// Resize changes the object's dimensions and creates a new canvas
func (o *Object) Resize(width, height int) {
	o.Width = width
	o.Height = height
	o.Canvas = canvas.NewCanvas(width, height)
}

// AddFrame adds an animation frame to the object
func (o *Object) AddFrame(frame *canvas.Canvas) error {
	if frame.Width != o.Width || frame.Height != o.Height {
		return fmt.Errorf("frame dimensions (%dx%d) don't match object dimensions (%dx%d)",
			frame.Width, frame.Height, o.Width, o.Height)
	}
	o.Frames = append(o.Frames, frame)
	o.Animated = len(o.Frames) > 1
	return nil
}

// NextFrame advances to the next animation frame
func (o *Object) NextFrame() {
	if len(o.Frames) > 0 {
		o.CurrentFrame = (o.CurrentFrame + 1) % len(o.Frames)
	}
}

// GetCurrentCanvas returns the current canvas (either static or current animation frame)
func (o *Object) GetCurrentCanvas() *canvas.Canvas {
	if o.Animated && len(o.Frames) > 0 {
		return o.Frames[o.CurrentFrame]
	}
	return o.Canvas
}

// AddObject adds an object to the scene
func (s *Scene) AddObject(object *Object) {
	s.Objects[object.ID] = object
}

// RemoveObject removes an object from the scene
func (s *Scene) RemoveObject(id string) {
	delete(s.Objects, id)
}

// GetObject retrieves an object by ID
func (s *Scene) GetObject(id string) (*Object, bool) {
	obj, exists := s.Objects[id]
	return obj, exists
}

// Render renders the complete scene to a single canvas
func (s *Scene) Render() *canvas.Canvas {
	result := canvas.NewCanvas(s.Width, s.Height)
	
	// Copy background
	for i, cell := range s.Background.Cells {
		result.Cells[i] = cell
	}
	
	// Sort objects by Z-index (simple bubble sort for small collections)
	objects := make([]*Object, 0, len(s.Objects))
	for _, obj := range s.Objects {
		if obj.Visible {
			objects = append(objects, obj)
		}
	}
	
	// Simple Z-index sorting
	for i := 0; i < len(objects); i++ {
		for j := i + 1; j < len(objects); j++ {
			if objects[i].ZIndex > objects[j].ZIndex {
				objects[i], objects[j] = objects[j], objects[i]
			}
		}
	}
	
	// Render objects in Z-index order
	for _, obj := range objects {
		s.renderObject(result, obj)
	}
	
	return result
}

// renderObject renders a single object onto the target canvas
func (s *Scene) renderObject(target *canvas.Canvas, obj *Object) {
	objCanvas := obj.GetCurrentCanvas()
	
	for y := 0; y < obj.Height; y++ {
		for x := 0; x < obj.Width; x++ {
			srcX, srcY := x, y
			destX, destY := obj.X+x, obj.Y+y
			
			// Check bounds
			if destX < 0 || destX >= target.Width || destY < 0 || destY >= target.Height {
				continue
			}
			if srcX < 0 || srcX >= objCanvas.Width || srcY < 0 || srcY >= objCanvas.Height {
				continue
			}
			
			// Get source cell
			srcCell, err := objCanvas.GetCell(srcX, srcY)
			if err != nil {
				continue
			}
			
			// Skip transparent cells (spaces with black background)
			if srcCell.Char == ' ' && srcCell.Background == canvas.Black {
				continue
			}
			
			// Copy cell to target
			target.SetCell(destX, destY, *srcCell)
		}
	}
}

// UpdateAnimations advances all animated objects to their next frame
func (s *Scene) UpdateAnimations() {
	for _, obj := range s.Objects {
		if obj.Animated {
			obj.NextFrame()
		}
	}
}

// ToJSON serializes the scene to JSON
func (s *Scene) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// FromJSON deserializes a scene from JSON
func FromJSON(data []byte) (*Scene, error) {
	scene := &Scene{
		Objects: make(map[string]*Object),
	}
	err := json.Unmarshal(data, scene)
	return scene, err
}

// Clone creates a deep copy of an object
func (o *Object) Clone(newID string) *Object {
	clone := &Object{
		ID:           newID,
		Name:         o.Name + "_copy",
		X:            o.X,
		Y:            o.Y,
		Width:        o.Width,
		Height:       o.Height,
		Canvas:       canvas.NewCanvas(o.Width, o.Height),
		Visible:      o.Visible,
		ZIndex:       o.ZIndex,
		Animated:     o.Animated,
		CurrentFrame: o.CurrentFrame,
	}
	
	// Copy canvas data
	copy(clone.Canvas.Cells, o.Canvas.Cells)
	
	// Copy animation frames
	if o.Animated {
		clone.Frames = make([]*canvas.Canvas, len(o.Frames))
		for i, frame := range o.Frames {
			clone.Frames[i] = canvas.NewCanvas(frame.Width, frame.Height)
			copy(clone.Frames[i].Cells, frame.Cells)
		}
	}
	
	return clone
}