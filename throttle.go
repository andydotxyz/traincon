package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	throttleLineWidth = float32(5)
	throttleHeight    = float32(25)
	throttleWidth     = float32(75)
)

var _ fyne.Draggable = (*throttle)(nil)
var _ fyne.Widget = (*throttle)(nil)

type throttle struct {
	widget.BaseWidget

	OnChanged func(f float64)
	val       float64
}

func newThrottle() *throttle {
	t := &throttle{}
	t.ExtendBaseWidget(t)
	return t
}

func (t *throttle) CreateRenderer() fyne.WidgetRenderer {
	return &throttleRenderer{t: t,
		handle: canvas.NewRectangle(theme.PrimaryColor()),
		line:   canvas.NewRectangle(theme.ShadowColor()),
	}
}

func (t *throttle) Dragged(ev *fyne.DragEvent) {
	max := t.Size().Height - throttleHeight
	pos := ev.Position.Y - throttleHeight/2
	ratio := 1 - (pos / max)

	val := ratio * 100
	if val < 0 {
		val = 0
	}
	if val > 100 {
		val = 100
	}
	t.val = float64(val)
	t.Refresh()
	if t.OnChanged != nil {
		t.OnChanged(t.val)
	}
}

func (t *throttle) DragEnd() {
}

func (t *throttle) SetValue(val float64) {
	diff := val - t.val
	inc := diff / 10
	go func() {
		for i := 0; i < 10; i++ {
			t.val += inc
			if i == 9 {
				t.val = val
			}

			t.Refresh()
			if t.OnChanged != nil {
				t.OnChanged(t.val)
			}

			time.Sleep(time.Millisecond * 10)
		}
	}()
}

var _ fyne.WidgetRenderer = (*throttleRenderer)(nil)

type throttleRenderer struct {
	t *throttle

	handle, line *canvas.Rectangle
}

func (r *throttleRenderer) Destroy() {
}

func (r *throttleRenderer) Layout(s fyne.Size) {
	middle := s.Width / 2
	inset := throttleHeight / 2

	r.line.Resize(fyne.NewSize(throttleLineWidth, s.Height-throttleHeight))
	r.line.Move(fyne.NewPos(middle-throttleLineWidth/2, inset))

	max := r.t.Size().Height - throttleHeight
	offset := float32(r.t.val/100) * max

	r.handle.Resize(fyne.NewSize(s.Width, throttleHeight))
	r.handle.Move(fyne.NewPos(0, s.Height-offset-throttleHeight))
	r.handle.Refresh()
}

func (r *throttleRenderer) MinSize() fyne.Size {
	return fyne.NewSize(throttleWidth, throttleHeight*2)
}

func (r *throttleRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.line, r.handle}
}

func (r *throttleRenderer) Refresh() {
	r.handle.FillColor = theme.PrimaryColor()
	r.line.FillColor = theme.ShadowColor()

	r.Layout(r.t.Size())
}
