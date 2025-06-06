package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"log"
	"os"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Title("XTunnel"))
		window.Option(app.Size(unit.Dp(800), unit.Dp(500)))
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops

	var hostInput widget.Editor
	var portInput widget.Editor
	var sshHostInput widget.Editor
	var sshPortInput widget.Editor
	var userNameInput widget.Editor
	var passwordInput widget.Editor
	var submitBtn widget.Clickable
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			if submitBtn.Clicked(gtx) {

				fmt.Println("submitBtn clicked")
				form := &submitForm{
					HostIp:   hostInput.Text(),
					Port:     portInput.Text(),
					SSHIp:    sshHostInput.Text(),
					SSHPort:  sshPortInput.Text(),
					UserName: userNameInput.Text(),
					Password: passwordInput.Text(),
				}

				form.Submit()
			}

			layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						dims := material.Body1(th, "固定宽度内容").Layout(gtx)
						return layout.Dimensions{
							Size:     image.Pt(gtx.Dp(200), dims.Size.Y),
							Baseline: dims.Baseline,
						}
					})
				}),

				// 右侧自适应部分
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { // Flexed(1)表示占满剩余空间
					return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								t := material.Subtitle1(th, "隧道配置")
								t.Alignment = text.Middle
								return t.Layout(gtx)
							}),

							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								t := material.Subtitle1(th, "主机配置")
								t.Alignment = text.Start
								return t.Layout(gtx)
							}),

							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Body1(th, "主机IP：").Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Editor(th, &hostInput, "请输入主机IP").Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Body1(th, "端口：").Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Editor(th, &portInput, "请输入端口").Layout(gtx)
									}),
								)
							}),

							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								t := material.Subtitle1(th, "SSH代理配置")
								t.Alignment = text.Start
								return t.Layout(gtx)
							}),

							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Body1(th, "代理主机IP：").Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Editor(th, &sshHostInput, "请输入代理主机IP").Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Body1(th, "端口：").Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Editor(th, &sshPortInput, "请输入端口").Layout(gtx)
									}),
								)
							}),

							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Body1(th, "用户名：").Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Editor(th, &userNameInput, "请输入用户名").Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Body1(th, "密码：").Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Editor(th, &passwordInput, "请输入密码").Layout(gtx)
									}),
								)
							}),

							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Top:    20,
									Bottom: 50,
									Left:   20,
									Right:  20,
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(th, &submitBtn, "保存")
									return btn.Layout(gtx)
								})
							}),
						)
					})
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

type submitForm struct {
	HostIp   string
	Port     string
	SSHIp    string
	SSHPort  string
	UserName string
	Password string
}

func (f *submitForm) Submit() {

}
