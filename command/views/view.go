package views

import (
	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
)

type View struct {
	ui       cli.Ui
	colorize *colorstring.Colorize
}

func NewView(ui cli.Ui, color bool) View {
	return View{
		ui: ui,
		colorize: &colorstring.Colorize{
			Colors:  colorstring.DefaultColors,
			Disable: !color,
			Reset:   true,
		},
	}
}

func (v *View) output(s string) {
	v.ui.Output(s)
}
