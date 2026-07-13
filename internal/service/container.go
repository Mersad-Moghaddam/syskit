package service

import (
	"sort"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

// Container groups cgroup-associated processes into a runtime-independent
// container view.
type Container struct{ collector ProcessCollector }

func NewContainer(c ProcessCollector) *Container { return &Container{collector: c} }

func (s *Container) List() (*model.ContainerList, error) {
	processes, err := s.collector.Collect()
	if err != nil {
		return nil, err
	}
	byID := make(map[string]model.ContainerInfo)
	for _, process := range processes.Processes {
		if process.ContainerID == "" {
			continue
		}
		container := byID[process.ContainerID]
		container.ID = process.ContainerID
		if container.Runtime == "" {
			container.Runtime = process.ContainerRuntime
		}
		container.PIDs++
		byID[container.ID] = container
	}
	result := &model.ContainerList{Containers: make([]model.ContainerInfo, 0, len(byID))}
	for _, container := range byID {
		result.Containers = append(result.Containers, container)
	}
	sort.Slice(result.Containers, func(i, j int) bool { return result.Containers[i].ID < result.Containers[j].ID })
	return result, nil
}
