package service

import "github.com/Mersad-Moghaddam/syskit/internal/plugin"

type Plugin struct{}

func NewPlugin() *Plugin                                  { return &Plugin{} }
func (*Plugin) List(dirs []string) ([]plugin.Info, error) { return plugin.Discover(dirs) }
