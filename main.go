//go:generate fyne bundle -o bundled.go Icon.png

package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/andydotxyz/traincon/rocrail"
)

var conn *rocrail.Connection

func connect() *rocrail.Connection {
	conn, err := rocrail.Connect("localhost",8051)
	if err != nil {
		log.Println("Failed to connect to RocRail")
		return nil
	}

	return conn
}

func reconnectOnErr(err error) {
	if err == nil {
		return
	}

	log.Println("Failed to send, retrying")
	conn = connect()
}

func main() {
	a := app.New()
	a.SetIcon(resourceIconPng)
	w := a.NewWindow("Train Con")
	conn = connect()
	loco := conn.Loco("0003")

	throttle := widget.NewSlider(0, 100)
	throttle.OnChanged = func(f float64) {
		reconnectOnErr(loco.SetVelocity(int(f)))
	}
	throttle.Orientation = widget.Vertical

	id := canvas.NewText("0003", theme.ErrorColor())
	id.TextStyle.Monospace = true
	id.TextSize = 32
	id.Alignment = fyne.TextAlignCenter

	w.SetContent(container.NewBorder(nil, nil, nil, throttle,
		container.NewGridWithRows(3,
			id,
			container.NewGridWithColumns(2,
				widget.NewButton("Rev", func() {
					reconnectOnErr(loco.SetDirection(rocrail.Reverse))
				}),
				widget.NewButton("Fwd", func() {
					reconnectOnErr(loco.SetDirection(rocrail.Forward))
				}),
			),
			widget.NewButton("Break", func() {
				throttle.SetValue(0)
			}),
		)))
	w.ShowAndRun()
}
