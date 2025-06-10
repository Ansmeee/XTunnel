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
	"xtunnel/service"
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

type TunnelItem struct {
	Config *service.ConfigFile
	Switch widget.Bool
	Click  widget.Clickable
	Tunnel *service.Tunnel
}

func run(window *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops

	var configNameInput widget.Editor
	var hostInput widget.Editor
	var portInput widget.Editor
	var sshHostInput widget.Editor
	var sshPortInput widget.Editor
	var userNameInput widget.Editor
	var passwordInput widget.Editor
	var submitBtn widget.Clickable

	var createBtn widget.Clickable

	var listState = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}

	cf := &service.ConfigFile{}
	files, err := cf.LoadConfigFile()
	if err != nil {
		return err
	}

	items := make([]*TunnelItem, len(files))
	for i, file := range files {
		items[i] = &TunnelItem{
			Config: file,
			Switch: widget.Bool{Value: false},
			Click:  widget.Clickable{},
			Tunnel: service.NewTunnel(&service.TunnelConfig{
				Username:   file.UserName,
				Password:   file.Password,
				LocalAddr:  fmt.Sprintf("127.0.0.1:%s", file.RemotePort),
				ServerAddr: fmt.Sprintf("%s:%s", file.SSHIp, file.SSHPort),
				RemoteAddr: fmt.Sprintf("%s:%s", file.RemoteIP, file.RemotePort),
			}),
		}
	}

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			if submitBtn.Clicked(gtx) {
				form := &submitForm{
					HostIp:   hostInput.Text(),
					Port:     portInput.Text(),
					SSHIp:    sshHostInput.Text(),
					SSHPort:  sshPortInput.Text(),
					UserName: userNameInput.Text(),
					Password: passwordInput.Text(),
				}

				form.Save()
			}

			layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints = layout.Exact(image.Pt(300, gtx.Constraints.Max.Y))
					return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Body1(th, "配置列表").Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return createBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return material.Body1(th, "新增").Layout(gtx)
										})
									}),
								)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Spacer{Height: 10}.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.List(th, listState).Layout(gtx, len(files), func(gtx layout.Context, index int) layout.Dimensions {
									item := items[index]

									if item.Click.Clicked(gtx) {
										configNameInput.SetText(item.Config.ConfigName)
										hostInput.SetText(item.Config.RemoteIP)
										portInput.SetText(item.Config.RemotePort)
										sshHostInput.SetText(item.Config.SSHIp)
										sshPortInput.SetText(item.Config.SSHPort)
										userNameInput.SetText(item.Config.UserName)
										passwordInput.SetText(item.Config.Password)
									}

									content := func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Body1(th, item.Config.ConfigName).Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												if item.Switch.Update(gtx) {
													if item.Switch.Value == true {
														go func() {
															if err := item.Tunnel.Start(); err != nil {
																log.Fatalf("隧道启动失败: %v", err)
															}
														}()
													} else {
														go func() {
															item.Tunnel.Stop()
														}()
													}
												}
												return material.Switch(th, &item.Switch, "").Layout(gtx)
											}),
										)
									}

									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return item.Click.Layout(gtx, content)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return layout.Spacer{Height: 10}.Layout(gtx)
										}),
									)
								})
							}),
						)
					})
				}),

				// 右侧自适应部分
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { // Flexed(1)表示占满剩余空间
					return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
										return material.Editor(th, &configNameInput, "请输入配置名称").Layout(gtx)
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
												return material.Editor(th, &hostInput, "请输入主机IP").Layout(gtx)
											}),
										)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Body1(th, "端口：").Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Editor(th, &portInput, "请输入端口").Layout(gtx)
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
												return material.Editor(th, &sshHostInput, "请输入代理主机IP").Layout(gtx)
											}),
										)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Body1(th, "端口：").Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Editor(th, &sshPortInput, "请输入端口").Layout(gtx)
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
										return material.Editor(th, &userNameInput, "请输入用户名").Layout(gtx)
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
										return material.Editor(th, &passwordInput, "请输入密码").Layout(gtx)
									}),
								)
							}),

							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Top:    20,
									Bottom: 20,
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

func readConfigFiles() []*service.ConfigFile {
	cf := &service.ConfigFile{}
	files, err := cf.LoadConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	return files
}

type submitForm struct {
	HostIp   string
	Port     string
	SSHIp    string
	SSHPort  string
	UserName string
	Password string
}

func (f *submitForm) Save() {
	config := service.ConfigFile{
		RemoteIP:   f.HostIp,
		RemotePort: f.Port,
		SSHIp:      f.SSHIp,
		SSHPort:    f.SSHPort,
		UserName:   f.UserName,
		Password:   f.Password,
	}

	err := config.SaveConfigFile()
	if err != nil {
		panic(err)
	}
}
