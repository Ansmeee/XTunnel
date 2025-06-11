package views

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"xtunnel/service"
)

const ModeCreate = 1
const ModeEdit = 2

type Editor struct {
	window          *Window
	mode            int
	configNameInput widget.Editor
	remoteIpInput   widget.Editor
	remotePortInput widget.Editor
	serverIpInput   widget.Editor
	serverPortInput widget.Editor
	usernameInput   widget.Editor
	passwordInput   widget.Editor
	saveButton      widget.Clickable
}

func NewEditor(w *Window) *Editor {
	return &Editor{
		window:          w,
		configNameInput: widget.Editor{},
		remoteIpInput:   widget.Editor{},
		remotePortInput: widget.Editor{},
		serverIpInput:   widget.Editor{},
		serverPortInput: widget.Editor{},
		usernameInput:   widget.Editor{},
		passwordInput:   widget.Editor{},
		saveButton:      widget.Clickable{},
	}
}

func (e *Editor) Layout() layout.Dimensions {
	th := e.window.th
	gtx := e.window.gtx

	gtx.Constraints = layout.Exact(image.Pt(600, gtx.Constraints.Max.Y))
	return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Subtitle1(th, "隧道配置")
				t.Alignment = text.Middle
				return t.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Body1(th, "配置名称：").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Editor(th, &e.configNameInput, "请输入配置名称").Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 30}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Subtitle1(th, "主机配置")
				t.Alignment = text.Start
				return t.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, "主机IP：").Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Editor(th, &e.remoteIpInput, "请输入主机IP").Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, "端口：").Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Editor(th, &e.remotePortInput, "请输入端口").Layout(gtx)
							}),
						)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 30}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Subtitle1(th, "SSH代理配置")
				t.Alignment = text.Start
				return t.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, "代理主机IP：").Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Editor(th, &e.serverIpInput, "请输入代理主机IP").Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, "端口：").Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Editor(th, &e.serverPortInput, "请输入端口").Layout(gtx)
							}),
						)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Body1(th, "用户名：").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Editor(th, &e.usernameInput, "请输入用户名").Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Body1(th, "密码：").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Editor(th, &e.passwordInput, "请输入密码").Layout(gtx)
					}),
				)
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    20,
					Bottom: 20,
					Left:   0,
					Right:  0,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &e.saveButton, "保存")
					return btn.Layout(gtx)
				})
			}),
		)
	})
}

func (e *Editor) SwitchCreateMode() {
	if !e.IsCreateMode() {
		e.mode = ModeCreate
		e.setCurItem()
	}
}

func (e *Editor) SwitchEditMode() {
	e.mode = ModeEdit
	e.setCurItem()
}

func (e *Editor) setCurItem() {
	config := &service.ConfigFile{}
	if e.IsEditMode() {
		curItem := e.window.ui.sidebar.SelectedItem
		if curItem == nil {
			return
		}
		config = curItem.config
	}

	e.configNameInput.SetText(config.ConfigName)
	e.remoteIpInput.SetText(config.RemoteIP)
	e.remotePortInput.SetText(config.RemotePort)
	e.serverIpInput.SetText(config.RemoteIP)
	e.serverPortInput.SetText(config.RemotePort)
	e.usernameInput.SetText(config.UserName)
	e.passwordInput.SetText(config.Password)
}

func (e *Editor) IsCreateMode() bool {
	return e.mode == ModeCreate
}

func (e *Editor) IsEditMode() bool {
	return e.mode == ModeEdit
}
