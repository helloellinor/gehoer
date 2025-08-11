package game

import (
	"gehoer/camera"
	"gehoer/engraver"
	"gehoer/grid"
	"gehoer/music"
	"gehoer/musicfont"
	"gehoer/settings"
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	grid     *grid.Grid
	camera   *camera.Controller
	engraver *engraver.Engraver
}

func New() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.grid = grid.New(units.GridSpacingPx, 4000, 4000, units.GridFontSizePx)
	g.camera = camera.NewController(1200, 800)

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
	//rl.DrawTextEx(g.engraver.MusicFont.Font, "\uE0A2", rl.NewVector2(50, 300), 300, 0, rl.Black)
	g.grid.Draw(g.camera.Camera)
	g.engraver.Draw(0, 0)
	rl.EndMode2D()
	rl.EndDrawing()
}
