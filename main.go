//go:generate fyne bundle -o bundled.go Icon.png

package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/andydotxyz/traincon/rocrail"
)

var (
	conn   *rocrail.Connection
	loco   *rocrail.Loco
	locoID = 3

	idDisplay *canvas.Text
	throttle  *widget.Slider
	win       fyne.Window
)

func connect() {
	d := dialog.NewProgressInfinite("Connecting",
		"Connecting to rocrail...", win)
	d.Show()
	go func() {
		host, port := userPref(fyne.CurrentApp())
		if host == "" {
			d.Hide()
			updatePref(fyne.CurrentApp())
			return
		}
		c, err := rocrail.Connect(host, port)
		if err != nil {
			d.Hide()
			dialog.ShowError(err, win)
			return
		}

		conn = c
		updateLoco(locoID)
		d.Hide()
	}()
}

func reconnectOnErr(err error) {
	if err == nil {
		return
	}

	connect()
}

func updateLoco(id int) {
	locoID = id
	str := fmt.Sprintf("%04d", locoID)
	idDisplay.Text = str
	idDisplay.Refresh()
	loco = conn.Loco(str)
	throttle.SetValue(float64(loco.Velocity()))
}

func updatePref(a fyne.App) {
	host := widget.NewEntry()
	host.SetPlaceHolder("localhost")
	port := widget.NewEntry()
	port.SetPlaceHolder("8051")

	dialog.ShowForm("Connect", "Go", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Hostname", host),
			widget.NewFormItem("Port", port),
		}, func(ok bool) {
			if !ok {
				a.Preferences().SetString("server.host", "")
				return
			}

			a.Preferences().SetString("server.host", host.Text)
			p, _ := strconv.Atoi(port.Text)
			a.Preferences().SetInt("server.port", p)

			connect()
		}, win)
}

func userPref(a fyne.App) (string, int) {
	return a.Preferences().String("server.host"),
		a.Preferences().Int("server.port")
}

func main() {
	a := app.NewWithID("xyz.andy.traincon")
	a.SetIcon(resourceIconPng)
	win = a.NewWindow("Train Con")
	connect()

	throttle = widget.NewSlider(0, 100)
	throttle.OnChanged = func(f float64) {
		reconnectOnErr(loco.SetVelocity(int(f)))
	}
	throttle.Orientation = widget.Vertical

	idDisplay = canvas.NewText("0000", theme.ErrorColor())
	idDisplay.TextStyle.Monospace = true
	idDisplay.TextSize = 32
	idDisplay.Alignment = fyne.TextAlignCenter
	updateLoco(3)

	win.SetContent(container.NewBorder(nil, nil, nil, throttle,
		container.NewGridWithRows(3,
			idDisplay,
			container.NewGridWithColumns(2,
				widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {
					if locoID > 1 {
						updateLoco(locoID - 1)
					}
				}),
				widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
					updateLoco(locoID+1)
				}),
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

	win.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Disconnect", func() {
				updatePref(fyne.CurrentApp())
			}))))
	win.ShowAndRun()
}
