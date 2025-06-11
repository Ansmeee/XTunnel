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
	ui     *UI
}

type UI struct {
	sidebar *Sidebar
	editor  *Editor
}

func NewWindow() *Window {
	window := new(app.Window)
	window.Option(app.Title("XTunnel"))
	window.Option(app.Size(unit.Dp(900), unit.Dp(500)))
	w := &Window{
		window: window,
		th:     material.NewTheme(),
		ops:    &op.Ops{},
		ui:     &UI{},
	}
	
	w.RegisterUI()
	return w
}

func (w *Window) RegisterUI() {
	w.ui.sidebar = NewSidebar(w)
	w.ui.editor = NewEditor(w)
}

func (w *Window) Run() {
	go func() {
		for {
			switch e := w.window.Event().(type) {
			case app.DestroyEvent:
				return
			case app.FrameEvent:
				w.gtx = app.NewContext(w.ops, e)
				layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceStart}.Layout(w.gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return w.ui.sidebar.Layout()
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return w.ui.editor.Layout()
					}),
				)
				e.Frame(w.gtx.Ops)
			}
		}
	}()
	app.Main()
}
