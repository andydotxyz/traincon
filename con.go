package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/andydotxyz/traincon/rocrail"
)

const defaultID = 3

type con struct {
	loco   *rocrail.Loco
	locoID int

	idDisplay *canvas.Text
	speed     *throttle
}

func newCon() *con {
	c := &con{locoID: defaultID + len(cons)}
	c.speed = newThrottle()
	c.speed.OnChanged = func(f float64) {
		reconnectOnErr(c.loco.SetVelocity(int(f)))
	}

	c.idDisplay = canvas.NewText("loading", theme.ErrorColor())
	c.idDisplay.TextStyle.Monospace = true
	c.idDisplay.TextSize = 32
	c.idDisplay.Alignment = fyne.TextAlignCenter

	if conn != nil {
		c.updateLoco(c.locoID)
	}
	return c
}

func (c *con) makeUI() fyne.CanvasObject {
	return container.NewBorder(nil, nil, nil, c.speed,
		container.NewGridWithColumns(1,
			c.idDisplay,
			container.NewGridWithColumns(2,
				widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {
					if c.locoID > 1 {
						c.updateLoco(c.locoID - 1)
					}
				}),
				widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
					c.updateLoco(c.locoID + 1)
				}),
			),
			container.NewGridWithColumns(2,
				widget.NewButton("Rev", func() {
					reconnectOnErr(c.loco.SetDirection(rocrail.Reverse))
				}),
				widget.NewButton("Fwd", func() {
					reconnectOnErr(c.loco.SetDirection(rocrail.Forward))
				}),
			),
			widget.NewButton("Break", func() {
				c.speed.SetValue(0)
			}),
		))
}

func (c *con) updateLoco(id int) {
	c.locoID = id
	str := fmt.Sprintf("%04d", id)
	c.idDisplay.Text = str
	c.idDisplay.Refresh()
	c.loco = conn.Loco(str)
	c.speed.SetValue(float64(c.loco.Velocity()))
}

func addControl() {
	// TODO what about tablet layout (startup with 4...)?
	c := newCon()
	cons = append(cons, c)
	grid.Add(c.makeUI())
	switch len(grid.Objects) { // TODO better intelligence based on size
	case 2, 3, 4:
		grid.Layout = layout.NewGridLayoutWithColumns(2)
	case 5, 6, 7, 8, 9:
		grid.Layout = layout.NewGridLayoutWithColumns(3)
	default:
		grid.Layout = layout.NewGridLayoutWithRows(4)
	}
	grid.Refresh()
}
