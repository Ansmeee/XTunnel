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

	var listState = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}

	cf := &service.TunnelConfig{}
	files, err := cf.LoadConfigFile()
	if err != nil {
		return err
	}

	switchs := make([]widget.Bool, len(files))
	for i, file := range files {
		fmt.Println(file.ConfigName)
		switchs[i] = widget.Bool{Value: file.Switch}
	}

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

				form.Save()
			}

			layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								dims := material.Body1(th, "固定宽度内容").Layout(gtx)
								return layout.Dimensions{
									Size:     image.Pt(gtx.Dp(200), dims.Size.Y),
									Baseline: dims.Baseline,
								}
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.List(th, listState).Layout(gtx, len(files), func(gtx layout.Context, index int) layout.Dimensions {
									cfg := files[index]

									if switchs[index].Update(gtx) {
										tunnel := service.Tunnel{
											Username:   cfg.UserName,
											Password:   cfg.Password,
											LocalAddr:  fmt.Sprintf("127.0.0.1:%s", cfg.RemotePort),
											ServerAddr: fmt.Sprintf("%s:%s", cfg.SSHIp, cfg.SSHPort),
											RemoteAddr: fmt.Sprintf("%s:%s", cfg.RemoteIP, cfg.RemotePort),
										}

										if switchs[index].Value == true {
											go tunnel.Start()
										} else {
											go tunnel.Stop()
										}
									}

									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return material.Body1(th, cfg.ConfigName).Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return material.Switch(th, &switchs[index], "").Layout(gtx)
										}),
									)
								})
							}),
						)
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

func readConfigFiles() []*service.TunnelConfig {
	cf := &service.TunnelConfig{}
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
	config := service.TunnelConfig{
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

func (f *submitForm) Create() {
	tunnel := &service.Tunnel{
		Username:   f.UserName,
		Password:   f.Password,
		LocalAddr:  fmt.Sprintf("127.0.0.1:%s", f.Port),
		RemoteAddr: fmt.Sprintf("%s:%s", f.HostIp, f.Port),
		ServerAddr: fmt.Sprintf("%s:%s", f.SSHIp, f.SSHPort),
	}

	go tunnel.Start()
}
