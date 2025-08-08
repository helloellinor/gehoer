package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Point struct{ X, Y float32 }
type BezierCurve struct {
	Start, Control1, Control2, End Point
}

func tokenizePath(d string) []string {
	var tokens []string
	var curr strings.Builder
	isNumChar := func(r rune) bool {
		return unicode.IsDigit(r) || r == '.' || r == '-' || r == 'e' || r == 'E'
	}
	for i, r := range d {
		if unicode.IsLetter(r) {
			if curr.Len() > 0 {
				tokens = append(tokens, curr.String())
				curr.Reset()
			}
			tokens = append(tokens, string(r))
		} else if isNumChar(r) {
			if r == '-' && curr.Len() > 0 && curr.String()[curr.Len()-1] != 'e' && curr.String()[curr.Len()-1] != 'E' {
				tokens = append(tokens, curr.String())
				curr.Reset()
			}
			curr.WriteRune(r)
		} else if r == ',' || unicode.IsSpace(r) {
			if curr.Len() > 0 {
				tokens = append(tokens, curr.String())
				curr.Reset()
			}
		}
		if i == len(d)-1 && curr.Len() > 0 {
			tokens = append(tokens, curr.String())
		}
	}
	return tokens
}

func parseSvgPathToBeziers(d string, scale float32, offset Point) []BezierCurve {
	tokens := tokenizePath(d)
	var i int
	var current, start, lastCP Point
	var lastCmd string
	var curves []BezierCurve

	for i < len(tokens) {
		cmd := tokens[i]
		i++
		isRel := (cmd == strings.ToLower(cmd))
		switch cmd {
		case "M", "m":
			x, _ := strconv.ParseFloat(tokens[i], 64)
			y, _ := strconv.ParseFloat(tokens[i+1], 64)
			i += 2
			if isRel {
				current.X += float32(x) * scale
				current.Y += float32(y) * scale
			} else {
				current = Point{float32(x)*scale + offset.X, float32(y)*scale + offset.Y}
			}
			start = current

		case "L", "l":
			x, _ := strconv.ParseFloat(tokens[i], 64)
			y, _ := strconv.ParseFloat(tokens[i+1], 64)
			i += 2
			var next Point
			if isRel {
				next = Point{current.X + float32(x)*scale, current.Y + float32(y)*scale}
			} else {
				next = Point{float32(x)*scale + offset.X, float32(y)*scale + offset.Y}
			}
			curves = append(curves, BezierCurve{Start: current, Control1: current, Control2: next, End: next})
			current = next

		case "C", "c":
			x1, _ := strconv.ParseFloat(tokens[i], 64)
			y1, _ := strconv.ParseFloat(tokens[i+1], 64)
			x2, _ := strconv.ParseFloat(tokens[i+2], 64)
			y2, _ := strconv.ParseFloat(tokens[i+3], 64)
			x, _ := strconv.ParseFloat(tokens[i+4], 64)
			y, _ := strconv.ParseFloat(tokens[i+5], 64)
			i += 6
			var p1, p2, p3 Point
			if isRel {
				p1 = Point{current.X + float32(x1)*scale, current.Y + float32(y1)*scale}
				p2 = Point{current.X + float32(x2)*scale, current.Y + float32(y2)*scale}
				p3 = Point{current.X + float32(x)*scale, current.Y + float32(y)*scale}
			} else {
				p1 = Point{float32(x1)*scale + offset.X, float32(y1)*scale + offset.Y}
				p2 = Point{float32(x2)*scale + offset.X, float32(y2)*scale + offset.Y}
				p3 = Point{float32(x)*scale + offset.X, float32(y)*scale + offset.Y}
			}
			curves = append(curves, BezierCurve{Start: current, Control1: p1, Control2: p2, End: p3})
			current = p3
			lastCP = p2

		case "S", "s":
			x2, _ := strconv.ParseFloat(tokens[i], 64)
			y2, _ := strconv.ParseFloat(tokens[i+1], 64)
			x, _ := strconv.ParseFloat(tokens[i+2], 64)
			y, _ := strconv.ParseFloat(tokens[i+3], 64)
			i += 4
			var reflect Point
			if lastCmd == "C" || lastCmd == "c" || lastCmd == "S" || lastCmd == "s" {
				reflect = Point{2*current.X - lastCP.X, 2*current.Y - lastCP.Y}
			} else {
				reflect = current
			}
			var p2, p3 Point
			if isRel {
				p2 = Point{current.X + float32(x2)*scale, current.Y + float32(y2)*scale}
				p3 = Point{current.X + float32(x)*scale, current.Y + float32(y)*scale}
			} else {
				p2 = Point{float32(x2)*scale + offset.X, float32(y2)*scale + offset.Y}
				p3 = Point{float32(x)*scale + offset.X, float32(y)*scale + offset.Y}
			}
			curves = append(curves, BezierCurve{Start: current, Control1: reflect, Control2: p2, End: p3})
			current = p3
			lastCP = p2

		case "Z", "z":
			curves = append(curves, BezierCurve{Start: current, Control1: current, Control2: start, End: start})
			current = start

		default:
			// ignoring unknown commands
		}
		lastCmd = cmd
	}
	return curves
}

func cubicBezier(p0, p1, p2, p3 Point, t float32) Point {
	u := 1 - t
	tt := t * t
	uu := u * u
	uuu := uu * u
	ttt := tt * t
	x := uuu*p0.X + 3*uu*t*p1.X + 3*u*tt*p2.X + ttt*p3.X
	y := uuu*p0.Y + 3*uu*t*p1.Y + 3*u*tt*p2.Y + ttt*p3.Y
	return Point{x, y}
}

func generateDrawCodeForCurve(curve BezierCurve, steps int) string {
	var sb strings.Builder
	last := curve.Start
	for i := 1; i <= steps; i++ {
		t := float32(i) / float32(steps)
		pt := cubicBezier(curve.Start, curve.Control1, curve.Control2, curve.End, t)
		sb.WriteString(fmt.Sprintf("rl.DrawLineV(rl.Vector2{%.2f, %.2f}, rl.Vector2{%.2f, %.2f}, rl.Black)\n", last.X, last.Y, pt.X, pt.Y))
		last = pt
	}
	return sb.String()
}

func GenerateRaylibGoSource(curves []BezierCurve, steps int) string {
	var sb strings.Builder
	sb.WriteString(`package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
	rl.InitWindow(800, 800, "Generated SVG Path")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
`)
	for i, curve := range curves {
		sb.WriteString(fmt.Sprintf("\t\t// Curve %d\n", i))
		drawCode := generateDrawCodeForCurve(curve, steps)
		for _, line := range strings.Split(drawCode, "\n") {
			if line != "" {
				sb.WriteString("\t\t" + line + "\n")
			}
		}
	}
	sb.WriteString(`
		rl.EndDrawing()
	}
}
`)
	return sb.String()
}

func main() {
	raw, err := os.ReadFile("path.txt")
	if err != nil {
		fmt.Println("Error reading path.txt:", err)
		return
	}
	d := strings.TrimSpace(string(raw))

	curves := parseSvgPathToBeziers(d, 5.0, Point{100, 100})

	source := GenerateRaylibGoSource(curves, 20)

	err = os.WriteFile("generated_svg_draw.go", []byte(source), 0644)
	if err != nil {
		fmt.Println("Error writing generated_svg_draw.go:", err)
		return
	}

	fmt.Println("Generated Go source saved to generated_svg_draw.go")
}
