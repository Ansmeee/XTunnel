package views

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Window struct {
	window *app.Window
	th     *material.Theme
	ops    *op.Ops
	gtx    layout.Context
}

func NewWindow() *Window {
	window := new(app.Window)
	window.Option(app.Title("XTunnel"))
	window.Option(app.Size(unit.Dp(900), unit.Dp(500)))
	return &Window{
		window: window,
		th:     material.NewTheme(),
		ops:    &op.Ops{},
	}
}

func (w *Window) Run() {
	go func() {
		for {
			switch e := w.window.Event().(type) {
			case app.DestroyEvent:
				return
			case app.FrameEvent:
				w.gtx = app.NewContext(w.ops, e)
				w.Layout()
				e.Frame(w.gtx.Ops)
			}
		}
	}()
	app.Main()
}

func (w *Window) Layout() {
	layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceStart}.Layout(w.gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.NewSidebar().Layout()
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.NewEditor().Layout()
		}),
	)
}

func (w *Window) NewSidebar() *Sidebar {
	return &Sidebar{window: w}
}

func (w *Window) NewEditor() *Editor {
	return &Editor{window: w}
}
