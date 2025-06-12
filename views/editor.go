package views

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"log"
	"xtunnel/service"
)

const ModeCreate = 1
const ModeEdit = 2

type Editor struct {
	window          *Window
	mode            int
	fileName        string
	configNameInput widget.Editor
	remoteIpInput   widget.Editor
	remotePortInput widget.Editor
	serverIpInput   widget.Editor
	serverPortInput widget.Editor
	usernameInput   widget.Editor
	passwordInput   widget.Editor
	saveButton      widget.Clickable
	deleteButton    widget.Clickable
}

func NewEditor(w *Window) *Editor {
	editor := &Editor{
		window:          w,
		configNameInput: widget.Editor{},
		remoteIpInput:   widget.Editor{},
		remotePortInput: widget.Editor{},
		serverIpInput:   widget.Editor{},
		serverPortInput: widget.Editor{},
		usernameInput:   widget.Editor{},
		passwordInput:   widget.Editor{},
		saveButton:      widget.Clickable{},
		deleteButton:    widget.Clickable{},
	}
	if w.ui.sidebar.SelectedItem != nil {
		editor.SwitchEditMode()
	}

	return editor
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
				t := material.Body1(th, "基本配置")
				t.TextSize = unit.Sp(12)
				t.Alignment = text.Start
				return t.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				input := &Input{
					gtx:        gtx,
					th:         th,
					e:          e,
					label:      "配置名称：",
					labelWidth: 80,
					hint:       "请输入配置名称",
					editor:     &e.configNameInput,
					width:      gtx.Constraints.Max.X,
				}
				return input.Layout()
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 30}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Subtitle1(th, "主机配置")
				t.TextSize = unit.Sp(12)
				t.Alignment = text.Start
				return t.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						input := &Input{
							gtx:        gtx,
							th:         th,
							e:          e,
							label:      "主机IP：",
							labelWidth: 80,
							hint:       "请输入主机IP",
							editor:     &e.remoteIpInput,
							width:      380,
						}

						return input.Layout()
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						input := &Input{
							gtx:        gtx,
							th:         th,
							e:          e,
							label:      "端口：",
							labelWidth: 60,
							hint:       "请输入主机端口",
							editor:     &e.remotePortInput,
							width:      200,
						}

						return input.Layout()
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 30}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t := material.Subtitle1(th, "SSH代理配置")
				t.Alignment = text.Start
				t.TextSize = unit.Sp(12)
				return t.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						input := &Input{
							gtx:        gtx,
							th:         th,
							e:          e,
							label:      "主机IP：",
							labelWidth: 80,
							hint:       "请输入主机IP",
							editor:     &e.serverIpInput,
							width:      380,
						}

						return input.Layout()
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						input := &Input{
							gtx:        gtx,
							th:         th,
							e:          e,
							label:      "端口：",
							labelWidth: 60,
							hint:       "请输入主机端口",
							editor:     &e.serverPortInput,
							width:      200,
						}

						return input.Layout()
					}),
				)

			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				input := &Input{
					gtx:        gtx,
					th:         th,
					e:          e,
					label:      "用户名：",
					labelWidth: 80,
					hint:       "请输入用户名",
					editor:     &e.usernameInput,
					width:      gtx.Constraints.Max.X,
				}
				return input.Layout()
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				input := &Input{
					gtx:        gtx,
					th:         th,
					e:          e,
					label:      "密码：",
					labelWidth: 80,
					hint:       "请输入密码",
					editor:     &e.passwordInput,
					width:      gtx.Constraints.Max.X,
				}
				return input.Layout()
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: 20, Bottom: 20, Left: 50, Right: 50}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if e.saveButton.Clicked(gtx) {
								e.OnSaveBtnClicked()
							}
							btn := material.Button(th, &e.saveButton, "保存")
							btn.Inset = layout.Inset{Top: 10, Bottom: 10, Left: 20, Right: 20}
							return btn.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if e.IsEditMode() {
								if e.deleteButton.Clicked(gtx) {
									e.OnDelBtnClicked()
								}
								btn := material.Button(th, &e.deleteButton, "删除")
								btn.Background = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
								btn.Inset = layout.Inset{Top: 8, Bottom: 8, Left: 20, Right: 20}
								return btn.Layout(gtx)
							}
							return layout.Dimensions{}
						}),
					)
				})
			}),
		)
	})
}

type Input struct {
	gtx        layout.Context
	th         *material.Theme
	e          *Editor
	label      string
	labelWidth int
	hint       string
	editor     *widget.Editor
	width      int
}

func (e *Input) Layout() layout.Dimensions {
	gtx := e.gtx
	th := e.th

	gtx.Constraints = layout.Exact(image.Pt(e.width, 30))
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints = layout.Exact(image.Pt(e.labelWidth, 30))
			return layout.UniformInset(5).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Body1(th, e.label).Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			border := widget.Border{
				Color:        color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
				Width:        unit.Dp(1),
				CornerRadius: unit.Dp(4),
			}
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(5).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					e := material.Editor(th, e.editor, e.hint)
					return e.Layout(gtx)
				})
			})
		}),
	)
}

func (e *Editor) OnSaveBtnClicked() {
	if err := e.validateForm(); err != nil {
		log.Printf("form validation error: %s", err)
		return
	}

	cf := &service.ConfigFile{
		ConfigName: e.configNameInput.Text(),
		RemoteIP:   e.remoteIpInput.Text(),
		RemotePort: e.remotePortInput.Text(),
		ServerIP:   e.serverIpInput.Text(),
		ServerPort: e.serverPortInput.Text(),
		UserName:   e.usernameInput.Text(),
		Password:   e.passwordInput.Text(),
	}

	var err error

	if e.IsEditMode() {
		cf.FileName = e.fileName
		err = cf.UpdateConfigFile()
	} else {
		err = cf.SaveConfigFile()
	}

	if err != nil {
		log.Printf("save config error: %s", err)
		return
	}

	if err := e.window.ui.sidebar.LoadSidebarItems(); err != nil {
		log.Printf("load sidebar error: %s", err)
		return
	}

	e.SwitchEditMode()
}

func (e *Editor) OnDelBtnClicked() {
	cf := &service.ConfigFile{
		FileName: e.fileName,
	}

	if err := cf.DeleteConfigFile(); err != nil {
		log.Printf("delete config error: %s", err.Error())
		return
	}

	if err := e.window.ui.sidebar.LoadSidebarItems(); err != nil {
		log.Printf("load sidebar error: %s", err)
	}

	e.SwitchCreateMode()
}

func (e *Editor) validateForm() error {
	return nil
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
		e.fileName = config.FileName
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
