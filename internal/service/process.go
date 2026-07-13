package service

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

type ProcessCollector interface {
	Collect() (*model.ProcessList, error)
}
type ProcessOptions struct {
	Filters    []Filter
	Sort       string
	Reverse    bool
	Limit      int
	Containers bool
}
type Process struct{ collector ProcessCollector }

func NewProcess(c ProcessCollector) *Process { return &Process{c} }
func (s *Process) List(o ProcessOptions) (*model.ProcessList, error) {
	list, err := s.collector.Collect()
	if err != nil {
		return nil, err
	}
	applyMemoryPercentages(list)
	if o.Containers {
		list = containerProcesses(list)
	}
	return s.list(list, o)
}
func containerProcesses(list *model.ProcessList) *model.ProcessList {
	result := *list
	result.Processes = nil
	for _, process := range list.Processes {
		if process.ContainerID != "" {
			result.Processes = append(result.Processes, process)
		}
	}
	return &result
}
func (s *Process) Sample(interval time.Duration, o ProcessOptions) (*model.ProcessList, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("process sample interval must be positive")
	}
	before, err := s.collector.Collect()
	if err != nil {
		return nil, err
	}
	time.Sleep(interval)
	after, err := s.collector.Collect()
	if err != nil {
		return nil, err
	}
	applyProcessPercentages(before, after)
	if o.Containers {
		after = containerProcesses(after)
	}
	return s.list(after, o)
}
func (s *Process) list(list *model.ProcessList, o ProcessOptions) (*model.ProcessList, error) {
	items, err := FilterItems(list.Processes, o.Filters, map[string]func(model.Process) string{"pid": func(p model.Process) string { return strconv.Itoa(p.PID) }, "user": func(p model.Process) string {
		if p.User != "" {
			return p.User
		}
		return strconv.FormatUint(p.UID, 10)
	}, "name": func(p model.Process) string { return p.Command }, "state": func(p model.Process) string { return p.State }})
	if err != nil {
		return nil, err
	}
	sortField := o.Sort
	if sortField == "" {
		sortField = "pid"
	}
	items, err = SortItems(items, sortField, map[string]func(model.Process, model.Process) bool{"pid": func(a, b model.Process) bool { return a.PID < b.PID }, "cpu": lessProcessCPU, "memory": func(a, b model.Process) bool { return a.ResidentBytes < b.ResidentBytes }, "name": func(a, b model.Process) bool { return a.Command < b.Command }}, o.Reverse)
	if err != nil {
		return nil, err
	}
	items, err = LimitItems(items, o.Limit)
	if err != nil {
		return nil, err
	}
	return &model.ProcessList{Processes: items, CPUTimeTotal: list.CPUTimeTotal, TotalMemoryBytes: list.TotalMemoryBytes, Partial: list.Partial}, nil
}
func lessProcessCPU(a, b model.Process) bool {
	if a.CPUPercent != nil && b.CPUPercent != nil {
		return *a.CPUPercent < *b.CPUPercent
	}
	return a.CPUTime < b.CPUTime
}
func applyMemoryPercentages(list *model.ProcessList) {
	if list.TotalMemoryBytes == 0 {
		return
	}
	for i := range list.Processes {
		percent := float64(list.Processes[i].ResidentBytes) * 100 / float64(list.TotalMemoryBytes)
		list.Processes[i].MemoryPercent = &percent
	}
}
func applyProcessPercentages(before, after *model.ProcessList) {
	applyMemoryPercentages(after)
	if after.CPUTimeTotal <= before.CPUTimeTotal {
		return
	}
	old := make(map[int]uint64, len(before.Processes))
	for _, process := range before.Processes {
		old[process.PID] = process.CPUTime
	}
	totalDelta := after.CPUTimeTotal - before.CPUTimeTotal
	for i := range after.Processes {
		previous, ok := old[after.Processes[i].PID]
		if !ok || after.Processes[i].CPUTime < previous {
			continue
		}
		percent := float64(after.Processes[i].CPUTime-previous) * 100 / float64(totalDelta)
		after.Processes[i].CPUPercent = &percent
	}
}
func ParseProcessFilters(raw []string) ([]Filter, error) {
	filters := make([]Filter, 0, len(raw))
	for _, value := range raw {
		f, err := ParseFilter(value)
		if err != nil {
			return nil, fmt.Errorf("process filter: %w", err)
		}
		filters = append(filters, f)
	}
	return filters, nil
}

// Tree orders processes depth-first by PPID. Orphans remain roots so a parent
// exiting or a PID namespace boundary never hides a process.
func (s *Process) Tree() ([]ProcessTreeNode, error) {
	list, err := s.collector.Collect()
	if err != nil {
		return nil, err
	}
	byPID := map[int]*ProcessTreeNode{}
	for _, p := range list.Processes {
		copy := p
		byPID[p.PID] = &ProcessTreeNode{Process: copy}
	}
	roots := make([]*ProcessTreeNode, 0)
	for _, node := range byPID {
		parent, ok := byPID[node.Process.PPID]
		if !ok || node.Process.PID == node.Process.PPID {
			roots = append(roots, node)
			continue
		}
		parent.Children = append(parent.Children, node)
	}
	sortTreeNodes(roots)
	for _, root := range roots {
		sortTree(root)
	}
	result := make([]ProcessTreeNode, len(roots))
	for i, root := range roots {
		result[i] = *root
	}
	return result, nil
}

type ProcessTreeNode struct {
	Process  model.Process      `json:"process"`
	Children []*ProcessTreeNode `json:"children,omitempty"`
}

func sortTreeNodes(nodes []*ProcessTreeNode) {
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Process.PID < nodes[j].Process.PID })
}
func sortTree(node *ProcessTreeNode) {
	sortTreeNodes(node.Children)
	for _, child := range node.Children {
		sortTree(child)
	}
}
