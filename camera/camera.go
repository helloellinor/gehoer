package camera

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Controller struct {
	Camera    rl.Camera2D
	Dragging  bool
	LastPos   rl.Vector2
	MinZoom   float32
	MaxZoom   float32
	MoveSpeed float32
}

func NewController(screenWidth, screenHeight int) *Controller {
	return &Controller{
		Camera: rl.Camera2D{
			Target:   rl.NewVector2(0, 0),
			Offset:   rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2),
			Zoom:     1.0,
			Rotation: 0,
		},
		MinZoom:   0.1,
		MaxZoom:   10.0,
		MoveSpeed: 10,
	}
}

// Update processes input and updates the camera accordingly.
// Should be called once per frame.
func (cc *Controller) Update() {
	// Keyboard pan
	speed := cc.MoveSpeed / cc.Camera.Zoom
	if rl.IsKeyDown(rl.KeyRight) {
		cc.Camera.Target.X += speed
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		cc.Camera.Target.X -= speed
	}
	if rl.IsKeyDown(rl.KeyUp) {
		cc.Camera.Target.Y -= speed
	}
	if rl.IsKeyDown(rl.KeyDown) {
		cc.Camera.Target.Y += speed
	}

	// Mouse drag pan
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		mp := rl.GetMousePosition()
		if !cc.Dragging {
			cc.Dragging = true
			cc.LastPos = mp
		} else {
			dx := (cc.LastPos.X - mp.X) / cc.Camera.Zoom
			dy := (cc.LastPos.Y - mp.Y) / cc.Camera.Zoom
			cc.Camera.Target.X += dx
			cc.Camera.Target.Y += dy
			cc.LastPos = mp
		}
	} else {
		cc.Dragging = false
	}

	// Mouse wheel zoom (zoom towards pointer)
	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		mouse := rl.GetMousePosition()
		worldBefore := rl.GetScreenToWorld2D(mouse, cc.Camera)
		cc.Camera.Zoom += wheel * 0.1
		if cc.Camera.Zoom < cc.MinZoom {
			cc.Camera.Zoom = cc.MinZoom
		} else if cc.Camera.Zoom > cc.MaxZoom {
			cc.Camera.Zoom = cc.MaxZoom
		}
		worldAfter := rl.GetScreenToWorld2D(mouse, cc.Camera)
		cc.Camera.Target.X += worldBefore.X - worldAfter.X
		cc.Camera.Target.Y += worldBefore.Y - worldAfter.Y
	}
}
