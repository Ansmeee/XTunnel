package views

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"log"
	"xtunnel/service"
)

type Sidebar struct {
	SelectedItem  *SidebarItem
	window        *Window
	items         []*SidebarItem
	tunnelManager *service.TunnelManager
	listState     *widget.List
	createBtn     *widget.Clickable
}

type SidebarItem struct {
	config       *service.ConfigFile
	tunnel       *service.Tunnel
	clickWidget  widget.Clickable
	switchWidget widget.Bool
}

func (s *Sidebar) LoadSidebarItems() error {
	cf := &service.ConfigFile{}
	files, err := cf.LoadConfigFile()
	if err != nil {
		log.Panicf("loading config file err: %s", err.Error())
		return err
	}

	tunnelManager := service.NewTunnelManager()
	items := make([]*SidebarItem, len(files))
	for i, file := range files {
		_, err := tunnelManager.AddTunnel(
			file.Identifier,
			&service.TunnelConfig{
				Username:   file.UserName,
				Password:   file.Password,
				LocalAddr:  fmt.Sprintf("127.0.0.1:%s", file.RemotePort),
				ServerAddr: fmt.Sprintf("%s:%s", file.ServerIP, file.ServerPort),
				RemoteAddr: fmt.Sprintf("%s:%s", file.RemoteIP, file.RemotePort),
			},
		)

		if err != nil {
			log.Printf("add tunnel err: %s", err.Error())
			continue
		}

		items[i] = &SidebarItem{
			config:       file,
			switchWidget: widget.Bool{Value: false},
			clickWidget:  widget.Clickable{},
		}
	}

	s.items = items
	s.tunnelManager = tunnelManager

	if len(items) > 0 {
		s.SelectedItem = items[0]
	}

	return nil
}

func NewSidebar(w *Window) *Sidebar {
	sidebar := &Sidebar{
		window:    w,
		createBtn: &widget.Clickable{},
		listState: &widget.List{List: layout.List{Axis: layout.Vertical}},
	}

	if err := sidebar.LoadSidebarItems(); err != nil {
		log.Printf("LoadSidebarItems err: %s", err.Error())
	}

	return sidebar
}

func (s *Sidebar) Layout() layout.Dimensions {
	th := s.window.th
	gtx := s.window.gtx

	gtx.Constraints = layout.Exact(image.Pt(300, gtx.Constraints.Max.Y))
	return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(10), Bottom: unit.Dp(40)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Body1(th, "配置列表").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:    0,
							Bottom: 0,
							Left:   0,
							Right:  10,
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							if s.createBtn.Clicked(gtx) {
								s.SelectedItem = nil
								s.window.ui.editor.SwitchCreateMode()
							}
							btn := material.Button(th, s.createBtn, "新增")
							btn.Inset = layout.Inset{Top: 2, Bottom: 2, Left: 10, Right: 10}
							return btn.Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: 10}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.List(th, s.listState).Layout(gtx, len(s.items), func(gtx layout.Context, index int) layout.Dimensions {
					item := s.items[index]

					bg := color.NRGBA{}
					if item.clickWidget.Clicked(gtx) {
						s.SelectedItem = item
						s.window.ui.editor.SwitchEditMode()
					}

					if s.SelectedItem != nil {
						if s.SelectedItem.config.ConfigName == item.config.ConfigName {
							bg = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
						}

						if item.clickWidget.Hovered() {
							bg = color.NRGBA{R: 220, G: 220, B: 220, A: 255}
						}
					}

					content := func(gtx layout.Context) layout.Dimensions {
						return layout.Stack{}.Layout(gtx,
							layout.Expanded(func(gtx layout.Context) layout.Dimensions {
								paint.Fill(gtx.Ops, bg)
								return layout.Dimensions{Size: gtx.Constraints.Min}
							}),
							layout.Stacked(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										gtx.Constraints = layout.Exact(image.Pt(210, 30))
										return layout.UniformInset(5).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return material.Body1(th, item.config.ConfigName).Layout(gtx)
										})
									}),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										if item.switchWidget.Update(gtx) {
											if item.switchWidget.Value == true {
												if err := s.tunnelManager.StartTunnel(item.config.Identifier); err != nil {
													log.Printf("start tunnel err: %s", err.Error())
												}
											} else {
												if err := s.tunnelManager.StopTunnel(item.config.Identifier); err != nil {
													log.Printf("stop tunnel err: %s", err.Error())
												}
											}
										}
										return layout.UniformInset(5).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											sw := material.Switch(th, &item.switchWidget, "")
											sw.Color.Enabled = color.NRGBA{R: 0, G: 128, B: 0, A: 255}
											return sw.Layout(gtx)
										})
									}),
								)
							}),
						)
					}

					return item.clickWidget.Layout(gtx, content)
				})
			}),
		)
	})
}
