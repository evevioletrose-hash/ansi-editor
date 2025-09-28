# Command Reference

This document provides a complete reference for the ANSI Editor command-line interface.

## Starting the Editor

```bash
go run cmd/editor/main.go
# or
go build -o ansi-editor cmd/editor/main.go
./ansi-editor
```

## Canvas Commands

### new `<width>` `<height>`
Create a new canvas with specified dimensions.

**Example:**
```
new 80 25    # Creates 80x25 canvas
new 60 20    # Creates 60x20 canvas
```

### show
Display the current canvas content.

**Example:**
```
show
```

### clear `[char]` `[fg]` `[bg]`
Clear the entire canvas. All parameters are optional.

**Default:** space character, white foreground, black background

**Examples:**
```
clear                    # Clear with defaults
clear . green black      # Fill with green dots
clear # blue white       # Fill with blue # on white
```

## Drawing Commands

### draw `<x>` `<y>` `<char>` `[fg]` `[bg]`
Draw a single character at the specified position.

**Examples:**
```
draw 10 5 @ red black        # Draw red @ at (10,5)
draw 0 0 # bright_yellow     # Draw bright yellow # at top-left
```

### text `<x>` `<y>` `<text>` `[fg]` `[bg]`
Draw a text string starting at the specified position.

**Examples:**
```
text 5 10 "Hello World" green black
text 0 0 Title bright_white blue
```

### rect `<x>` `<y>` `<width>` `<height>` `<char>` `[fg]` `[bg]`
Draw a rectangle outline.

**Examples:**
```
rect 5 5 20 10 # blue black     # Blue rectangle
rect 0 0 80 25 * yellow black   # Screen border
```

### line `<x1>` `<y1>` `<x2>` `<y2>` `<char>` `[fg]` `[bg]`
Draw a line between two points.

**Examples:**
```
line 0 0 79 24 - white black    # Diagonal line
line 10 5 30 5 = red black      # Horizontal line
```

## File Operations

### load `<filename>`
Load an ANSI art file into the current canvas.

**Supported formats:** .ans, .txt with ANSI escape sequences

**Examples:**
```
load artwork.ans
load ../assets/logo.txt
```

### save `<filename>`
Save the current canvas as an ANSI file.

**Examples:**
```
save myart.ans
save output/scene1.txt
```

## Object Commands

### object add `<id>` `<name>` `<x>` `<y>` `<width>` `<height>`
Create a new object and add it to the scene.

**Examples:**
```
object add player Player 10 5 3 3
object add ui StatusBar 0 0 80 1
```

### object move `<id>` `<x>` `<y>`
Move an existing object to new coordinates.

**Examples:**
```
object move player 15 8
object move ui 0 24
```

### object list
List all objects in the current scene.

**Example:**
```
object list
```

### object show `<id>`
Display the canvas content of a specific object.

**Examples:**
```
object show player
object show ui
```

## Export Commands

### export wasm `<output_dir>` `<name>` `<version>`
Export the current scene as a WebAssembly bundle.

**Examples:**
```
export wasm ./my-game "My Game" 1.0.0
export wasm ../web-app "Demo App" 0.1.0
```

The exported bundle will contain:
- `index.html` - Web page
- `ansi-runtime.js` - JavaScript runtime
- `main.go` - Go WASM source
- `build.sh` / `build.bat` - Build scripts
- `bundle.json` - Scene data

## Color Names

### Standard Colors
- `black`, `red`, `green`, `yellow`
- `blue`, `magenta`, `cyan`, `white`

### Bright Colors
- `bright_black`, `bright_red`, `bright_green`, `bright_yellow`
- `bright_blue`, `bright_magenta`, `bright_cyan`, `bright_white`

**Alternative names:**
- `brightblack`, `brightred`, etc. (no underscore)

## General Commands

### help, h
Show the help menu with available commands.

### quit, q, exit
Exit the editor.

## Command Examples

### Creating a Simple Scene
```
new 40 15
text 5 2 "Welcome!" bright_green black
rect 3 1 20 4 # yellow black
object add sprite Character 25 8 3 3
show
```

### Drawing a House
```
new 30 20
rect 10 10 10 8 # brown black
line 10 10 15 6 # red black
line 15 6 20 10 # red black
draw 12 15 D yellow black
text 5 18 "My House" white black
show
```

### Creating a Game Scene
```
new 60 20
clear . green black
text 2 1 "Level 1" white black
object add player Player 30 10 1 1
object add enemy Enemy1 10 15 1 1
object add enemy Enemy2 45 8 1 1
object list
export wasm ./my-game "Simple Game" 1.0
```

## Tips

1. **Coordinates** are 0-based (top-left is 0,0)
2. **Colors** can be abbreviated (e.g., 'r' for red in some contexts)
3. **Strings with spaces** should be quoted: `"Hello World"`
4. **File paths** can be relative or absolute
5. **Object IDs** must be unique within a scene
6. **Canvas size** affects export performance - keep reasonable

## Error Messages

- **"coordinates out of bounds"** - Position is outside canvas
- **"object not found"** - Object ID doesn't exist
- **"invalid color"** - Color name not recognized
- **"file not found"** - File path is incorrect
- **"invalid dimensions"** - Width or height is invalid

## Keyboard Shortcuts

While in the editor:
- **Tab** - Auto-complete command names
- **Up/Down arrows** - Command history
- **Ctrl+C** - Exit editor
- **Ctrl+L** - Clear screen (terminal-dependent)