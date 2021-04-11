//go:generate fyne bundle -o bundled.go Icon.png

package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/andydotxyz/traincon/rocrail"
)

var (
	conn *rocrail.Connection
	cons []*con

	grid *fyne.Container
	win  fyne.Window
)

func connect() {
	d := dialog.NewProgressInfinite("Connecting",
		"Connecting to rocrail...", win)
	d.Show()
	go func() {
		host, port := userPref(fyne.CurrentApp())
		c, err := rocrail.Connect(host, port)
		if err != nil {
			d.Hide()
			d := dialog.NewError(err, win)
			d.SetOnClosed(func() {
				showLogin(fyne.CurrentApp())
			})
			d.Show()
			return
		}

		conn = c
		for _, ctrl := range cons {
			ctrl.updateLoco(ctrl.locoID)
		}
		d.Hide()
	}()
}

func reconnectOnErr(err error) {
	if err == nil {
		return
	}

	connect()
}

func showLogin(a fyne.App) {
	myHost, myPort := userPref(a)
	host := widget.NewEntry()
	host.SetPlaceHolder("localhost")
	if myHost != "" {
		host.SetText(myHost)
	}
	port := widget.NewEntry()
	port.SetPlaceHolder("8051")
	if myPort != 0 {
		port.SetText(fmt.Sprintf("%d", myPort))
	}

	dialog.ShowForm("Connect to rocrail", "Go", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Hostname", host),
			widget.NewFormItem("Port", port),
		}, func(ok bool) {
			if !ok {
				d := dialog.NewInformation("Connection",
					"A connection is required\nplease try again.", win)
				d.SetOnClosed(func() {
					showLogin(a)
				})
				d.Show()
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
	ver := a.Preferences().Int("app.version")
	if ver < 1 {
		d := dialog.NewInformation("Welcome",
			"To use this app you need a\nrocrail server running.\n\nPlease enter the details\non the next screen", win)
		d.SetOnClosed(func() {
			showLogin(a)
		})
		d.Show()
		a.Preferences().SetInt("app.version", 1)
	} else {
		host, _ := userPref(a)
		if host == "" {
			go showLogin(a)
		} else {
			connect()
		}
	}

	c := newCon()
	cons = []*con{c}
	grid = container.NewGridWithColumns(1, c.makeUI())
	win.SetContent(grid)

	win.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("New", addControl),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Disconnect", func() {
				conn.Disconnect()
				showLogin(fyne.CurrentApp())
			}))))
	win.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: desktop.ControlModifier,
	}, func(_ fyne.Shortcut) {
		addControl()
	})
	win.ShowAndRun()
}
