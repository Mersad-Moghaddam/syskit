package service

import (
	"fmt"
	"sort"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

// ContainerMetricsReader reads optional resource counters for a process's
// cgroup. It is injected so service code does not depend on the platform layer.
type ContainerMetricsReader func(model.Process) (*model.ContainerMetrics, error)

type Container struct {
	collector ProcessCollector
	metrics   ContainerMetricsReader
}

func NewContainer(c ProcessCollector, readers ...ContainerMetricsReader) *Container {
	s := &Container{collector: c}
	if len(readers) > 0 {
		s.metrics = readers[0]
	}
	return s
}

func (s *Container) List() (*model.ContainerList, error) {
	processes, err := s.collector.Collect()
	if err != nil {
		return nil, err
	}
	byID := make(map[string]model.ContainerInfo)
	representatives := make(map[string]model.Process)
	for _, process := range processes.Processes {
		if process.ContainerID == "" {
			continue
		}
		container := byID[process.ContainerID]
		container.ID = process.ContainerID
		if container.Runtime == "" {
			container.Runtime = process.ContainerRuntime
		}
		if _, ok := representatives[process.ContainerID]; !ok {
			representatives[process.ContainerID] = process
		}
		container.PIDs++
		byID[container.ID] = container
	}
	result := &model.ContainerList{Containers: make([]model.ContainerInfo, 0, len(byID))}
	for _, container := range byID {
		container.Metrics = s.readMetrics(representatives[container.ID])
		result.Containers = append(result.Containers, container)
	}
	sort.Slice(result.Containers, func(i, j int) bool { return result.Containers[i].ID < result.Containers[j].ID })
	return result, nil
}

func (s *Container) Inspect(id string) (*model.ContainerDetail, error) {
	processes, err := s.collector.Collect()
	if err != nil {
		return nil, err
	}
	detail := &model.ContainerDetail{}
	for _, process := range processes.Processes {
		if process.ContainerID != id {
			continue
		}
		if detail.ID == "" {
			detail.ID, detail.Runtime = process.ContainerID, process.ContainerRuntime
		}
		detail.PIDs++
		detail.Processes = append(detail.Processes, process)
	}
	if detail.ID == "" {
		return nil, fmt.Errorf("container %q not found", id)
	}
	detail.Metrics = s.readMetrics(detail.Processes[0])
	sort.Slice(detail.Processes, func(i, j int) bool { return detail.Processes[i].PID < detail.Processes[j].PID })
	return detail, nil
}

func (s *Container) readMetrics(process model.Process) *model.ContainerMetrics {
	if s.metrics == nil {
		return nil
	}
	metrics, err := s.metrics(process)
	if err != nil {
		return nil
	}
	return metrics
}
