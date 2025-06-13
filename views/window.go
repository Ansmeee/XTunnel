package views

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"image/color"
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
	window.Option(app.MaxSize(unit.Dp(900), unit.Dp(500)))
	window.Option(app.MinSize(unit.Dp(900), unit.Dp(500)))
	th := material.NewTheme()
	th.TextSize = unit.Sp(14)
	w := &Window{
		window: window,
		th:     th,
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
				layout.Stack{}.Layout(w.gtx,
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						paint.Fill(gtx.Ops, color.NRGBA{R: 240, G: 240, B: 240, A: 255})
						return layout.Dimensions{Size: gtx.Constraints.Min}
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: 20, Bottom: 20, Left: 20, Right: 20}.Layout(w.gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(w.gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return w.ui.sidebar.Layout()
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return w.ui.editor.Layout()
								}),
							)
						})
					}),
				)
				e.Frame(w.gtx.Ops)
			}
		}
	}()
	app.Main()
}
