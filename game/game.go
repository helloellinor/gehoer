package game

import (
	"fmt"
	"gehoer/camera"
	"gehoer/engraver"
	"gehoer/grid"
	"gehoer/music"
	"gehoer/musicfont"
	"gehoer/renderer"
	"gehoer/settings"
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	grid          *grid.Grid
	camera        *camera.Controller
	engraver      *engraver.Engraver
	renderer      renderer.Renderer
	commandBuffer *renderer.CommandBuffer

	// Score interaction state
	scorePosition rl.Vector2 // Current position of the score
	isDragging    bool       // Is the score currently being dragged
	dragOffset    rl.Vector2 // Offset between mouse and score position when dragging started
	scoreSelected bool       // Is the score currently selected
}

func New() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.grid = grid.New(units.GridSpacingPx, 4000, 4000, units.GridFontSizePx)
	g.camera = camera.NewController(1200, 800)
	g.renderer = renderer.NewRaylibRenderer()
	g.commandBuffer = renderer.NewCommandBuffer()

	// Initialize score position
	g.scorePosition = rl.Vector2{X: 100, Y: 200}
	g.isDragging = false
	g.scoreSelected = false

	// Load sample score from JSON file
	score, err := music.LoadScoreFromJSON("assets/scores/simple_piano.json")
	if err != nil {
		panic("Failed to load score JSON: " + err.Error())
	}

	font, err := musicfont.LoadMusicFont("external/smufl", "assets/fonts/Leland/leland_metadata.json", "assets/fonts/Leland/Leland.otf", settings.MusicFontSizePx)
	if err != nil {
		panic("Failed to load music font: " + err.Error())
	}

	g.engraver = engraver.NewEngraver(score, font)
}

func (g *Game) Run() {
	for !rl.WindowShouldClose() {
		g.Update()
		g.Draw()
	}
}

func (g *Game) Update() {
	g.camera.Update()
	g.handleScoreInteraction()
}

// handleScoreInteraction handles mouse interaction with the score
func (g *Game) handleScoreInteraction() {
	mousePos := rl.GetMousePosition()
	worldMousePos := rl.GetScreenToWorld2D(mousePos, g.camera.Camera)

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		// Check if mouse is over the score (simple bounding box check)
		scoreRect := g.getScoreBoundingRect()
		if g.pointInRect(worldMousePos, scoreRect) {
			g.scoreSelected = true
			g.isDragging = true
			g.dragOffset = rl.Vector2Subtract(g.scorePosition, worldMousePos)
		} else {
			g.scoreSelected = false
		}
	}

	if rl.IsMouseButtonDown(rl.MouseLeftButton) && g.isDragging {
		// Drag the score
		g.scorePosition = rl.Vector2Add(worldMousePos, g.dragOffset)
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		g.isDragging = false
	}
}

// getScoreBoundingRect returns an approximate bounding rectangle for the score
func (g *Game) getScoreBoundingRect() rl.Rectangle {
	// Estimate score dimensions (this is simplified - in a real implementation
	// you would calculate the actual rendered dimensions)
	return rl.Rectangle{
		X:      g.scorePosition.X,
		Y:      g.scorePosition.Y - 50, // Account for staff height
		Width:  800,                    // Approximate width
		Height: 100,                    // Approximate height for single staff
	}
}

// pointInRect checks if a point is inside a rectangle
func (g *Game) pointInRect(point rl.Vector2, rect rl.Rectangle) bool {
	return point.X >= rect.X && point.X <= rect.X+rect.Width &&
		point.Y >= rect.Y && point.Y <= rect.Y+rect.Height
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode2D(g.camera.Camera)

	// Clear the command buffer for this frame
	g.commandBuffer.Clear()

	// Generate all drawing commands
	g.grid.GenerateDrawCommands(g.camera.Camera, g.commandBuffer)

	// Draw the score at its current position
	g.engraver.GenerateDrawCommands(g.scorePosition.X, g.scorePosition.Y, g.commandBuffer)

	// Draw selection indicator if score is selected
	if g.scoreSelected {
		g.drawSelectionIndicator()
	}

	// Execute all commands
	g.commandBuffer.Execute(g.renderer)

	rl.EndMode2D()

	// Draw UI information
	g.drawUI()

	rl.EndDrawing()
}

// drawSelectionIndicator draws a visual indicator around the selected score
func (g *Game) drawSelectionIndicator() {
	scoreRect := g.getScoreBoundingRect()

	// Draw a selection rectangle
	start1 := renderer.Vector2{X: scoreRect.X - 5, Y: scoreRect.Y - 5}
	end1 := renderer.Vector2{X: scoreRect.X + scoreRect.Width + 5, Y: scoreRect.Y - 5}
	start2 := renderer.Vector2{X: scoreRect.X + scoreRect.Width + 5, Y: scoreRect.Y - 5}
	end2 := renderer.Vector2{X: scoreRect.X + scoreRect.Width + 5, Y: scoreRect.Y + scoreRect.Height + 5}
	start3 := renderer.Vector2{X: scoreRect.X + scoreRect.Width + 5, Y: scoreRect.Y + scoreRect.Height + 5}
	end3 := renderer.Vector2{X: scoreRect.X - 5, Y: scoreRect.Y + scoreRect.Height + 5}
	start4 := renderer.Vector2{X: scoreRect.X - 5, Y: scoreRect.Y + scoreRect.Height + 5}
	end4 := renderer.Vector2{X: scoreRect.X - 5, Y: scoreRect.Y - 5}

	selectionColor := renderer.Color{R: 100, G: 150, B: 255, A: 255}
	thickness := float32(2.0)

	g.commandBuffer.AddCommand(renderer.NewLineCommand(start1, end1, thickness, selectionColor))
	g.commandBuffer.AddCommand(renderer.NewLineCommand(start2, end2, thickness, selectionColor))
	g.commandBuffer.AddCommand(renderer.NewLineCommand(start3, end3, thickness, selectionColor))
	g.commandBuffer.AddCommand(renderer.NewLineCommand(start4, end4, thickness, selectionColor))
}

// drawUI draws UI information in screen space
func (g *Game) drawUI() {
	// Draw score position and selection status
	infoText := ""
	if g.scoreSelected {
		infoText = "Score Selected - "
	}

	// Format position using fmt.Sprintf
	positionStr := fmt.Sprintf("Position: (%.1f, %.1f)", g.scorePosition.X, g.scorePosition.Y)
	infoText += positionStr

	if g.isDragging {
		infoText += " [DRAGGING]"
	}

	rl.DrawText(infoText, 10, 10, 20, rl.Black)
}
