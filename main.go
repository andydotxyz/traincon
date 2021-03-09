//go:generate fyne bundle -o bundled.go Icon.png

package main

import (
	"fmt"
	"log"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func connect() net.Conn {
	conn, err := net.Dial("tcp", "localhost:8051")
	if err != nil {
		log.Println("Failed to connect to RocRail")
		return nil
	}

	return conn
}

func send(conn net.Conn, cmd, args string) {
	str := "<xmlh><xml size=\"%d\" name=\"%s\"/></xmlh>%s"
	_, err := fmt.Fprintf(conn, str, len(args), cmd, args)

	if err != nil {
		log.Println("Failed to send, retrying", err)
		connect()
		_, _ = fmt.Fprintf(conn, str, len(args), cmd, args)
	}
}

func main() {
	a := app.New()
	a.SetIcon(resourceIconPng)
	w := a.NewWindow("Train Con")
	conn := connect()
	fwd := true

	throttle := widget.NewSlider(0, 100)
	throttle.OnChanged = func(f float64) {
		cmd := fmt.Sprintf("<lc id=\"0003\" V=\"%d\" dir=\"%t\" cmd=\"velocity\" />",
			int(f), fwd)
		send(conn, "lc", cmd)
	}
	throttle.Orientation = widget.Vertical

	loco := canvas.NewText("0003", theme.ErrorColor())
	loco.TextStyle.Monospace = true
	loco.TextSize = 32
	loco.Alignment = fyne.TextAlignCenter

	w.SetContent(container.NewBorder(nil, nil, nil, throttle,
		container.NewGridWithRows(3,
			loco,
			container.NewGridWithColumns(2,
				widget.NewButton("Rev", func() {
					fwd = false
					send(conn, "lc", "<lc id=\"0003\" dir=\"false\" cmd=\"velocity\" />")
				}),
				widget.NewButton("Fwd", func() {
					fwd = true
					send(conn, "lc", "<lc id=\"0003\" dir=\"true\" cmd=\"velocity\" />")
				}),
			),
			widget.NewButton("Break", func() {
				throttle.SetValue(0)
			}),
		)))
	w.ShowAndRun()
}
