package views

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"log"
	"strings"
	"xtunnel/service"
)

const ModeCreate = 1
const ModeEdit = 2

type Editor struct {
	window          *Window
	mode            int
	fileName        string
	originPassword  string
	passwordChanged bool
	configNameInput widget.Editor
	remoteIpInput   widget.Editor
	remotePortInput widget.Editor
	serverIpInput   widget.Editor
	serverPortInput widget.Editor
	usernameInput   widget.Editor
	passwordInput   widget.Editor
	saveButton      widget.Clickable
	deleteButton    widget.Clickable

	configNameInputWidget *InputWidget
	remoteIpInputWidget   *InputWidget
	remotePortInputWidget *InputWidget
	serverIpInputWidget   *InputWidget
	serverPortInputWidget *InputWidget
	usernameInputWidget   *InputWidget
	passwordInputWidget   *InputWidget
}

type InputWidget struct {
	Input    *Input
	Editor   widget.Editor
	ValidErr string
}

type Input struct {
	gtx         layout.Context
	th          *material.Theme
	e           *Editor
	label       string
	labelWidth  int
	hint        string
	hintColor   color.NRGBA
	editor      *widget.Editor
	width       int
	borderColor color.NRGBA
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

		configNameInputWidget: &InputWidget{
			Input:  &Input{},
			Editor: widget.Editor{},
		},
		remoteIpInputWidget: &InputWidget{
			Input:  &Input{},
			Editor: widget.Editor{},
		},
		remotePortInputWidget: &InputWidget{
			Input:  &Input{},
			Editor: widget.Editor{},
		},
		serverIpInputWidget: &InputWidget{
			Input:  &Input{},
			Editor: widget.Editor{},
		},
		serverPortInputWidget: &InputWidget{
			Input:  &Input{},
			Editor: widget.Editor{},
		},
		usernameInputWidget: &InputWidget{
			Input:  &Input{},
			Editor: widget.Editor{},
		},
		passwordInputWidget: &InputWidget{
			Input:  &Input{},
			Editor: widget.Editor{},
		},
	}
	if w.ui.sidebar.SelectedItem != nil {
		editor.SwitchEditMode()
	}

	return editor
}

func (e *Editor) Layout() layout.Dimensions {
	th := e.window.th
	gtx := e.window.gtx

	gtx.Constraints = layout.Exact(image.Pt(560, gtx.Constraints.Max.Y))
	return layout.Inset{Left: unit.Dp(10), Right: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
				e.configNameInputWidget.Input = &Input{
					gtx:         gtx,
					th:          th,
					e:           e,
					label:       "配置名称：",
					labelWidth:  80,
					hint:        "请输入配置名称",
					hintColor:   color.NRGBA{R: 169, G: 169, B: 169, A: 255},
					editor:      &e.configNameInput,
					width:       gtx.Constraints.Max.X,
					borderColor: color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
				}

				return e.configNameInputWidget.Input.Layout()
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
						e.remoteIpInputWidget.Input = &Input{
							gtx:         gtx,
							th:          th,
							e:           e,
							label:       "主机IP：",
							labelWidth:  80,
							hint:        "请输入主机IP",
							hintColor:   color.NRGBA{R: 169, G: 169, B: 169, A: 255},
							editor:      &e.remoteIpInput,
							width:       340,
							borderColor: color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
						}

						if e.remoteIpInputWidget.ValidErr != "" {
							e.remoteIpInputWidget.Input.borderColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
							e.remoteIpInputWidget.Input.hintColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
						}

						return e.remoteIpInputWidget.Input.Layout()
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						e.remotePortInputWidget.Input = &Input{
							gtx:         gtx,
							th:          th,
							e:           e,
							label:       "端口：",
							labelWidth:  60,
							hint:        "请输入主机端口",
							hintColor:   color.NRGBA{R: 169, G: 169, B: 169, A: 255},
							editor:      &e.remotePortInput,
							width:       190,
							borderColor: color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
						}

						if e.remotePortInputWidget.ValidErr != "" {
							e.remotePortInputWidget.Input.borderColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
							e.remotePortInputWidget.Input.hintColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
						}

						return e.remotePortInputWidget.Input.Layout()
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
						e.serverIpInputWidget.Input = &Input{
							gtx:         gtx,
							th:          th,
							e:           e,
							label:       "主机IP：",
							labelWidth:  80,
							hint:        "请输入主机IP",
							hintColor:   color.NRGBA{R: 169, G: 169, B: 169, A: 255},
							editor:      &e.serverIpInput,
							width:       340,
							borderColor: color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
						}

						if e.serverIpInputWidget.ValidErr != "" {
							e.serverIpInputWidget.Input.borderColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
							e.serverIpInputWidget.Input.hintColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
						}

						return e.serverIpInputWidget.Input.Layout()
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						e.serverPortInputWidget.Input = &Input{
							gtx:         gtx,
							th:          th,
							e:           e,
							label:       "端口：",
							labelWidth:  60,
							hint:        "请输入主机端口",
							hintColor:   color.NRGBA{R: 169, G: 169, B: 169, A: 255},
							editor:      &e.serverPortInput,
							width:       190,
							borderColor: color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
						}

						if e.serverPortInputWidget.ValidErr != "" {
							e.serverPortInputWidget.Input.borderColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
							e.serverPortInputWidget.Input.hintColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
						}

						return e.serverPortInputWidget.Input.Layout()
					}),
				)

			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e.usernameInputWidget.Input = &Input{
					gtx:         gtx,
					th:          th,
					e:           e,
					label:       "用户名：",
					labelWidth:  80,
					hint:        "请输入用户名",
					hintColor:   color.NRGBA{R: 169, G: 169, B: 169, A: 255},
					editor:      &e.usernameInput,
					width:       gtx.Constraints.Max.X,
					borderColor: color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
				}

				if e.usernameInputWidget.ValidErr != "" {
					e.usernameInputWidget.Input.borderColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
					e.usernameInputWidget.Input.hintColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
				}

				return e.usernameInputWidget.Input.Layout()
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if _, c := e.passwordInput.Update(gtx); c {
					e.passwordChanged = e.passwordInput.Text() != e.originPassword
				}

				e.passwordInputWidget.Input = &Input{
					gtx:         gtx,
					th:          th,
					e:           e,
					label:       "密码：",
					labelWidth:  80,
					hint:        "请输入密码",
					hintColor:   color.NRGBA{R: 169, G: 169, B: 169, A: 255},
					editor:      &e.passwordInput,
					width:       gtx.Constraints.Max.X,
					borderColor: color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
				}
				if e.passwordInputWidget.ValidErr != "" {
					e.passwordInputWidget.Input.borderColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
					e.passwordInputWidget.Input.hintColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
				}
				return e.passwordInputWidget.Input.Layout()
			}),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: 20, Bottom: 20, Left: 50, Right: 50}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceSides}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    20,
								Bottom: 20,
								Left:   50,
								Right:  50,
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								if e.saveButton.Clicked(gtx) {
									e.OnSaveBtnClicked()
								}
								btn := material.Button(th, &e.saveButton, "保存")
								btn.Inset = layout.Inset{Top: 10, Bottom: 10, Left: 20, Right: 20}
								return btn.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if e.IsEditMode() {
								return layout.Inset{
									Top:    20,
									Bottom: 20,
									Left:   50,
									Right:  50,
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									if e.deleteButton.Clicked(gtx) {
										e.OnDelBtnClicked()
									}
									btn := material.Button(th, &e.deleteButton, "删除")
									btn.Background = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
									btn.Inset = layout.Inset{Top: 8, Bottom: 8, Left: 20, Right: 20}
									return btn.Layout(gtx)
								})
							}
							return layout.Dimensions{}
						}),
					)
				})
			}),
		)
	})
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
				Color:        e.borderColor,
				Width:        unit.Dp(1),
				CornerRadius: unit.Dp(4),
			}
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(5).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, e.editor, e.hint)
					editor.HintColor = e.hintColor
					return editor.Layout(gtx)
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

	password := e.originPassword
	if e.passwordChanged {
		password = e.passwordInput.Text()
	}

	cf := &service.ConfigFile{
		ConfigName: e.configNameInput.Text(),
		RemoteIP:   e.remoteIpInput.Text(),
		RemotePort: e.remotePortInput.Text(),
		ServerIP:   e.serverIpInput.Text(),
		ServerPort: e.serverPortInput.Text(),
		UserName:   e.usernameInput.Text(),
		Password:   password,
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
	hasErr := false
	e.remoteIpInputWidget.ValidErr = ""
	if e.remoteIpInput.Text() == "" {
		e.remoteIpInputWidget.ValidErr = "remote ip is empty"
		hasErr = true
	}

	e.remotePortInputWidget.ValidErr = ""
	if e.remotePortInput.Text() == "" {
		e.remotePortInputWidget.ValidErr = "remote port is empty"
		hasErr = true
	}

	e.serverIpInputWidget.ValidErr = ""
	if e.serverIpInput.Text() == "" {
		e.serverIpInputWidget.ValidErr = "server ip is empty"
		hasErr = true
	}
	e.serverPortInputWidget.ValidErr = ""
	if e.serverPortInput.Text() == "" {
		e.serverPortInputWidget.ValidErr = "server port is empty"
		hasErr = true
	}

	e.usernameInputWidget.ValidErr = ""
	if e.usernameInput.Text() == "" {
		e.usernameInputWidget.ValidErr = "username is empty"
		hasErr = true
	}

	e.passwordInputWidget.ValidErr = ""
	if e.passwordInput.Text() == "" {
		e.passwordInputWidget.ValidErr = "password is empty"
		hasErr = true
	}

	if hasErr {
		return fmt.Errorf("form validation error")
	}

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
	e.serverIpInput.SetText(config.ServerIP)
	e.serverPortInput.SetText(config.ServerPort)
	e.usernameInput.SetText(config.UserName)
	e.originPassword = config.Password
	e.passwordInput.SetText(strings.Repeat("*", 10))
}

func (e *Editor) IsCreateMode() bool {
	return e.mode == ModeCreate
}

func (e *Editor) IsEditMode() bool {
	return e.mode == ModeEdit
}
