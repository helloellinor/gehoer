package game

import (
	"gehoer/camera"
	"gehoer/engraver"
	"gehoer/grid"
	"gehoer/music"
	"gehoer/musicfont"
	"gehoer/smufl"
	"gehoer/units"
	"log"

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

	font, err := musicfont.LoadMusicFont("assets/fonts/Leland/Leland.otf", int32(units.FontRenderSizePx), nil)
	if err != nil {
		panic("Failed to load music font: " + err.Error())
	}
	meta, err := smufl.LoadMetadata("/Users/eg/gehoer/external/smufl")
	if err != nil {
		log.Fatal(err)
	}

	g.engraver = engraver.NewEngraver(score, font, meta)
}

func (g *Game) Run() {
	for !rl.WindowShouldClose() {
		g.Update()
		g.Draw()
	}
}

func (g *Game) Update() {
	g.camera.Update()
	g.engraver.Update()
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.BeginMode2D(g.camera.Camera)
	g.grid.Draw(g.camera.Camera)
	g.engraver.Draw(-1)
	rl.EndMode2D()
	rl.EndDrawing()
}
