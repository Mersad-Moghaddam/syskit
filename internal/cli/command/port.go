package command

import (
	"fmt"
	"strings"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
	"github.com/Mersad-Moghaddam/syskit/internal/service"
	"github.com/spf13/cobra"
)

type PortService interface {
	List(service.PortOptions) (*model.PortInfo, error)
}
type PortOptions struct {
	Format   func() string
	NoHeader func() bool
}

func NewPortCmd(s PortService, o PortOptions) *cobra.Command {
	var listening bool
	var protocol string
	var port int
	var pid int
	cmd := &cobra.Command{Use: "ports", Short: "Show TCP and UDP sockets", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		info, err := s.List(service.PortOptions{Listening: listening, Protocol: protocol, Port: port, PID: pid})
		if err != nil {
			return fmt.Errorf("collecting ports: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), portTable(info))
		}
		return r.Render(c.OutOrStdout(), info)
	}}
	cmd.Flags().BoolVar(&listening, "listening", false, "show listening sockets only")
	cmd.Flags().StringVar(&protocol, "protocol", "", "filter by protocol (tcp, tcp6, udp, udp6, unix)")
	cmd.Flags().IntVar(&port, "port", 0, "filter by local port")
	cmd.Flags().IntVar(&pid, "pid", 0, "filter by owning process ID")
	return cmd
}
func portTable(info *model.PortInfo) render.Table {
	t := render.Table{Headers: []string{"PROTO", "LOCAL", "REMOTE", "STATE", "INODE", "PID", "COMMAND"}}
	for _, s := range info.Sockets {
		local, remote := socketEndpoint(s.LocalAddress, s.LocalPort), socketEndpoint(s.RemoteAddress, s.RemotePort)
		var pids, commands []string
		for _, owner := range s.Owners {
			pids = append(pids, fmt.Sprint(owner.PID))
			commands = append(commands, owner.Command)
		}
		t.Rows = append(t.Rows, []string{s.Protocol, local, remote, s.State, fmt.Sprint(s.Inode), strings.Join(pids, ","), strings.Join(commands, ", ")})
	}
	return t
}
func socketEndpoint(address string, port uint16) string {
	if port == 0 {
		if address == "" {
			return "*"
		}
		return address
	}
	return fmt.Sprintf("%s:%d", address, port)
}
