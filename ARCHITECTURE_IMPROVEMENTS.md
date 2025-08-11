# Architecture Improvements: Separation of Drawing Logic from Engraving/Layout

## Problem Statement Addressed

The original codebase had several architectural issues:

1. **Tight Coupling**: Drawing logic was directly embedded in layout/engraving code
2. **Redundant Methods**: Multiple `DrawGlyph` implementations with different signatures
3. **Inconsistent Go Patterns**: Mixed responsibilities and non-standard interface usage
4. **Untestable Layout**: Layout logic required graphics context, making testing difficult

## Solution Implemented

### 1. Command Pattern for Drawing Operations

Created a command-based rendering system in the `renderer/` package:

- `DrawCommand` interface: Represents any drawing operation
- `Renderer` interface: Abstracts the actual drawing backend
- Specific command types: `LineCommand`, `TextCommand`, `GlyphCommand`, `RectangleLinesCommand`
- `CommandBuffer`: Batches commands for efficient execution

### 2. Clean Separation of Concerns

**Before**: Engraver directly made drawing calls
```go
func (e *Engraver) Draw(x, y float32) {
    rl.DrawLineEx(...)  // Direct drawing
    rl.DrawTextEx(...)  // Mixed with layout
}
```

**After**: Engraver generates commands, renderer executes them
```go
func (e *Engraver) GenerateDrawCommands(x, y float32, buffer *CommandBuffer) {
    buffer.AddCommand(NewLineCommand(...))  // Pure layout logic
}
```

### 3. Interface-Based Design

- `Renderer` interface allows multiple backends (currently Raylib, easily extendable)
- `DrawCommand` interface enables polymorphic command handling
- Existing `MusicElement` interface preserved and enhanced

### 4. Eliminated Redundancies

**Before**: Multiple drawing methods with inconsistent signatures
- `DrawGlyph(font, glyph, x, y, offset, color)` in `glyph.go`
- `DrawGlyph(glyphName, x, y, color)` in `engraver.go`

**After**: Single command creation pattern
- `CreateGlyphCommand(...)` generates commands consistently
- `CalculateGlyphPosition(...)` pure layout calculation
- All drawing unified through command pattern

## Benefits Achieved

### 1. Testability
Layout logic can now be tested without graphics context:
```go
buffer := renderer.NewCommandBuffer()
engraver.GenerateDrawCommands(x, y, buffer)
// Verify commands generated correctly without rendering
```

### 2. Flexibility
Multiple rendering backends possible:
- Current: `RaylibRenderer`
- Possible: `SVGRenderer`, `PDFRenderer`, `TestRenderer`

### 3. Performance
Command batching allows optimization:
- Batch similar operations
- Sort by render state
- Minimize backend calls

### 4. Maintainability
Clear separation of responsibilities:
- **Engraver**: Layout and positioning calculations only
- **Renderer**: Drawing operations only
- **Game**: Orchestration and command buffer management

## Files Modified

### New Files
- `renderer/renderer.go`: Core interfaces and command types
- `renderer/raylib.go`: Raylib-specific renderer implementation
- `renderer/utils.go`: Conversion utilities between types

### Modified Files
- `engraver/engraver.go`: Removed direct drawing, added command generation
- `engraver/glyph.go`: Converted to pure layout calculation utilities
- `engraver/note.go`: Command-based note rendering
- `engraver/staff.go`: Command-based staff rendering
- `engraver/clef.go`: Simplified using command pattern
- `grid/grid.go`: Dual mode (command generation + backwards compatibility)
- `game/game.go`: Uses command buffer and renderer

## Go Best Practices Implemented

1. **Interface Segregation**: Small, focused interfaces (`Renderer`, `DrawCommand`)
2. **Dependency Inversion**: High-level modules don't depend on low-level details
3. **Single Responsibility**: Each component has one clear purpose
4. **Open/Closed Principle**: Easy to add new renderers without modifying existing code

## Example Usage

```go
// Pure layout calculation (testable)
buffer := renderer.NewCommandBuffer()
engraver.GenerateDrawCommands(x, y, buffer)

// Rendering execution (when graphics context available)
renderer := renderer.NewRaylibRenderer()
buffer.Execute(renderer)
```

This architecture enables music notation applications to:
- Generate layouts server-side without graphics
- Support multiple output formats
- Test layout logic thoroughly
- Maintain clean, separated concerns