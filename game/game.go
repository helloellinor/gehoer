package game

import (
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

	// Load sample score from JSON file
	score, err := music.LoadScoreFromJSON("assets/scores/lisa_gikk_til_skolen.json")
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
	//g.engraver.Update()
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode2D(g.camera.Camera)

	// Clear the command buffer for this frame
	g.commandBuffer.Clear()

	// Generate all drawing commands
	g.grid.GenerateDrawCommands(g.camera.Camera, g.commandBuffer)
	g.engraver.GenerateDrawCommands(0, 0, g.commandBuffer)

	// Execute all commands
	g.commandBuffer.Execute(g.renderer)

	rl.EndMode2D()
	rl.EndDrawing()
}
