package service

import (
	"context"

	"github.com/Mersad-Moghaddam/syskit/internal/plugin"
)

type Plugin struct{}

func NewPlugin() *Plugin                                  { return &Plugin{} }
func (*Plugin) List(dirs []string) ([]plugin.Info, error) { return plugin.Discover(dirs) }
func (*Plugin) Inspect(dirs []string, name string) (*plugin.Info, error) {
	return plugin.Inspect(dirs, name)
}
func (*Plugin) Run(ctx context.Context, dirs []string, name string) (any, error) {
	return plugin.Run(ctx, dirs, name)
}
