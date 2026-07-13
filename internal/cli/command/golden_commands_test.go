package command

import (
	"bytes"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/golden"
	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type goldenCPUService struct{ info *model.CPUInfo }

func (s goldenCPUService) Collect() (*model.CPUInfo, error)             { return s.info, nil }
func (s goldenCPUService) Sample(time.Duration) (*model.CPUInfo, error) { return s.info, nil }

type goldenMemoryService struct{ info *model.MemoryInfo }

func (s goldenMemoryService) Collect() (*model.MemoryInfo, error) { return s.info, nil }

type goldenDiskService struct{ info *model.DiskInfo }

func (s goldenDiskService) Collect() (*model.DiskInfo, error)             { return s.info, nil }
func (s goldenDiskService) Sample(time.Duration) (*model.DiskInfo, error) { return s.info, nil }

func TestCoreCommandGoldens(t *testing.T) {
	util := 15.0
	available, used := uint64(60), uint64(40)
	total, avail := uint64(100), uint64(60)
	pct := 40.0
	inodes, free := uint64(100), uint64(80)
	for _, tc := range []struct {
		name string
		cmd  func(string) *cobra.Command
	}{
		{"cpu", func(format string) *cobra.Command {
			return NewCPUCmd(goldenCPUService{&model.CPUInfo{LogicalCores: 1, Model: "Fixture CPU", Architecture: "amd64", Times: []model.CPUTime{{CPUID: "all", User: 10, System: 5, Idle: 85, Total: 100, Utilization: &util}}}}, CPUOptions{Format: func() string { return format }, NoHeader: func() bool { return false }, Color: func() bool { return false }})
		}},
		{"memory", func(format string) *cobra.Command {
			return NewMemoryCmd(goldenMemoryService{&model.MemoryInfo{TotalBytes: 100, UsedBytes: &used, AvailableBytes: &available, FreeBytes: 10}}, MemoryOptions{Format: func() string { return format }, NoHeader: func() bool { return false }, Color: func() bool { return false }})
		}},
		{"disk", func(format string) *cobra.Command {
			return NewDiskCmd(goldenDiskService{&model.DiskInfo{Mounts: []model.MountInfo{{Source: "/dev/test", FilesystemType: "ext4", MountPoint: "/", TotalBytes: &total, UsedBytes: &used, AvailableBytes: &avail, UsePercent: &pct}}}}, DiskOptions{Format: func() string { return format }, NoHeader: func() bool { return false }, Color: func() bool { return false }})
		}},
		{"filesystem", func(format string) *cobra.Command {
			mount := model.MountInfo{Source: "/dev/test", FilesystemType: "ext4", MountPoint: "/", TotalInodes: &inodes, FreeInodes: &free, Options: []string{"rw"}}
			return NewFilesystemCmd(goldenDiskService{&model.DiskInfo{Mounts: []model.MountInfo{mount}}}, FilesystemOptions{Format: func() string { return format }, NoHeader: func() bool { return false }, Color: func() bool { return false }})
		}},
	} {
		for _, format := range []string{"table", "json"} {
			t.Run(tc.name+"_"+format, func(t *testing.T) {
				cmd := tc.cmd(format)
				var out bytes.Buffer
				cmd.SetOut(&out)
				cmd.SetArgs([]string{})
				if tc.name == "cpu" {
					cmd.SetArgs([]string{"--interval", "1ms"})
				}
				assert.NoError(t, cmd.Execute())
				golden.Assert(t, out.Bytes(), tc.name+"_"+format+".golden")
			})
		}
	}
}
